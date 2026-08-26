package services

import (
	"testing"

	cmsmodels "reflexcms/backend/modules/cms/models"
)

func TestCanTransition(t *testing.T) {
	legal := [][2]string{
		{"draft", "published"}, {"draft", "archived"},
		{"published", "draft"}, {"published", "archived"},
		{"archived", "draft"},
	}
	for _, tc := range legal {
		if !CanTransition(tc[0], tc[1]) {
			t.Errorf("expected %s → %s to be legal", tc[0], tc[1])
		}
	}

	illegal := [][2]string{
		{"draft", "draft"}, {"published", "published"},
		{"archived", "published"}, {"archived", "archived"},
		{"unknown", "draft"}, {"draft", ""},
	}
	for _, tc := range illegal {
		if CanTransition(tc[0], tc[1]) {
			t.Errorf("expected %s → %s to be rejected", tc[0], tc[1])
		}
	}
}

func TestTransitionAction(t *testing.T) {
	for action, want := range map[string]string{
		"publish": cmsmodels.ArticlePublished, "unpublish": cmsmodels.ArticleDraft,
		"archive": cmsmodels.ArticleArchived, "restore": cmsmodels.ArticleDraft,
	} {
		got, ok := TransitionAction(action)
		if !ok || got != want {
			t.Errorf("TransitionAction(%q) = %q,%v; want %q", action, got, ok, want)
		}
	}
	if _, ok := TransitionAction("delete"); ok {
		t.Error(`TransitionAction("delete") should not resolve`)
	}
}
