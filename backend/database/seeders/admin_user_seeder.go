package seeders

import (
	"os"

	"reflexcms/backend/app/facades"
	accessmodels "reflexcms/backend/modules/access/models"
)

type AdminUserSeeder struct{}

func (s *AdminUserSeeder) Signature() string { return "admin-user-seeder" }

// Run creates the initial super-admin account exactly once. Credentials come
// from env with safe-for-dev defaults; change them before any real deployment.
func (s *AdminUserSeeder) Run() error {
	count, err := facades.Orm().Query().Model(&accessmodels.User{}).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var role accessmodels.Role
	if err := facades.Orm().Query().Where("name = ?", "super-admin").FirstOrFail(&role); err != nil {
		return err
	}

	hashed, hashErr := facades.Hash().Make(adminPassword())
	if hashErr != nil {
		return hashErr
	}

	admin := accessmodels.User{
		Username: adminEnv("ADMIN_USERNAME", "admin"),
		Email:    adminEnv("ADMIN_EMAIL", "admin@reflexcms.dev"),
		Password: hashed,
		RoleID:   role.ID,
		Status:   "active",
	}

	return facades.Orm().Query().Create(&admin)
}

func adminPassword() string {
	if p := os.Getenv("ADMIN_PASSWORD"); p != "" {
		return p
	}
	return "ReflexCMS@2026"
}

func adminEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
