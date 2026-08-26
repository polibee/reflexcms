package events

import (
	"github.com/goravel/framework/contracts/event"
)

// ReplyCreated is dispatched after a reply is committed. Payload order:
//
//	args[0] = reply ID, args[1] = topic ID,
//	args[2] = topic author user ID, args[3] = reply author user ID
//
// Composition roots forward this to notification; the forum module itself
// never imports it. Value receiver + value dispatch: the framework matches
// listeners by interface equality, so pointer instances would never match.
type ReplyCreated struct{}

func (e ReplyCreated) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
