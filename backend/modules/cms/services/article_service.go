package services

import (
	"fmt"

	safefetch "reflexcms/backend/app/support/safefetch"
	cmsmodels "reflexcms/backend/modules/cms/models"
)

var validator = safefetch.NewValidator()

// ValidateCover enforces the outbound-URL policy on the cover field: only
// public http(s) URLs may be stored. No fetch happens here, but rejecting
// unsafe targets at write time keeps every future consumer safe.
func ValidateCover(raw string) error {
	if raw == "" {
		return nil
	}
	if _, err := validator.ValidateURL(raw); err != nil {
		return fmt.Errorf("cover url rejected: %w", err)
	}
	return nil
}

// EnsureSlug derives a unique slug from the title when none was provided,
// or validates/uniquifies a hand-written one.
func EnsureSlug(article *cmsmodels.Article, exists func(slug string) (bool, error)) error {
	base := article.Slug
	if base == "" {
		base = Slugify(article.Title)
	}
	slug, err := EnsureUnique(base, 191, exists)
	if err != nil {
		return err
	}
	article.Slug = slug
	return nil
}

// SyncArticleTags replaces the article's tag links with the given tag names.
// Missing tags are created (slug derived). Empty input clears all links.
func SyncArticleTags(articleID uint64, names []string) error {
	pivot := cmsmodels.ArticleTag{}
	if _, err := facadesQuery().Table(pivot.TableName()).Where("article_id = ?", articleID).Delete(); err != nil {
		return err
	}
	if len(names) == 0 {
		return nil
	}

	tagIDs := make([]uint64, 0, len(names))
	for _, name := range names {
		if name == "" {
			continue
		}
		var tag cmsmodels.Tag
		// FirstOrFail, not First: goravel's First returns a zero struct with
		// nil error when no row matches, which would forge tag IDs of 0.
		err := facadesQuery().Where("name = ?", name).FirstOrFail(&tag)
		fmt.Println("DEBUG syncTag", name, "firstErr:", err, "resolvedID:", tag.ID)
		if err != nil {
			slug, slugErr := EnsureUnique(Slugify(name), 191, func(s string) (bool, error) {
				var count int64
				count, qErr := facadesQuery().Model(cmsmodels.Tag{}).Where("slug = ?", s).Count()
				return count > 0, qErr
			})
			if slugErr != nil {
				return slugErr
			}
			tag = cmsmodels.Tag{Name: name, Slug: slug}
			if err := facadesQuery().Create(&tag); err != nil {
				return err
			}
			fmt.Println("DEBUG created tag", name, "id:", tag.ID)
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	if len(tagIDs) == 0 {
		return nil
	}
	rows := make([]cmsmodels.ArticleTag, 0, len(tagIDs))
	for _, id := range tagIDs {
		rows = append(rows, cmsmodels.ArticleTag{ArticleID: articleID, TagID: id})
	}
	return facadesQuery().Create(&rows)
}

// Transition moves an article to the target status through the state machine
// and dispatches ArticlePublished on entry into published.
func Transition(article *cmsmodels.Article, to string) error {
	if !CanTransition(article.Status, to) {
		return fmt.Errorf("illegal transition %s → %s", article.Status, to)
	}
	article.Status = to
	if to == cmsmodels.ArticlePublished && article.PublishedAt == nil {
		article.PublishedAt = carbonNowPtr()
	}
	if err := facadesQuery().Save(article); err != nil {
		return err
	}
	if to == cmsmodels.ArticlePublished {
		dispatchPublished(article.ID)
	}
	return nil
}
