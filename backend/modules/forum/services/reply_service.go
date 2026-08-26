package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/event"

	"reflexcms/backend/app/facades"
	forumevents "reflexcms/backend/modules/forum/events"
	forummodels "reflexcms/backend/modules/forum/models"
)

var (
	ErrTopicClosed  = errors.New("topic is closed for replies")
	ErrContentEmpty = errors.New("reply content is required")
)

func query() orm.Query { return facades.Orm().Query() }

// ReplyToTopic inserts a reply and keeps the denormalised topic counters
// consistent inside a single transaction:
//
//  1. lock the topic row (serialises concurrent floors)
//  2. insert the reply with floor = max(existing)+1
//  3. bump reply_count, refresh last_reply_*
//
// The raw statements use bound parameters exclusively ($n placeholders).
func ReplyToTopic(topicID, userID, parentID uint64, content string) (*forummodels.Reply, error) {
	if strings.TrimSpace(content) == "" {
		return nil, ErrContentEmpty
	}

	reply := &forummodels.Reply{}
	txErr := facades.Orm().Transaction(func(tx orm.Query) error {
		// 1. Serialise concurrent replies on this topic.
		if _, err := tx.Exec(
			"SELECT id FROM topics WHERE id = ? AND deleted_at IS NULL FOR UPDATE",
			topicID,
		); err != nil {
			return err
		}

		var statuses []string
		if err := tx.Table("topics").
			Where("id = ?", topicID).
			Pluck("status", &statuses); err == nil && len(statuses) > 0 &&
			statuses[0] == forummodels.TopicClosed {
			return ErrTopicClosed
		}

		// 2. Insert with the next free floor.
		if _, err := tx.Exec(`
			INSERT INTO replies (topic_id, user_id, parent_id, content, floor, created_at, updated_at)
			VALUES (?, ?, ?, ?,
				(SELECT COALESCE(MAX(floor), 0) + 1 FROM replies WHERE topic_id = ? AND deleted_at IS NULL),
				NOW(), NOW())
			RETURNING id, floor
		`, topicID, userID, parentID, content, topicID); err != nil {
			return err
		}
		if err := tx.Table("replies").
			Where("topic_id = ?", topicID).
			OrderBy("id", "desc").
			Limit(1).
			First(reply); err != nil {
			return fmt.Errorf("reload inserted reply: %w", err)
		}

		// 3. Refresh counters in the same transaction.
		if _, err := tx.Exec(`
			UPDATE topics SET
				reply_count = reply_count + 1,
				last_reply_at = NOW(),
				last_reply_user_id = ?
			WHERE id = ?
		`, userID, topicID); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	// Notify wiring data for the composition root: topic author + replier.
	var topicAuthor uint64
	if err := query().Table("topics").
		Where("id = ?", topicID).
		Pluck("user_id", &topicAuthor); err == nil {
		facades.Event().Job(forumevents.ReplyCreated{}, []event.Arg{
			{Value: reply.ID},
			{Value: topicID},
			{Value: topicAuthor},
			{Value: userID},
		}).Dispatch()
	}

	return reply, nil
}
