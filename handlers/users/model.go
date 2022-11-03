package users

import (
	"encoding/base64"
	"time"

	"github.com/minio/sha256-simd"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"

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
	*data.Model      `bson:",inline"`
	Email            string     `json:"email" bson:"email"`
	Username         string     `json:"username" bson:"username"`
	Password         string     `json:"-" bson:"password"`
	Name             string     `json:"name" bson:"name"`
	Bio              string     `json:"bio" bson:"bio"`
	Roles            []string   `json:"-" bson:"roles"`
	RefreshToken     string     `json:"-" bson:"refresh_token"`
	RefreshTokenHash string     `json:"-" bson:"refresh_token_hash"`
	LastLoginAt      *time.Time `json:"-" bson:"last_login_at"`
	LastLogoutAt     *time.Time `json:"-" bson:"last_logout_at"`
	LastRefreshAt    *time.Time `json:"-" bson:"last_refresh_at"`
}

func NewUser(email string, username string) *User {
	return &User{
		Model:    data.NewModel(),
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
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(b)

	return nil
}

func (u *User) CheckPassword(s string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(s))
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

	h := sha256.New()
	h.Write(refresh)
	b := h.Sum(nil)

	u.RefreshToken = string(refresh)
	u.RefreshTokenHash = base64.StdEncoding.EncodeToString(b)

	return access, refresh, nil
}

func (u *User) Logout() {
	t := time.Now()
	u.LastLogoutAt = &t

	u.RefreshToken = ""
	u.RefreshTokenHash = ""
}

func (u *User) Refresh() ([]byte, []byte, error) {
	access, refresh, err := util.GenerateTokens(u.Id, map[string]any{"roles": u.Roles})
	if err != nil {
		return nil, nil, err
	}

	t := time.Now()
	u.LastRefreshAt = &t

	h := sha256.New()
	h.Write(refresh)
	b := h.Sum(nil)

	u.RefreshToken = string(refresh)
	u.RefreshTokenHash = base64.StdEncoding.EncodeToString(b)

	return access, refresh, nil
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}
