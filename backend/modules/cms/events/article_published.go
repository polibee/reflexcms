package events

import (
	"github.com/goravel/framework/contracts/event"
)

// ArticlePublished is dispatched whenever an article enters the published
// state (first publish or re-publish). Value receiver + value dispatch: the
// framework matches listeners by interface equality, so pointer instances
// would never match.
type ArticlePublished struct{}

func (e ArticlePublished) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
