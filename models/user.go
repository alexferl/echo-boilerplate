package models

import (
	"errors"
	"slices"
	"time"

	"github.com/alexferl/echo-boilerplate/util/jwt"
	"github.com/alexferl/echo-boilerplate/util/password"
)

type Role int8

const (
	UserRole Role = iota + 1
	AdminRole
	SuperRole
)

func (r Role) String() string {
	return [...]string{"user", "admin", "super"}[r-1]
}

var RolesMap = map[string]Role{
	"user":  UserRole,
	"admin": AdminRole,
	"super": SuperRole,
}

var (
	ErrAdminRoleRequired = errors.New("admin or greater role required")

	ErrBanSelf           = errors.New("cannot ban self")
	ErrBanExist          = errors.New("user already banned")
	ErrBanMorePrivileged = errors.New("cannot ban user with higher permissions")

	ErrUnbanSelf           = errors.New("cannot unban self")
	ErrUnbanExist          = errors.New("user already unbanned")
	ErrUnbanMorePrivileged = errors.New("cannot unban user with higher permissions")

	ErrLockSelf           = errors.New("cannot lock self")
	ErrLockExist          = errors.New("user already locked")
	ErrLockMorePrivileged = errors.New("cannot lock user with higher permissions")

	ErrUnlockSelf           = errors.New("cannot unlock self")
	ErrUnlockExist          = errors.New("user already unlocked")
	ErrUnlockMorePrivileged = errors.New("cannot unlock user with higher permissions")

	ErrRoleSelf = errors.New("cannot modify own roles")

	ErrRoleAddExist          = errors.New("user already has role")
	ErrRoleAddMorePrivileged = errors.New("cannot add a more privileged role")

	ErrRoleRemoveNotExist       = errors.New("user doesn't have role")
	ErrRoleRemoveMorePrivileged = errors.New("cannot remove a more privileged role")
)

type User struct {
	*Model        `bson:",inline"`
	BannedAt      *time.Time `json:"banned_at" bson:"banned_at"`
	BannedBy      any        `json:"banned_by" bson:"banned_by"`
	Bio           string     `json:"bio" bson:"bio"`
	Email         string     `json:"email" bson:"email"`
	IsBanned      bool       `json:"is_banned" bson:"is_banned"`
	IsLocked      bool       `json:"is_locked" bson:"is_locked"`
	LastLoginAt   *time.Time `json:"-" bson:"last_login_at"`
	LastLogoutAt  *time.Time `json:"-" bson:"last_logout_at"`
	LastRefreshAt *time.Time `json:"-" bson:"last_refresh_at"`
	LockedAt      *time.Time `json:"locked_at" bson:"locked_at"`
	LockedBy      any        `json:"locked_by" bson:"locked_by"`
	Name          string     `json:"name" bson:"name"`
	Password      string     `json:"-" bson:"password"`
	RefreshToken  string     `json:"-" bson:"refresh_token"`
	Roles         []string   `json:"-" bson:"roles"`
	UnbannedAt    *time.Time `json:"unbanned_at" bson:"unbanned_at"`
	UnbannedBy    any        `json:"unbanned_by" bson:"unbanned_by"`
	UnlockedAt    *time.Time `json:"unlocked_at" bson:"unlocked_at"`
	UnlockedBy    any        `json:"unlocked_by" bson:"unlocked_by"`
	Username      string     `json:"username" bson:"username"`
}

type Users []User

