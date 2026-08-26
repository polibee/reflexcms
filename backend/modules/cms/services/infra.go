package services

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/support/carbon"

	"reflexcms/backend/app/facades"
	cmsevents "reflexcms/backend/modules/cms/events"
)

func facadesQuery() orm.Query { return facades.Orm().Query() }

func carbonNowPtr() *carbon.DateTime {
	return carbon.NewDateTime(carbon.Now())
}

// dispatchPublished fires the ArticlePublished event through the framework
// queue (sync connection runs it inline, which keeps the M2 DoD observable
// in logs). Value dispatch: listeners match by struct equality.
func dispatchPublished(articleID uint64) {
	facades.Event().Job(cmsevents.ArticlePublished{}, []event.Arg{
		{Value: articleID, Type: "uint64"},
	}).Dispatch()
}
