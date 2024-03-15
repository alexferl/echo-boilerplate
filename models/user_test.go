package models

import (
	"errors"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	email := "test@example.com"
	username := "test"
	id := "1"
	user := NewUserWithRole(email, username, AdminRole)
	assert.NotEqual(t, "", user.Id)
	user.Id = id

	hasUser := user.HasRoleOrHigher(UserRole)
	assert.True(t, hasUser)
	hasAdmin := user.HasRoleOrHigher(AdminRole)
	assert.True(t, hasAdmin)
	hasSuper := user.HasRoleOrHigher(SuperRole)
	assert.False(t, hasSuper)

	user.addRole(SuperRole)
	assert.True(t, slices.Contains(user.Roles, SuperRole.String()))

	pwd := "abcdefghijlk"
	err := user.SetPassword(pwd)
	assert.NoError(t, err)

	err = user.ValidatePassword(pwd)
	assert.NoError(t, err)

	user.Create(id)
	assert.Equal(t, id, user.CreatedBy.(*Ref).Id)
	assert.NotNil(t, user.CreatedAt)

	user.Update(id)
	assert.Equal(t, id, user.UpdatedBy.(*Ref).Id)
	assert.NotNil(t, user.UpdatedAt)

	user.Delete(id)
	assert.Equal(t, id, user.DeletedBy.(*Ref).Id)
	assert.NotNil(t, user.DeletedAt)

	access, refresh, err := user.Login()
	assert.NoError(t, err)
	assert.NotEqual(t, "", string(access))
	assert.NotEqual(t, "", string(refresh))
	assert.NotNil(t, user.LastLoginAt)

	err = user.ValidateRefreshToken(string(refresh))
	assert.NoError(t, err)

	access, refresh, err = user.Refresh()
	assert.NoError(t, err)
	assert.NotEqual(t, "", string(access))
	assert.NotEqual(t, "", string(refresh))
	assert.NotNil(t, user.LastRefreshAt)

	user.Logout()
	assert.Equal(t, "", user.RefreshToken)
	assert.NotNil(t, user.LastLogoutAt)

	resp := user.Response()
	assert.Equal(t, id, resp.Id)
	assert.Equal(t, email, resp.Email)
	assert.Equal(t, username, resp.Username)

	admin := user.AdminResponse()
	assert.False(t, admin.IsLocked)
}

func TestUsers(t *testing.T) {
	user1 := NewUser("test1@example.com", "test1")
	user2 := NewUser("test2@example.com", "test2")

	users := Users{*user1, *user2}
	resp := users.Response()

	assert.Len(t, resp.Users, 2)
}

func TestBan(t *testing.T) {
	user := NewUser("test@example.com", "test")
	bannedUser := NewUser("banned@example.com", "banned")
	bannedUser.IsBanned = true
	admin := NewUserWithRole("admin@example.com", "admin", AdminRole)
	admin1 := NewUserWithRole("admin1@example.com", "admin1", AdminRole)
	super := NewUserWithRole("super@example.com", "super", SuperRole)
	super1 := NewUserWithRole("super1@example.com", "super1", SuperRole)

	testCases := []struct {
		name   string
		user   *User
		target *User
		err    error
		kind   Kind
	}{
		{"need admin or higher role", user, user, ErrAdminRoleRequired, Permission},
		{"cannot ban self", admin, admin, ErrBanSelf, Conflict},
		{"target already banned", admin, bannedUser, ErrBanExist, Conflict},
		{"target cannot be more privileged", admin, super, ErrBanMorePrivileged, Permission},
		{"admin cannot ban admin", admin, admin1, ErrBanMorePrivileged, Permission},
		{"super can ban super", super, super1, nil, 0},
		{"success", admin, user, nil, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.target.Ban(tc.user)
			if tc.err != nil {
				assert.Error(t, err)
				var e *Error
				assert.ErrorAs(t, err, &e)
				if errors.As(err, &e) {
					assert.Equal(t, tc.err.Error(), e.Message)
					assert.Equal(t, tc.kind, e.Kind)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tc.target.BannedAt)
				assert.Equal(t, tc.user.Id, tc.target.BannedBy.(*Ref).Id)
			}
		})
	}
}

func TestUnban(t *testing.T) {
	user := NewUser("test@example.com", "test")
	user.IsBanned = true
	unbannedUser := NewUser("unbanned@example.com", "unbanned")
	admin := NewUserWithRole("admin@example.com", "admin", AdminRole)
	admin1 := NewUserWithRole("admin1@example.com", "admin1", AdminRole)
	super := NewUserWithRole("super@example.com", "super", SuperRole)
	super1 := NewUserWithRole("super1@example.com", "super1", SuperRole)
	super1.IsBanned = true

	testCases := []struct {
		name   string
		user   *User
		target *User
		err    error
		kind   Kind
	}{
		{"need admin or higher role", user, user, ErrAdminRoleRequired, Permission},
		{"cannot unban self", admin, admin, ErrUnbanSelf, Conflict},
		{"target already unbanned", admin, unbannedUser, ErrUnbanExist, Conflict},
		{"target cannot be more privileged", admin, super, ErrUnbanMorePrivileged, Permission},
		{"admin cannot unban admin", admin, admin1, ErrUnbanMorePrivileged, Permission},
		{"super can unban super", super, super1, nil, 0},
		{"success", admin, user, nil, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.target.Unban(tc.user)
			if tc.err != nil {
				assert.Error(t, err)
				var e *Error
				assert.ErrorAs(t, err, &e)
				if errors.As(err, &e) {
					assert.Equal(t, tc.err.Error(), e.Message)
					assert.Equal(t, tc.kind, e.Kind)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tc.target.UnbannedAt)
				assert.Equal(t, tc.user.Id, tc.target.UnbannedBy.(*Ref).Id)
			}
		})
	}
}

