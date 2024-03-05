package users

import (
	"fmt"
	"slices"
	"time"

	"github.com/alexferl/echo-boilerplate/data"
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
	*data.Model   `bson:",inline"`
	Bio           string     `json:"bio" bson:"bio"`
	Email         string     `json:"email" bson:"email"`
	IsBanned      bool       `json:"is_banned" bson:"is_banned"`
	IsLocked      bool       `json:"is_locked" bson:"is_locked"`
	Name          string     `json:"name" bson:"name"`
	Password      string     `json:"-" bson:"password"`
	RefreshToken  string     `json:"-" bson:"refresh_token"`
	Roles         []string   `json:"-" bson:"roles"`
	Username      string     `json:"username" bson:"username"`
	LastLoginAt   *time.Time `json:"-" bson:"last_login_at"`
	LastLogoutAt  *time.Time `json:"-" bson:"last_logout_at"`
	LastRefreshAt *time.Time `json:"-" bson:"last_refresh_at"`
	BannedAt      *time.Time `json:"banned_at" bson:"banned_at"`
	BannedBy      string     `json:"banned_by" bson:"banned_by"`
	UnbannedAt    *time.Time `json:"unbanned_at" bson:"unbanned_at"`
	UnbannedBy    string     `json:"unbanned_by" bson:"unbanned_by"`
	LockedAt      *time.Time `json:"locked_at" bson:"locked_at"`
	LockedBy      string     `json:"locked_by" bson:"locked_by"`
	UnlockedAt    *time.Time `json:"unlocked_at" bson:"unlocked_at"`
	UnlockedBy    string     `json:"unlocked_by" bson:"unlocked_by"`
}

type Users []User

func (users Users) Response() []*Response {
	res := make([]*Response, 0)
	for _, user := range users {
		res = append(res, user.Response())
	}
	return res
}

func (users Users) Public() []*Public {
	res := make([]*Public, 0)
	for _, user := range users {
		res = append(res, user.Public())
	}
	return res
}

type Response struct {
	Id        string     `json:"id" bson:"id"`
	Href      string     `json:"href" bson:"href"`
	Bio       string     `json:"bio" bson:"bio"`
	Email     string     `json:"email" bson:"email"`
	Name      string     `json:"name" bson:"name"`
	Username  string     `json:"username" bson:"username"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
}

type AdminResponse struct {
	Id            string     `json:"id" bson:"id"`
	Href          string     `json:"href" bson:"href"`
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

type Public struct {
	Id       string `json:"id" bson:"id"`
	Href     string `json:"href" bson:"href"`
	Username string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
	Bio      string `json:"bio" bson:"bio"`
}

func NewUser(email string, username string) *User {
	return &User{
		Model:    data.NewModel(""),
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

func (u *User) Response() *Response {
	return &Response{
		Id:        u.Id,
		Href:      util.GetFullURL(fmt.Sprintf("/users/%s", u.Id)),
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
		Href:      util.GetFullURL(fmt.Sprintf("/users/%s", u.Id)),
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

func (u *User) Public() *Public {
	return &Public{
		Id:       u.Id,
		Href:     util.GetFullURL(fmt.Sprintf("/users/%s", u.Id)),
		Username: u.Username,
		Name:     u.Name,
		Bio:      u.Bio,
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
	u.BannedBy = id
	u.UnbannedAt = nil
	u.UnbannedBy = ""
}

func (u *User) Unban(id string) {
	u.IsBanned = false
	t := time.Now()
	u.BannedAt = nil
	u.BannedBy = ""
	u.UnbannedAt = &t
	u.UnbannedBy = id
}

func (u *User) Lock(id string) {
	u.IsLocked = true
	t := time.Now()
	u.LockedAt = &t
	u.LockedBy = id
	u.UnlockedAt = nil
	u.UnlockedBy = ""
}

func (u *User) Unlock(id string) {
	u.IsLocked = false
	t := time.Now()
	u.LockedAt = nil
	u.LockedBy = ""
	u.UnlockedAt = &t
	u.UnlockedBy = id
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
