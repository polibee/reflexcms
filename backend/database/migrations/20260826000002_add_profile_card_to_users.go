package migrations

import (
	"reflexcms/backend/app/facades"
)

type M20260826000002AddProfileCardToUsers struct{}

func (r *M20260826000002AddProfileCardToUsers) Signature() string {
	return "20260826000002_add_profile_card_to_users"
}

func (r *M20260826000002AddProfileCardToUsers) Up() error {
	if !facades.Schema().HasColumn("users", "profile_card") {
		// Pure constant DDL — no user input involved, so binding rules do not
		// apply; raw SQL goes through the ORM connection on purpose (see
		// create_articles_table for the facades.DB() vs facades.Orm() caveat).
		if _, err := facades.Orm().Query().Exec(`
			ALTER TABLE users ADD COLUMN profile_card TEXT NULL
		`); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260826000002AddProfileCardToUsers) Down() error {
	_, err := facades.Orm().Query().Exec(`
		ALTER TABLE users DROP COLUMN IF EXISTS profile_card
	`)
	return err
}