func TestLock(t *testing.T) {
	user := NewUser("test@example.com", "test")
	lockedUser := NewUser("locked@example.com", "locked")
	lockedUser.IsLocked = true
	admin := NewUserWithRole("admin@example.com", "admin", AdminRole)
	admin1 := NewUserWithRole("admin1@example.com", "admin1", AdminRole)
	super := NewUserWithRole("super@example.com", "super", SuperRole)
	super1 := NewUserWithRole("super1@example.com", "super1", SuperRole)

	testCases := []struct {
		name   string
		user   *User
		target *User
		err    error
		kind   Kind
	}{
		{"need admin or higher role", user, user, ErrAdminRoleRequired, Permission},
		{"cannot lock self", admin, admin, ErrLockSelf, Conflict},
		{"target already locked", admin, lockedUser, ErrLockExist, Conflict},
		{"target cannot be more privileged", admin, super, ErrLockMorePrivileged, Permission},
		{"admin cannot lock admin", admin, admin1, ErrLockMorePrivileged, Permission},
		{"super can lock super", super, super1, nil, 0},
		{"success", admin, user, nil, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.target.Lock(tc.user)
			if tc.err != nil {
				assert.Error(t, err)
				var e *Error
				assert.ErrorAs(t, err, &e)
				if errors.As(err, &e) {
					assert.Equal(t, tc.err.Error(), e.Message)
					assert.Equal(t, tc.kind, e.Kind)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tc.target.LockedAt)
				assert.Equal(t, tc.user.Id, tc.target.LockedBy.(*Ref).Id)
			}
		})
	}
}

func TestUnlock(t *testing.T) {
	user := NewUser("test@example.com", "test")
	user.IsLocked = true
	unlockedUser := NewUser("unlocked@example.com", "unlocked")
	admin := NewUserWithRole("admin@example.com", "admin", AdminRole)
	admin1 := NewUserWithRole("admin1@example.com", "admin1", AdminRole)
	super := NewUserWithRole("super@example.com", "super", SuperRole)
	super1 := NewUserWithRole("super1@example.com", "super1", SuperRole)
	super1.IsLocked = true

	testCases := []struct {
		name   string
		user   *User
		target *User
		err    error
		kind   Kind
	}{
		{"need admin or higher role", user, user, ErrAdminRoleRequired, Permission},
		{"cannot unlock self", admin, admin, ErrUnlockSelf, Conflict},
		{"target already unlocked", admin, unlockedUser, ErrUnlockExist, Conflict},
		{"target cannot be more privileged", admin, super, ErrUnlockMorePrivileged, Permission},
		{"admin cannot unlock admin", admin, admin1, ErrUnlockMorePrivileged, Permission},
		{"super can unlock super", super, super1, nil, 0},
		{"success", admin, user, nil, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.target.Unlock(tc.user)
			if tc.err != nil {
				assert.Error(t, err)
				var e *Error
				assert.ErrorAs(t, err, &e)
				if errors.As(err, &e) {
					assert.Equal(t, tc.err.Error(), e.Message)
					assert.Equal(t, tc.kind, e.Kind)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tc.target.UnlockedAt)
				assert.Equal(t, tc.user.Id, tc.target.UnlockedBy.(*Ref).Id)
			}
		})
	}
}

func TestAddRole(t *testing.T) {
	user := NewUser("test@example.com", "test")
	admin := NewUserWithRole("admin@example.com", "admin", AdminRole)

	testCases := []struct {
		name   string
		user   *User
		target *User
		role   Role
		err    error
		kind   Kind
	}{
		{"need admin or higher role", user, user, UserRole, ErrAdminRoleRequired, Permission},
		{"cannot modify self", admin, admin, AdminRole, ErrRoleSelf, Conflict},
		{"cannot add more privileged role", admin, user, SuperRole, ErrRoleAddMorePrivileged, Permission},
		{"cannot add existing role", admin, user, UserRole, ErrRoleAddExist, Conflict},
		{"success", admin, user, AdminRole, nil, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.target.AddRole(tc.user, tc.role)
			if tc.err != nil {
				assert.Error(t, err)
				var e *Error
				assert.ErrorAs(t, err, &e)
				if errors.As(err, &e) {
					assert.Equal(t, tc.err.Error(), e.Message)
					assert.Equal(t, tc.kind, e.Kind)
				}
			} else {
				assert.NoError(t, err)
				assert.Contains(t, tc.target.Roles, tc.role.String())
			}
		})
	}
}

func TestRemoveRole(t *testing.T) {
	user := NewUser("test@example.com", "test")
	admin := NewUserWithRole("admin@example.com", "admin", AdminRole)

	testCases := []struct {
		name   string
		user   *User
		target *User
		role   Role
		err    error
		kind   Kind
	}{
		{"need admin or higher role", user, user, UserRole, ErrAdminRoleRequired, Permission},
		{"cannot modify self", admin, admin, AdminRole, ErrRoleSelf, Conflict},
		{"cannot remove more privileged role", admin, user, SuperRole, ErrRoleRemoveMorePrivileged, Permission},
		{"cannot remove non existing role", admin, user, AdminRole, ErrRoleRemoveNotExist, Conflict},
		{"success", admin, user, UserRole, nil, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.target.RemoveRole(tc.user, tc.role)
			if tc.err != nil {
				assert.Error(t, err)
				var e *Error
				assert.ErrorAs(t, err, &e)
				if errors.As(err, &e) {
					assert.Equal(t, tc.err.Error(), e.Message)
					assert.Equal(t, tc.kind, e.Kind)
				}
			} else {
				assert.NoError(t, err)
				assert.NotContains(t, tc.target.Roles, tc.role.String())
			}
		})
	}
}
