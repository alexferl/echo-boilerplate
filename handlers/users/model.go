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
)

func (r Role) String() string {
	return [...]string{"user", "admin"}[r-1]
}

type User struct {
	*data.Model   `bson:",inline"`
	Email         string     `json:"email" bson:"email"`
	Username      string     `json:"username" bson:"username"`
	Password      string     `json:"-" bson:"password"`
	Name          string     `json:"name" bson:"name"`
	Bio           string     `json:"bio" bson:"bio"`
	Roles         []string   `json:"-" bson:"roles"`
	RefreshToken  string     `json:"-" bson:"refresh_token"`
	LastLoginAt   *time.Time `json:"-" bson:"last_login_at"`
	LastLogoutAt  *time.Time `json:"-" bson:"last_logout_at"`
	LastRefreshAt *time.Time `json:"-" bson:"last_refresh_at"`
}

type Users []User

func (users Users) Response() []*UserResponse {
	res := make([]*UserResponse, 0)
	for _, user := range users {
		res = append(res, user.Response())
	}

	return res
}

func (users Users) Public() []*UserResponsePublic {
	res := make([]*UserResponsePublic, 0)
	for _, user := range users {
		res = append(res, user.Public())
	}

	return res
}

type UserResponse struct {
	Id        string     `json:"id" bson:"id"`
	Href      string     `json:"href" bson:"href"`
	Username  string     `json:"username" bson:"username"`
	Email     string     `json:"email" bson:"email"`
	Name      string     `json:"name" bson:"name"`
	Bio       string     `json:"bio" bson:"bio"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
}

type UserResponsePublic struct {
	Id       string `json:"id" bson:"id"`
	Href     string `json:"href" bson:"href"`
	Username string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
}

func NewUser(email string, username string) *User {
	return &User{
		Model:    data.NewModel(""),
		Email:    email,
		Username: username,
		Roles:    []string{UserRole.String()},
	}
}

func NewAdminUser(email string, username string) *User {
	user := NewUser(email, username)
	user.AddRole(AdminRole)

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
		Href:      util.GetFullURL(fmt.Sprintf("/users/%s", u.Id)),
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) Public() *UserResponsePublic {
	return &UserResponsePublic{
		Id:       u.Id,
		Href:     util.GetFullURL(fmt.Sprintf("/users/%s", u.Id)),
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

func (u *User) Login() ([]byte, []byte, error) {
	access, refresh, err := util.GenerateTokens(u.Id, map[string]any{"roles": u.Roles})
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
	access, refresh, err := util.GenerateTokens(u.Id, map[string]any{"roles": u.Roles})
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

func (u *User) encryptRefreshToken(token []byte) error {
	b, err := util.HashPassword(token)
	if err != nil {
		return err
	}

	u.RefreshToken = b

	return nil
}
