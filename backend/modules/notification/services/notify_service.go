package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/event"

	"reflexcms/backend/app/facades"
	notificationmodels "reflexcms/backend/modules/notification/models"
)

// Notify persists an in-app notification for the user. data is encoded as
// JSON so each trigger decides its own payload shape.
func Notify(userID uint64, ntype string, data map[string]any) error {
	if userID == 0 {
		return nil // system-scoped events have no recipient
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	notification := notificationmodels.Notification{
		UserID: userID,
		Type:   ntype,
		Data:   string(payload),
	}
	return facades.Orm().Query().Create(&notification)
}

// MarkAllRead stamps every unread notification of the caller.
func MarkAllRead(userID uint64) (int64, error) {
	result, err := facades.Orm().Query().
		Model(&notificationmodels.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", now())
	if err != nil {
		return 0, err
	}
	rows := int64(0)
	if result != nil {
		rows = result.RowsAffected
	}
	return rows, nil
}

/* ---------------- composition-ready listeners ---------------- */

// FuncListener adapts a closure to the framework listener contract so
// composition roots can wire domain events to Notify without modules
// importing each other.
type FuncListener struct {
	sig string
	fn  func(args ...any) error
}

func NewFuncListener(signature string, fn func(args ...any) error) event.Listener {
	return &FuncListener{sig: signature, fn: fn}
}

func (l *FuncListener) Signature() string { return l.sig }

func (l *FuncListener) Queue(args ...any) event.Queue { return event.Queue{Enable: false} }

func (l *FuncListener) Handle(args ...any) error { return l.fn(args...) }

var _ = fmt.Sprintf // keep fmt until listeners format payloads

/* ---------------- @mentions ---------------- */

// mentionPattern captures @username tokens; usernames are \w plus CJK and
// hyphen, capped at 32 chars.
var mentionPattern = regexp.MustCompile(`@([\p{Han}\w-]{1,32})`)

const maxMentionsPerPost = 5

// NotifyMentions scans text for @username tokens, resolves them to active
// users (excluding the author), and drops a mention notification for each —
// at most maxMentionsPerPost per post, self-mentions skipped. Best-effort:
// failures never block the calling action.
func NotifyMentions(fromUserID uint64, fromName, text, linkType string, linkID uint64, excerpt string) {
	matches := mentionPattern.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return
	}

	seen := map[string]bool{}
	names := make([]string, 0, maxMentionsPerPost)
	for _, m := range matches {
		name := m[1]
		if name == fromName || seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
		if len(names) >= maxMentionsPerPost {
			break
		}
	}
	if len(names) == 0 {
		return
	}

	args := make([]any, 0, len(names))
	placeholders := make([]string, 0, len(names))
	for _, n := range names {
		args = append(args, n)
		placeholders = append(placeholders, "?")
	}

	var users []map[string]any
	if err := facades.Orm().Query().Table("users").
		Where("username IN ("+strings.Join(placeholders, ", ")+") AND status = 'active' AND deleted_at IS NULL", args...).
		Select("id", "username").
		Get(&users); err != nil {
		return
	}

	excerptRunes := []rune(excerpt)
	if len(excerptRunes) > 80 {
		excerptRunes = excerptRunes[:80]
	}
	for _, u := range users {
		id, ok := u["id"].(int64)
		if !ok || uint64(id) == fromUserID {
			continue
		}
		_ = Notify(uint64(id), "mention", map[string]any{
			"from_user_id": fromUserID,
			"from_name":    fromName,
			"link_type":    linkType,
			"link_id":      linkID,
			"excerpt":      string(excerptRunes),
		})
	}
}
