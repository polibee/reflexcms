package services

// Slug helpers moved to the shared support package so the forum module can
// reuse them without importing cms internals (plan §7.2).

import slug "reflexcms/backend/app/support/slug"

var (
	Slugify      = slug.Slugify
	EnsureUnique = slug.EnsureUnique
)
