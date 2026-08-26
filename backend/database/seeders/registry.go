package seeders

import (
	"github.com/goravel/framework/contracts/database/seeder"
)

// AccessSeeders are owned by the access module (roles + initial admin).
func AccessSeeders() []seeder.Seeder {
	return []seeder.Seeder{
		&RoleSeeder{},
		&AdminUserSeeder{},
	}
}
