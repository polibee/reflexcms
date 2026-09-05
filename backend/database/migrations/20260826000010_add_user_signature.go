package migrations

import (
	"reflexcms/backend/app/facades"
)

// M20260826000010AddUserSignature adds the user signature (short markdown,
// text + links only; rendered client-side through the escaping renderer).
type M20260826000010AddUserSignature struct{}

func (r *M20260826000010AddUserSignature) Signature() string {
	return "20260826000010_add_user_signature"
}

func (r *M20260826000010AddUserSignature) Up() error {
	if _, err := facades.Orm().Query().Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS signature VARCHAR(500) NULL
	`); err != nil {
		return err
	}
	return nil
}

func (r *M20260826000010AddUserSignature) Down() error {
	_, err := facades.Orm().Query().Exec(`
		ALTER TABLE users DROP COLUMN IF EXISTS signature
	`)
	return err
}
