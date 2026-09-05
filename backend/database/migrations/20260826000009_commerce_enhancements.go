package migrations

import (
	"reflexcms/backend/app/facades"
)

// M20260826000009CommerceEnhancements adds:
//   - products.type ('general' | 'invite') for invite-code goods
//   - orders.granted_code for deliverable goods (invite codes)
//   - users.muted_until / banned_until for moderator governance
type M20260826000009CommerceEnhancements struct{}

func (r *M20260826000009CommerceEnhancements) Signature() string {
	return "20260826000009_commerce_enhancements"
}

func (r *M20260826000009CommerceEnhancements) Up() error {
	// Pure constant DDL statements — no user input involved, so binding
	// rules do not apply; raw SQL goes through the ORM connection on
	// purpose (see create_articles_table for the facades.DB() caveat).
	if _, err := facades.Orm().Query().Exec(`
		ALTER TABLE products ADD COLUMN IF NOT EXISTS type VARCHAR(16) DEFAULT 'general'
	`); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(`
		ALTER TABLE orders ADD COLUMN IF NOT EXISTS granted_code VARCHAR(64) NULL
	`); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS muted_until TIMESTAMPTZ NULL
	`); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS banned_until TIMESTAMPTZ NULL
	`); err != nil {
		return err
	}
	return nil
}

func (r *M20260826000009CommerceEnhancements) Down() error {
	return nil
}
