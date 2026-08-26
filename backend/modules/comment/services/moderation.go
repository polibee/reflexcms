package services

import (
	"fmt"

	"github.com/goravel/framework/contracts/event"

	"reflexcms/backend/app/facades"
	commentevents "reflexcms/backend/modules/comment/events"
	commentmodels "reflexcms/backend/modules/comment/models"
)

// Approve / Reject move a comment through the moderation state machine and
// keep the host article's comment_count exact (only approved comments count).
// Every decision dispatches CommentModerated so composition roots can notify
// the comment author without this module knowing about notifications.
func Approve(commentID uint64) (*commentmodels.Comment, error) {
	return moderate(commentID, commentmodels.CommentApproved)
}

func Reject(commentID uint64) (*commentmodels.Comment, error) {
	return moderate(commentID, commentmodels.CommentRejected)
}

func moderate(commentID uint64, to string) (*commentmodels.Comment, error) {
	var comment commentmodels.Comment
	if err := facades.Orm().Query().FindOrFail(&comment, commentID); err != nil {
		return nil, fmt.Errorf("comment not found")
	}
	if !CanTransition(comment.Status, to) {
		return nil, fmt.Errorf("illegal transition %s → %s", comment.Status, to)
	}

	comment.Status = to
	if err := facades.Orm().Query().Save(&comment); err != nil {
		return nil, err
	}

	if err := RecountArticle(comment.ArticleID); err != nil {
		return nil, err
	}

	facades.Event().Job(commentevents.CommentModerated{}, []event.Arg{
		{Value: comment.ID},
		{Value: comment.UserID},
		{Value: to},
	}).Dispatch()

	return &comment, nil
}

// CanTransition: pending → approved|rejected; later flips allowed between
// approved/rejected; spam is terminal.
func CanTransition(from, to string) bool {
	if from == to {
		return false
	}
	switch from {
	case commentmodels.CommentPending:
		return to == commentmodels.CommentApproved || to == commentmodels.CommentRejected || to == commentmodels.CommentSpam
	case commentmodels.CommentApproved:
		return to == commentmodels.CommentRejected || to == commentmodels.CommentSpam
	case commentmodels.CommentRejected:
		return to == commentmodels.CommentApproved || to == commentmodels.CommentSpam
	}
	return false // spam or unknown
}

// RecountArticle re-derives article.comment_count from approved rows, so the
// denormalised counter can never drift from the truth.
func RecountArticle(articleID uint64) error {
	_, err := facades.Orm().Query().Exec(`
		UPDATE articles SET comment_count =
			(SELECT COUNT(*) FROM comments WHERE article_id = ? AND status = ? AND deleted_at IS NULL)
		WHERE id = ?
	`, articleID, commentmodels.CommentApproved, articleID)
	return err
}