func (users Users) Response() *UsersResponse {
	res := make([]UserResponse, 0)
	for _, user := range users {
		res = append(res, *user.Response())
	}
	return &UsersResponse{Users: res}
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

type PublicUsersResponse struct {
	Users []UserPublic `json:"users"`
}

func (users Users) Public() *PublicUsersResponse {
	res := make([]UserPublic, 0)
	for _, user := range users {
		res = append(res, *user.Public())
	}
	return &PublicUsersResponse{Users: res}
}

type UserResponse struct {
	Id        string     `json:"id" bson:"id"`
	Bio       string     `json:"bio" bson:"bio"`
	Email     string     `json:"email" bson:"email"`
	Name      string     `json:"name" bson:"name"`
	Username  string     `json:"username" bson:"username"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
}

type UserAdminResponse struct {
	Id            string     `json:"id" bson:"id"`
	Bio           string     `json:"bio" bson:"bio"`
	Email         string     `json:"email" bson:"email"`
	IsBanned      bool       `json:"is_banned" bson:"is_banned"`
	IsLocked      bool       `json:"is_locked" bson:"is_locked"`
	Name          string     `json:"name" bson:"name"`
	Roles         []string   `json:"roles" bson:"roles"`
	Username      string     `json:"username" bson:"username"`
	CreatedAt     *time.Time `json:"created_at" bson:"created_at"`
	CreatedBy     *User      `json:"created_by" bson:"created_by"`
	DeletedAt     *time.Time `json:"-" bson:"deleted_at"`
	DeletedBy     string     `json:"-" bson:"deleted_by"`
	UpdatedAt     *time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy     *User      `json:"updated_by" bson:"updated_by"`
	LastLoginAt   *time.Time `json:"last_login_at" bson:"last_login_at"`
	LastLogoutAt  *time.Time `json:"last_logout_at" bson:"last_logout_at"`
	LastRefreshAt *time.Time `json:"last_refresh_at" bson:"last_refresh_at"`
	BannedAt      *time.Time `json:"banned_at" bson:"banned_at"`
	BannedBy      string     `json:"banned_by" bson:"banned_by"`
	UnbannedAt    *time.Time `json:"unbanned_at" bson:"unbanned_at"`
	UnbannedBy    string     `json:"unbanned_by" bson:"unbanned_by"`
	LockedAt      *time.Time `json:"locked_at" bson:"locked_at"`
	LockedBy      string     `json:"locked_by" bson:"locked_by"`
	UnlockedAt    *time.Time `json:"unlocked_at" bson:"unlocked_at"`
	UnlockedBy    string     `json:"unlocked_by" bson:"unlocked_by"`
}

type UserPublic struct {
	Id       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
}

func NewUser(email string, username string) *User {
	return &User{
		Model:    NewModel(),
		Email:    email,
		Username: username,
		Roles:    []string{UserRole.String()},
	}
}

func NewUserWithRole(email string, username string, role Role) *User {
	user := NewUser(email, username)
	user.addRole(role)
	return user
}

func (u *User) SetPassword(s string) error {
	b, err := password.Hash([]byte(s))
	if err != nil {
		return err
	}
	u.Password = b
	return nil
}

func (u *User) Response() *UserResponse {
	return &UserResponse{
		Id:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) AdminResponse() *UserAdminResponse {
	return &UserAdminResponse{
		Id:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		Bio:       u.Bio,
		IsBanned:  u.IsBanned,
		IsLocked:  u.IsLocked,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) Public() *UserPublic {
	return &UserPublic{
		Id:       u.Id,
		Username: u.Username,
		Name:     u.Name,
	}
}

func (u *User) ValidatePassword(s string) error {
	return password.Verify([]byte(u.Password), []byte(s))
}

func (u *User) HasRoleOrHigher(role Role) bool {
	if slices.Max(stringSliceToRolesSlice(u.Roles)) >= role {
		return true
	}
	return false
}

func stringSliceToRolesSlice(roles []string) []Role {
	var rs []Role
	for _, r := range roles {
		rs = append(rs, RolesMap[r])
	}
	return rs
}

// compare checks if u has a higher role than user
// and returns true if it does.
// It also returns true if both highest roles equals the SuperRole
// as users with the SuperRole are allowed to interact (ban, lock etc.)
// with each others.
func (u *User) compare(user *User) bool {
	roles := stringSliceToRolesSlice(u.Roles)
	otherRoles := stringSliceToRolesSlice(user.Roles)
	highestRole := slices.Max(roles)
	otherHighestRole := slices.Max(otherRoles)

	if highestRole == SuperRole && otherHighestRole == SuperRole {
		return false
	}

	if highestRole >= otherHighestRole {
		return true
	}

	return false
}

// hasRoleOrHigher check if user as at least role and returns true if it does.
func hasRoleOrHigher(user *User, role Role) bool {
	if slices.Max(stringSliceToRolesSlice(user.Roles)) >= role {
		return true
	}
	return false
}

func (u *User) addRole(role Role) {
	u.Roles = append(u.Roles, role.String())
}

func (u *User) AddRole(user *User, role Role) error {
	if slices.Max(stringSliceToRolesSlice(user.Roles)) < AdminRole {
		return NewError(ErrAdminRoleRequired, Permission)
	}

	if user.Id == u.Id {
		return NewError(ErrRoleSelf, Conflict)
	}

	if !hasRoleOrHigher(user, role) {
		return NewError(ErrRoleAddMorePrivileged, Permission)
	}

	if slices.Contains(u.Roles, role.String()) {
		return NewError(ErrRoleAddExist, Conflict)
	}

	u.addRole(role)

	return nil
}

func (u *User) removeRole(role Role) {
	idx := slices.Index(u.Roles, role.String())
	u.Roles = slices.Delete(u.Roles, idx, idx+1)
}

func (u *User) RemoveRole(user *User, role Role) error {
	if slices.Max(stringSliceToRolesSlice(user.Roles)) < AdminRole {
		return NewError(ErrAdminRoleRequired, Permission)
	}

	if user.Id == u.Id {
		return NewError(ErrRoleSelf, Conflict)
	}

	if !hasRoleOrHigher(user, role) {
		return NewError(ErrRoleRemoveMorePrivileged, Permission)
	}

	if !slices.Contains(u.Roles, role.String()) {
		return NewError(ErrRoleRemoveNotExist, Conflict)
	}

	u.removeRole(role)

	return nil
}

func (u *User) Ban(user *User) error {
	if slices.Max(stringSliceToRolesSlice(user.Roles)) < AdminRole {
		return NewError(ErrAdminRoleRequired, Permission)
	}

	if user.Id == u.Id {
		return NewError(ErrBanSelf, Conflict)
	}

	if u.compare(user) {
		return NewError(ErrBanMorePrivileged, Permission)
	}

	if u.IsBanned {
		return NewError(ErrBanExist, Conflict)
	}

	u.IsBanned = true
	t := time.Now()
	u.BannedAt = &t
	u.BannedBy = &Ref{Id: user.Id}
	u.UnbannedAt = nil
	u.UnbannedBy = nil

	return nil
}

func (u *User) Unban(user *User) error {
	if slices.Max(stringSliceToRolesSlice(user.Roles)) < AdminRole {
		return NewError(ErrAdminRoleRequired, Permission)
	}

	if user.Id == u.Id {
		return NewError(ErrUnbanSelf, Conflict)
	}

	if u.compare(user) {
		return NewError(ErrUnbanMorePrivileged, Permission)
	}

	if !u.IsBanned {
		return NewError(ErrUnbanExist, Conflict)
	}

	u.IsBanned = false
	t := time.Now()
	u.BannedAt = nil
	u.BannedBy = nil
	u.UnbannedAt = &t
	u.UnbannedBy = &Ref{Id: user.Id}

	return nil
}

func (u *User) Lock(user *User) error {
	if slices.Max(stringSliceToRolesSlice(user.Roles)) < AdminRole {
		return NewError(ErrAdminRoleRequired, Permission)
	}

	if user.Id == u.Id {
		return NewError(ErrLockSelf, Conflict)
	}

	if u.compare(user) {
		return NewError(ErrLockMorePrivileged, Permission)
	}

	if u.IsLocked {
		return NewError(ErrLockExist, Conflict)
	}

	u.IsLocked = true
	t := time.Now()
	u.LockedAt = &t
	u.LockedBy = &Ref{Id: user.Id}
	u.UnlockedAt = nil
	u.UnlockedBy = nil

	return nil
}

func (u *User) Unlock(user *User) error {
	if slices.Max(stringSliceToRolesSlice(user.Roles)) < AdminRole {
		return NewError(ErrAdminRoleRequired, Permission)
	}

	if user.Id == u.Id {
		return NewError(ErrUnlockSelf, Conflict)
	}

	if u.compare(user) {
		return NewError(ErrUnlockMorePrivileged, Permission)
	}

	if !u.IsLocked {
		return NewError(ErrUnlockExist, Conflict)
	}

	u.IsLocked = false
	t := time.Now()
	u.LockedAt = nil
	u.LockedBy = nil
	u.UnlockedAt = &t
	u.UnlockedBy = &Ref{Id: user.Id}

	return nil
}

func (u *User) Login() ([]byte, []byte, error) {
	access, refresh, err := u.getTokens()
	if err != nil {
		return nil, nil, err
	}

	t := time.Now()
	u.LastLoginAt = &t

	err = u.encryptRefreshToken(refresh)
	if err != nil {
		return nil, nil, err
	}

	return access, refresh, nil
}

func (u *User) Logout() {
	t := time.Now()
	u.LastLogoutAt = &t
	u.RefreshToken = ""
}

func (u *User) Refresh() ([]byte, []byte, error) {
	access, refresh, err := u.getTokens()
	if err != nil {
		return nil, nil, err
	}

	t := time.Now()
	u.LastRefreshAt = &t

	err = u.encryptRefreshToken(refresh)
	if err != nil {
		return nil, nil, err
	}

	return access, refresh, nil
}

func (u *User) ValidateRefreshToken(s string) error {
	return password.Verify([]byte(u.RefreshToken), []byte(s))
}

func (u *User) getTokens() ([]byte, []byte, error) {
	access, refresh, err := jwt.GenerateTokens(u.Id, nil)
	if err != nil {
		return nil, nil, err
	}

	return access, refresh, nil
}

func (u *User) encryptRefreshToken(token []byte) error {
	b, err := password.Hash(token)
	if err != nil {
		return err
	}

	u.RefreshToken = b

	return nil
}

type UserSearchParams struct {
	Limit int
	Skip  int
}
