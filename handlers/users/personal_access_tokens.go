package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexferl/echo-boilerplate/util"
)

const PATCollection = "personal_access_tokens"

type PersonalAccessToken struct {
	Id        string     `json:"id" bson:"id"`
	Name      string     `json:"name" bson:"name"`
	Revoked   bool       `json:"revoked" bson:"revoked"`
	UserId    string     `json:"user_id" bson:"user_id"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	ExpiresAt *time.Time `json:"expires_at" bson:"expires_at"`
	Token     string     `json:"token" bson:"token"`
}

func (pat *PersonalAccessToken) Encrypt() error {
	b, err := util.HashPassword(pat.Token)
	if err != nil {
		return err
	}

	pat.Token = b

	return nil
}

func (pat *PersonalAccessToken) Validate(s string) error {
	return util.VerifyPassword(pat.Token, s)
}

type PATWithoutToken struct {
	Id        string     `json:"id" bson:"id"`
	Name      string     `json:"name" bson:"name"`
	Revoked   bool       `json:"revoked" bson:"revoked"`
	UserId    string     `json:"user_id" bson:"user_id"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	ExpiresAt *time.Time `json:"expires_at" bson:"expires_at"`
}

var ErrExpiresAtPast = errors.New("expires_at cannot be in the past")

func NewPersonalAccessToken(token jwt.Token, name string, expiresAt string) (*PersonalAccessToken, error) {
	t, err := time.Parse("2006-01-02", expiresAt)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if t.Before(now) {
		return nil, ErrExpiresAtPast
	}

	roles := util.GetRoles(token)
	pat, err := util.GeneratePersonalToken(token.Subject(), t.Sub(now), map[string]any{"roles": roles})
	if err != nil {
		return nil, err
	}

	return &PersonalAccessToken{
		Id:        xid.New().String(),
		Name:      name,
		UserId:    token.Subject(),
		Token:     string(pat),
		CreatedAt: &now,
		ExpiresAt: &t,
	}, nil
}

func (pat *PersonalAccessToken) MakeResponse() *PATWithoutToken {
	return &PATWithoutToken{
		Id:        pat.Id,
		Name:      pat.Name,
		Revoked:   pat.Revoked,
		UserId:    pat.UserId,
		CreatedAt: pat.CreatedAt,
		ExpiresAt: pat.ExpiresAt,
	}
}

type CreatePATRequest struct {
	Name      string `json:"name" bson:"name"`
	ExpiresAt string `json:"expires_at" bson:"expires_at"`
}

func (h *Handler) CreatePersonalAccessToken(c echo.Context) error {
	body := &CreatePATRequest{}
	if err := c.Bind(body); err != nil {
		return err
	}

	token := c.Get("token").(jwt.Token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"user_id", token.Subject()}, {"name", body.Name}}
	result, err := h.Mapper.Collection(PATCollection).FindOne(ctx, filter, &PersonalAccessToken{})
	if err != nil {
		if !errors.Is(err, ErrNoDocuments) {
			return fmt.Errorf("failed getting personal access token: %v", err)
		}
	}

	if result != nil {
		return h.Validate(c, http.StatusConflict, echo.Map{"message": "token name already in-use"})
	}

	newPAT, err := NewPersonalAccessToken(token, body.Name, body.ExpiresAt)
	if err != nil {
		if errors.Is(err, ErrExpiresAtPast) {
			m := echo.Map{
				"message": "Validation error",
				"errors":  []string{ErrExpiresAtPast.Error()},
			}
			return h.Validate(c, http.StatusUnprocessableEntity, m)
		}
		return fmt.Errorf("failed generating personal access token: %v", err)
	}

	decodedToken := newPAT.Token
	if err = newPAT.Encrypt(); err != nil {
		return fmt.Errorf("failed encrypting personal access token: %v", err)
	}

	opts := options.FindOneAndUpdate().SetUpsert(true)
	upsert, err := h.Mapper.Collection(PATCollection).Upsert(ctx, filter, newPAT, &PersonalAccessToken{}, opts)
	if err != nil {
		return fmt.Errorf("failed inserting personal access token: %v", err)
	}

	pat := upsert.(*PersonalAccessToken)
	pat.Token = decodedToken

	return h.Validate(c, http.StatusOK, pat)
}

type ListPATResponse struct {
	Tokens []*PATWithoutToken `json:"personal_access_tokens"`
}

func (h *Handler) ListPersonalAccessTokens(c echo.Context) error {
	token := c.Get("token").(jwt.Token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"user_id", token.Subject()}}
	result, err := h.Mapper.Collection(PATCollection).Find(ctx, filter, []*PATWithoutToken{})
	if err != nil {
		return fmt.Errorf("failed getting personal access token: %v", err)
	}

	return h.Validate(c, http.StatusOK, ListPATResponse{Tokens: result.([]*PATWithoutToken)})
}

func (h *Handler) GetPersonalAccessToken(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pat, errResp := h.getToken(ctx, c)
	if errResp != nil {
		return errResp()
	}

	return h.Validate(c, http.StatusOK, pat)
}

func (h *Handler) RevokePersonalAccessToken(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pat, errResp := h.getToken(ctx, c)
	if errResp != nil {
		return errResp()
	}

	pat.Revoked = true
	_, err := h.Mapper.Collection(PATCollection).UpdateById(ctx, c.Param("id"), pat, nil)
	if err != nil {
		return fmt.Errorf("failed inserting personal access token: %v", err)
	}

	return h.Validate(c, http.StatusNoContent, nil)
}

func (h *Handler) getToken(ctx context.Context, c echo.Context) (*PATWithoutToken, func() error) {
	taskId := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"id", taskId}}
	result, err := h.Mapper.Collection(PATCollection).FindOne(ctx, filter, &PATWithoutToken{})
	if err != nil {
		if errors.Is(err, ErrNoDocuments) {
			return nil, wrap(h.Validate(c, http.StatusNotFound, echo.Map{"message": "personal access token not found"}))
		}
		return nil, wrap(fmt.Errorf("failed getting personal access token: %v", err))
	}

	pat := result.(*PATWithoutToken)

	return pat, nil
}

func wrap(err error) func() error {
	return func() error { return err }
}
