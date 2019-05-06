package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/admiralobvious/echo-boilerplate/util"
)

type (
	// User represents the structure of our resource
	User struct {
		Id       bson.ObjectId `json:"id" bson:"_id"`
		Email    string        `json:"email" bson:"email"`
		Name     string        `json:"name,omitempty" bson:"name,omitempty"`
		Password string        `json:"-" bson:"password"`
		Type     string        `json:"type" bson:"type"`
		Username string        `json:"username,omitempty" bson:"username,omitempty"`
	}
)

// NewUser creates a new User
func NewUser(email, password string) *User {
	p, _ := util.HashPassword(password)
	return &User{
		Id:       bson.NewObjectId(),
		Email:    email,
		Password: p,
		Type:     "user",
	}
}

// GenerateJWT generates a JSON web token
func (u *User) GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = u.Email
	claims["name"] = u.Name
	claims["type"] = u.Type
	claims["username"] = u.Username
	claims["exp"] = time.Now().Add(time.Hour * 168).Unix()

	t, err := token.SignedString([]byte(viper.GetString("token-secret")))
	if err != nil {
		logrus.Panicf("Error generating JWT: %v", err)
		return "", err
	}

	return t, nil
}
