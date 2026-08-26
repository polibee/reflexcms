package events

import (
	"github.com/goravel/framework/contracts/event"
)

// CommentModerated is dispatched after a moderation decision changes a
// comment's status. Value receiver + value dispatch: the framework matches
// listeners by interface equality, so pointer instances would never match.
type CommentModerated struct{}

func (e CommentModerated) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
