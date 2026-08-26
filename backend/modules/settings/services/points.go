package services

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"

	"reflexcms/backend/app/facades"
)

// Point reasons for the transaction log.
const (
	PointReasonArticle      = "article_publish"
	PointReasonTopic        = "topic_create"
	PointReasonReply        = "reply_create"
	PointReasonComment      = "comment_approved"
	PointReasonSignin       = "daily_signin"
	PointReasonInviteCode   = "invite_code_generated"
	PointReasonInviteUsed   = "invite_code_used"
)

// settingInt reads an integer-valued settings key.
func settingInt(key string, fallback int) int {
	raw := Get(key)
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(strings.Trim(raw, `"`))
	if err != nil {
		return fallback
	}
	return n
}

// Points for each action (configurable via admin settings).
func EarnArticle() int      { return settingInt("points.earn_article", 10) }
func EarnTopic() int        { return settingInt("points.earn_topic", 5) }
func EarnReply() int        { return settingInt("points.earn_reply", 2) }
func EarnComment() int      { return settingInt("points.earn_comment", 3) }
func EarnSignin() int       { return settingInt("points.earn_signin", 1) }
func EarnInvite() int       { return settingInt("points.earn_invite", 20) }
func CostInviteCode() int   { return settingInt("points.cost_invite_code", 50) }

// CurrencyName returns the display name of the community currency.
func CurrencyName() string { return Get("points.currency_name") }

// Daily earn caps per source (currency amount per day, admin-settable).
func TopicDailyCap() int { return settingInt("points.topic_daily_cap", 20) }
func ReplyDailyCap() int { return settingInt("points.reply_daily_cap", 20) }

// LevelThresholds parses the admin-set CSV of cumulative earn thresholds
// ("0,100,300,600,1000"), ascending; Lv1 starts at the first value.
func LevelThresholds() []int {
	raw := strings.Trim(Get("points.level_thresholds"), `"`)
	if raw == "" {
		raw = "0,100,300,600,1000,1500,2200,3000"
	}
	var out []int
	for _, part := range strings.Split(raw, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			continue
		}
		out = append(out, n)
	}
	if len(out) == 0 {
		out = []int{0, 100, 300, 600, 1000, 1500, 2200, 3000}
	}
	sort.Ints(out)
	return out
}

// toInt normalises the numeric types different drivers surface for INTEGER
// columns scanned into map[string]any.
func toInt(v any) int {
	switch n := v.(type) {
	case int64:
		return int(n)
	case int32:
		return int(n)
	case int:
		return n
	case uint64:
		return int(n)
	case uint:
		return int(n)
	case float64:
		return int(n)
	}
	return 0
}

// DailyEarned sums today's positive earnings for one reason.
func DailyEarned(userID uint64, reason string) int {
	var rows []map[string]any
	if err := facades.Orm().Query().Table("point_transactions").
		Where("user_id = ? AND reason = ? AND amount > 0 AND created_at >= CURRENT_DATE", userID, reason).
		Select("COALESCE(SUM(amount), 0) AS total").
		Get(&rows); err != nil || len(rows) == 0 {
		return 0
	}
	return toInt(rows[0]["total"])
}

// HasReasonToday reports whether the user already has a transaction with
// this reason today (used for once-per-day actions like sign-in).
func HasReasonToday(userID uint64, reason string) bool {
	var count int64
	count, _ = facades.Orm().Query().Table("point_transactions").
		Where("user_id = ? AND reason = ? AND created_at >= CURRENT_DATE", userID, reason).
		Count()
	return count > 0
}

// TotalEarned returns the lifetime earned total driving the level.
func TotalEarned(userID uint64) int {
	var rows []map[string]any
	if err := facades.Orm().Query().Table("user_points").
		Where("user_id = ?", userID).Limit(1).Get(&rows); err != nil || len(rows) == 0 {
		return 0
	}
	return toInt(rows[0]["total_earned"])
}

// AwardDaily grants `earn` points for `reason` unless today's earnings for
// that reason already reached `cap`. Returns the granted amount (0 = capped).
func AwardDaily(userID uint64, reason string, earn, cap int) (int, error) {
	if earn <= 0 || cap <= 0 {
		return 0, nil
	}
	if DailyEarned(userID, reason) >= cap {
		return 0, nil
	}
	if err := GrantPoints(userID, earn, reason); err != nil {
		return 0, err
	}
	return earn, nil
}

// LevelInfo locates total within the threshold ladder: 1-based level, the
// current threshold, and the next threshold (0 = max level reached).
func LevelInfo(total int) (level, curThreshold, nextThreshold int) {
	th := LevelThresholds()
	level = 1
	curThreshold = th[0]
	nextThreshold = 0
	for i := 1; i < len(th); i++ {
		if total >= th[i] {
			level = i + 1
			curThreshold = th[i]
		} else {
			nextThreshold = th[i]
			break
		}
	}
	return level, curThreshold, nextThreshold
}

// Grant points to a user inside a DB transaction.
func GrantPoints(userID uint64, amount int, reason string) error {
	if userID == 0 || amount == 0 {
		return nil
	}
	return facades.Orm().Transaction(func(tx orm.Query) error {
		// Upsert user_points balance.
		if _, err := tx.Exec(`
			INSERT INTO user_points (user_id, balance, total_earned, updated_at)
			VALUES (?, ?, ?, NOW())
			ON CONFLICT (user_id) DO UPDATE
			SET balance = user_points.balance + ?,
			    total_earned = user_points.total_earned + GREATEST(?, 0),
			    updated_at = NOW()
		`, userID, amount, amount, amount, amount); err != nil {
			return err
		}

		_, err := tx.Exec(`
			INSERT INTO point_transactions (user_id, amount, reason, created_at)
			VALUES (?, ?, ?, NOW())
		`, userID, amount, reason)
		return err
	})
}

// SpendPoints deducts points if the balance is sufficient; returns error on insufficient funds.
func SpendPoints(userID uint64, amount int, reason string) error {
	var balance int
	err := facades.Orm().Query().Table("user_points").
		Where("user_id = ?", userID).
		Pluck("balance", &balance)
	if err != nil || balance < amount {
		return fmt.Errorf("余额不足（当前 %d，需要 %d）", balance, amount)
	}

	return facades.Orm().Transaction(func(tx orm.Query) error {
		if _, err := tx.Exec(`
			UPDATE user_points SET balance = balance - ?, updated_at = NOW()
			WHERE user_id = ?
		`, amount, userID); err != nil {
			return err
		}
		_, err := tx.Exec(`
			INSERT INTO point_transactions (user_id, amount, reason, created_at)
			VALUES (?, -?, ?, NOW())
		`, userID, amount, reason)
		return err
	})
}

// GetBalance returns the current point balance for a user.
func GetBalance(userID uint64) int {
	var rows []map[string]any
	if err := facades.Orm().Query().Table("user_points").
		Where("user_id = ?", userID).Limit(1).Get(&rows); err != nil || len(rows) == 0 {
		return 0
	}
	return toInt(rows[0]["balance"])
}

var _ = facades.Cache
