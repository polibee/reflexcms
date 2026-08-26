// Package products is the ONLY place allowed to know every module. Each
// constructor returns one product composition; each becomes a deployable
// binary via cmd/<product>/main.go (docs/模块化架构与产品组合.md).
package products

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/event"

	ainew "reflexcms/backend/modules/ai"
	"reflexcms/backend/modules/access"
	auditmodule "reflexcms/backend/modules/audit"
	authmodule "reflexcms/backend/modules/auth"
	cmsmodule "reflexcms/backend/modules/cms"
	cmsevents "reflexcms/backend/modules/cms/events"
	commentevents "reflexcms/backend/modules/comment/events"
	commentmodule "reflexcms/backend/modules/comment"
	forumevents "reflexcms/backend/modules/forum/events"
	forummodule "reflexcms/backend/modules/forum"
	layoutmodule "reflexcms/backend/modules/layout"
	notificationmodule "reflexcms/backend/modules/notification"
	notificationservices "reflexcms/backend/modules/notification/services"
	pagesmodule "reflexcms/backend/modules/pages"
	settingsmodule "reflexcms/backend/modules/settings"
	"reflexcms/backend/kernel"
)

/* ---------------- products ---------------- */

// Full is the default combined CMS+BBS product.
func Full() []kernel.Module {
	mods := core()
	mods = append(mods,
		cmsmodule.New(), commentmodule.New(),
		forummodule.New(),
		settingsmodule.New(), notificationmodule.New(), auditmodule.New(),
		layoutmodule.New(), pagesmodule.New(),
		aiWiring{}, notifyWiring{},
	)
	return mods
}

// Blog is the articles-only product.
func Blog() []kernel.Module {
	mods := core()
	mods = append(mods,
		cmsmodule.New(), commentmodule.New(),
		settingsmodule.New(), notificationmodule.New(), auditmodule.New(),
		layoutmodule.New(), pagesmodule.New(),
		aiWiring{}, notifyWiring{},
	)
	return mods
}

// Forum is the boards-only product.
func Forum() []kernel.Module {
	mods := core()
	mods = append(mods,
		forummodule.New(),
		settingsmodule.New(), notificationmodule.New(), auditmodule.New(),
		layoutmodule.New(), pagesmodule.New(),
		notifyWiring{},
	)
	return mods
}

func core() []kernel.Module {
	return []kernel.Module{access.New(), authmodule.New()}
}

/* ---------------- composition-time wirings ---------------- */

// aiWiring connects published articles to the AI placeholder listener.
type aiWiring struct{}

func (aiWiring) Name() string                   { return "ai-wiring" }
func (aiWiring) Migrations() []schema.Migration { return nil }
func (aiWiring) Routes()                        {}
func (aiWiring) Boot()                          {}

func (aiWiring) Events() map[event.Event][]event.Listener {
	return map[event.Event][]event.Listener{
		cmsevents.ArticlePublished{}: {ainew.NewSummaryListener()},
	}
}

// notifyWiring forwards domain events into in-app notifications. Emitters
// stay unaware of notification; this is the only coupling point.
type notifyWiring struct{}

func (notifyWiring) Name() string                   { return "notification-wiring" }
func (notifyWiring) Migrations() []schema.Migration { return nil }
func (notifyWiring) Routes()                        {}
func (notifyWiring) Boot()                          {}

func (notifyWiring) Events() map[event.Event][]event.Listener {
	return map[event.Event][]event.Listener{
		forumevents.ReplyCreated{}:    {notificationservices.NewFuncListener("wiring.reply-created", onReplyCreated)},
		commentevents.CommentModerated{}: {notificationservices.NewFuncListener("wiring.comment-moderated", onCommentModerated)},
	}
}

// args: replyID, topicID, topicAuthorUserID, replyAuthorUserID
func onReplyCreated(args ...any) error {
	if len(args) < 4 {
		return nil
	}
	topicAuthor := toU64(args[2])
	replyAuthor := toU64(args[3])
	if topicAuthor == 0 || topicAuthor == replyAuthor {
		return nil // no self-notification
	}
	return notificationservices.Notify(topicAuthor, "forum.reply", map[string]any{
		"title":    "Your topic has a new reply",
		"topic_id": toU64(args[1]),
		"reply_id": toU64(args[0]),
	})
}

// args: commentID, commentAuthorUserID, action
func onCommentModerated(args ...any) error {
	if len(args) < 3 {
		return nil
	}
	author := toU64(args[1])
	if author == 0 {
		return nil
	}
	action := ""
	if s, ok := args[2].(string); ok {
		action = s
	}
	return notificationservices.Notify(author, "comment."+action, map[string]any{
		"title":      "Your comment was " + action,
		"comment_id": toU64(args[0]),
		"status":     action,
	})
}

func toU64(v any) uint64 {
	switch n := v.(type) {
	case uint64:
		return n
	case int:
		return uint64(n)
	case int64:
		return uint64(n)
	case uint:
		return uint64(n)
	}
	return 0
}
