package services

import (
	"errors"

	"reflexcms/backend/app/facades"
	accessmodels "reflexcms/backend/modules/access/models"
)

var ErrUserNotFound = errors.New("user not found")

// FindActiveUserByEmail loads a non-banned account together with its role.
// All values are bound parameters — never interpolated (plan §10.3).
func FindActiveUserByEmail(email string) (*accessmodels.User, error) {
	var user accessmodels.User
	// FirstOrFail: goravel's First silently zero-fills on empty results.
	if err := facades.Orm().Query().
		Where("email = ? AND status = ?", email, "active").
		FirstOrFail(&user); err != nil {
		return nil, ErrUserNotFound
	}

	if user.RoleID > 0 {
		var role accessmodels.Role
		if err := facades.Orm().Query().Where("id = ?", user.RoleID).FirstOrFail(&role); err == nil {
			user.Role = &role
		}
	}

	return &user, nil
}
