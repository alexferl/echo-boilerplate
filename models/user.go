package models

import (
	"slices"
	"time"

	"github.com/alexferl/echo-boilerplate/util"
)

type Role int

const (
	UserRole Role = iota + 1
	AdminRole
	SuperRole
)

func (r Role) String() string {
	return [...]string{"user", "admin", "super"}[r-1]
}

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

type AdminResponse struct {
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
	user.AddRole(role)
	return user
}

func (u *User) SetPassword(s string) error {
	b, err := util.HashPassword([]byte(s))
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

func (u *User) AdminResponse() *AdminResponse {
	return &AdminResponse{
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
	return util.VerifyPassword([]byte(u.Password), []byte(s))
}

func (u *User) AddRole(role Role) {
	if !slices.Contains(u.Roles, role.String()) {
		u.Roles = append(u.Roles, role.String())
	}
}

func (u *User) Ban(id string) {
	u.IsBanned = true
	t := time.Now()
	u.BannedAt = &t
	u.BannedBy = &Ref{Id: id}
	u.UnbannedAt = nil
	u.UnbannedBy = nil
}

func (u *User) Unban(id string) {
	u.IsBanned = false
	t := time.Now()
	u.BannedAt = nil
	u.BannedBy = nil
	u.UnbannedAt = &t
	u.UnbannedBy = &Ref{Id: id}
}

func (u *User) Lock(id string) {
	u.IsLocked = true
	t := time.Now()
	u.LockedAt = &t
	u.LockedBy = &Ref{Id: id}
	u.UnlockedAt = nil
	u.UnlockedBy = nil
}

func (u *User) Unlock(id string) {
	u.IsLocked = false
	t := time.Now()
	u.LockedAt = nil
	u.LockedBy = nil
	u.UnlockedAt = &t
	u.UnlockedBy = &Ref{Id: id}
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
	return util.VerifyPassword([]byte(u.RefreshToken), []byte(s))
}

func (u *User) getTokens() ([]byte, []byte, error) {
	claims := map[string]any{
		"roles":     u.Roles,
		"is_banned": u.IsBanned,
		"is_locked": u.IsLocked,
	}
	access, refresh, err := util.GenerateTokens(u.Id, claims)
	if err != nil {
		return nil, nil, err
	}

	return access, refresh, nil
}

func (u *User) encryptRefreshToken(token []byte) error {
	b, err := util.HashPassword(token)
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
