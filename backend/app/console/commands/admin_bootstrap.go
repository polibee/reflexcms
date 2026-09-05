// Package commands holds platform-level artisan commands. Module-owned
// commands (if any) live with their modules; these belong to the product
// binary itself.
package commands

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"

	"github.com/goravel/framework/contracts/console"
	consolecommand "github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"

	accessmodels "reflexcms/backend/modules/access/models"
)

// AdminBootstrap creates or resets the super-admin account. Intended for
// production bring-up: `./reflexcms artisan admin:bootstrap` prints a
// securely generated password once; `--password`/`--email` override for
// scripted setups. Never touches an existing super-admin unless flags are
// given, so re-running it during upgrades is a no-op.
type AdminBootstrap struct{}

func (c *AdminBootstrap) Signature() string { return "admin:bootstrap" }
func (c *AdminBootstrap) Description() string {
	return "Create or reset the super-admin account (generates a secure password when omitted)"
}

func (c *AdminBootstrap) Extend() consolecommand.Extend {
	return consolecommand.Extend{
		Category: "admin",
		Flags: []consolecommand.Flag{
			&consolecommand.StringFlag{Name: "email", Usage: "admin email (default: existing super-admin, else ADMIN_EMAIL, else admin@reflexcms.dev)"},
			&consolecommand.StringFlag{Name: "username", Usage: "admin username when creating (default: ADMIN_USERNAME, else admin)"},
			&consolecommand.StringFlag{Name: "password", Usage: "explicit password (min 8 chars); omit to generate a strong random one"},
		},
	}
}

func (c *AdminBootstrap) Handle(ctx console.Context) error {
	roleID, err := ensureSuperAdminRole()
	if err != nil {
		return err
	}

	// Resolve target: --email > existing super-admin > fresh account.
	email := strings.TrimSpace(ctx.Option("email"))
	var user accessmodels.User
	if email != "" {
		err = facades.Orm().Query().Where("email = ?", email).FirstOrFail(&user)
		if err != nil && !isRecordNotFound(err) {
			return err
		}
	} else {
		err = facades.Orm().Query().
			Where("role_id = ?", roleID).OrderBy("id", "asc").FirstOrFail(&user)
		if err != nil && !isRecordNotFound(err) {
			return err
		}
	}

	created := isRecordNotFound(err)
	if created {
		if email == "" {
			email = adminEnvOr("ADMIN_EMAIL", "admin@reflexcms.dev")
		}
		user = accessmodels.User{
			Username: strings.TrimSpace(ctx.Option("username")),
			Email:    email,
			Status:   "active",
		}
		if user.Username == "" {
			user.Username = adminEnvOr("ADMIN_USERNAME", "admin")
		}
	}

	// Password: explicit flag wins, else env, else generate. On an existing
	// account only reset when explicitly requested via flags/env.
	password := strings.TrimSpace(ctx.Option("password"))
	if created {
		if password == "" {
			password = os.Getenv("ADMIN_PASSWORD")
		}
	} else if password == "" {
		ctx.Info("Super-admin already exists: " + user.Email)
		ctx.Info("To reset its password, re-run with --password (or omit all flags and target a fresh install).")
		return nil
	}
	if password == "" {
		password, err = generatePassword(20)
		if err != nil {
			return err
		}
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	hashed, hashErr := facades.Hash().Make(password)
	if hashErr != nil {
		return hashErr
	}
	user.Password = hashed
	user.RoleID = roleID
	user.Status = "active"

	if created {
		if err := facades.Orm().Query().Create(&user); err != nil {
			return err
		}
	} else if err := facades.Orm().Query().Save(&user); err != nil {
		return err
	}

	ctx.NewLine()
	ctx.Info("✓ Super-admin ready")
	ctx.Line("  Email:    " + user.Email)
	ctx.Line("  Username: " + user.Username)
	if passwordIsGenerated := strings.TrimSpace(ctx.Option("password")) == "" && os.Getenv("ADMIN_PASSWORD") == ""; passwordIsGenerated {
		ctx.Line("  Password: " + password)
		ctx.NewLine()
		ctx.Info("Store it now — it is NOT written anywhere else.")
	} else {
		ctx.Line("  Password: (as provided)")
	}
	ctx.NewLine()
	return nil
}

// ensureSuperAdminRole returns the super-admin role id, creating the role
// when the deployment skipped db:seed.
func ensureSuperAdminRole() (uint64, error) {
	var role accessmodels.Role
	if err := facades.Orm().Query().Where("name = ?", "super-admin").FirstOrFail(&role); err == nil {
		return role.ID, nil
	}
	role = accessmodels.Role{
		Name:        "super-admin",
		DisplayName: "超级管理员",
		Permissions: accessmodels.StringList{"*"},
		Sort:        10,
	}
	if err := facades.Orm().Query().Create(&role); err != nil {
		return 0, err
	}
	return role.ID, nil
}

func isRecordNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "record not found")
}

// generatePassword returns n characters from an unambiguous alphabet —
// no 0/O/1/l/I — via crypto/rand.
func generatePassword(n int) (string, error) {
	const alphabet = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789!@#$%^&*"
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	out := make([]byte, n)
	for i, b := range buf {
		out[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(out), nil
}

func adminEnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
