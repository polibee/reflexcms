package services

import (
	cmsmodels "reflexcms/backend/modules/cms/models"
)

// Article lifecycle transitions (plan §M2):
//
//	draft ⇄ published, published → archived, archived → draft
//
// Every other change is rejected; status is only ever mutated through the
// dedicated action endpoints, never through the generic update payload.
var articleTransitions = map[string][]string{
	cmsmodels.ArticleDraft:     {cmsmodels.ArticleScheduled, cmsmodels.ArticlePublished, cmsmodels.ArticleArchived},
	cmsmodels.ArticleScheduled: {cmsmodels.ArticlePublished, cmsmodels.ArticleDraft, cmsmodels.ArticleArchived},
	cmsmodels.ArticlePublished: {cmsmodels.ArticleDraft, cmsmodels.ArticleArchived},
	cmsmodels.ArticleArchived:  {cmsmodels.ArticleDraft},
}

func CanTransition(from, to string) bool {
	if from == to {
		return false
	}
	for _, next := range articleTransitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

func IsArticleStatus(s string) bool {
	switch s {
	case cmsmodels.ArticleDraft, cmsmodels.ArticleScheduled, cmsmodels.ArticlePublished, cmsmodels.ArticleArchived:
		return true
	}
	return false
}

// TransitionAction maps admin action names onto target statuses.
func TransitionAction(action string) (to string, ok bool) {
	switch action {
	case "publish":
		return cmsmodels.ArticlePublished, true
	case "schedule":
		return cmsmodels.ArticleScheduled, true
	case "unpublish":
		return cmsmodels.ArticleDraft, true
	case "archive":
		return cmsmodels.ArticleArchived, true
	case "restore":
		return cmsmodels.ArticleDraft, true
	}
	return "", false
}
