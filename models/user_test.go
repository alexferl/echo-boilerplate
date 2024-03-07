package models

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	email := "test@example.com"
	username := "text"
	id := "1"
	user := NewUserWithRole(email, username, AdminRole)
	assert.NotEqual(t, "", user.Id)

	user.AddRole(SuperRole)
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

	user.Ban(id)
	assert.Equal(t, id, user.BannedBy.(*Ref).Id)
	assert.NotNil(t, user.BannedAt)

	user.Unban(id)
	assert.Equal(t, id, user.UnbannedBy.(*Ref).Id)
	assert.NotNil(t, user.UnbannedAt)
	assert.Nil(t, user.BannedBy)
	assert.Nil(t, user.BannedAt)

	user.Lock(id)
	assert.Equal(t, id, user.LockedBy.(*Ref).Id)
	assert.NotNil(t, user.LockedAt)

	user.Unlock(id)
	assert.Equal(t, id, user.UnlockedBy.(*Ref).Id)
	assert.NotNil(t, user.UnlockedAt)
	assert.Nil(t, user.LockedBy)
	assert.Nil(t, user.LockedAt)

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
}
