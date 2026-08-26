package services

import (
	"errors"
	"fmt"

	"github.com/goravel/framework/contracts/database/orm"

	"reflexcms/backend/app/facades"
	forummodels "reflexcms/backend/modules/forum/models"
)

var ErrInvalidLikeable = errors.New("invalid likeable type")

// ToggleLike flips the caller's like on a topic or reply. The composite
// unique index makes the insert idempotent: an affected row means "now
// liked", no affected row means it already existed → delete to unlike. The
// denormalised counter is then recomputed exactly, so any drift self-heals.
//
// The target table name never comes from user input — each branch uses its
// own hard-coded statement; every value is a bound parameter.
func ToggleLike(userID uint64, likeableType string, likeableID uint64) (liked bool, count int64, err error) {
	switch likeableType {
	case forummodels.LikeableTopic, forummodels.LikeableReply:
	default:
		return false, 0, ErrInvalidLikeable
	}

	txErr := facades.Orm().Transaction(func(tx orm.Query) error {
		result, execErr := tx.Exec(`
			INSERT INTO likes (user_id, likeable_type, likeable_id, created_at)
			VALUES (?, ?, ?, NOW())
			ON CONFLICT DO NOTHING
		`, userID, likeableType, likeableID)
		if execErr != nil {
			return execErr
		}

		inserted := result != nil && result.RowsAffected > 0
		if !inserted {
			if _, err := tx.Exec(`
				DELETE FROM likes
				WHERE user_id = ? AND likeable_type = ? AND likeable_id = ?
			`, userID, likeableType, likeableID); err != nil {
				return err
			}
			liked = false
		} else {
			liked = true
		}

		return refreshLikeCount(tx, likeableType, likeableID)
	})
	if txErr != nil {
		return false, 0, txErr
	}

	count, err = currentLikeCount(likeableType, likeableID)
	return liked, count, err
}

func refreshLikeCount(tx orm.Query, likeableType string, likeableID uint64) error {
	switch likeableType {
	case forummodels.LikeableTopic:
		_, err := tx.Exec(`
			UPDATE topics SET like_count =
				(SELECT COUNT(*) FROM likes WHERE likeable_type = ? AND likeable_id = ?)
			WHERE id = ?
		`, likeableType, likeableID, likeableID)
		return err
	case forummodels.LikeableReply:
		_, err := tx.Exec(`
			UPDATE replies SET like_count =
				(SELECT COUNT(*) FROM likes WHERE likeable_type = ? AND likeable_id = ?)
			WHERE id = ?
		`, likeableType, likeableID, likeableID)
		return err
	}
	return ErrInvalidLikeable
}

func currentLikeCount(likeableType string, likeableID uint64) (int64, error) {
	var rows []map[string]any
	err := query().Table(likeableTableForRead(likeableType)).
		Where("id = ?", likeableID).
		Limit(1).
		Get(&rows)
	if err != nil || len(rows) == 0 {
		return 0, err
	}
	switch v := rows[0]["like_count"].(type) {
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case []byte:
		return parseInt64(string(v))
	case string:
		return parseInt64(v)
	}
	return 0, nil
}

func likeableTableForRead(t string) string {
	if t == forummodels.LikeableReply {
		return "replies"
	}
	return "topics"
}

func parseInt64(s string) (int64, error) {
	var n int64
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
