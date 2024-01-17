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
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/util"
)

const PATCollection = "personal_access_tokens"

var ErrExpiresAtPast = errors.New("expires_at cannot be in the past")

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
	b, err := util.HashPassword([]byte(pat.Token))
	if err != nil {
		return err
	}

	pat.Token = b

	return nil
}

func (pat *PersonalAccessToken) Validate(s string) error {
	return util.VerifyPassword([]byte(pat.Token), []byte(s))
}

func (pat *PersonalAccessToken) Response() *PATResponse {
	return &PATResponse{
		Id:        pat.Id,
		Href:      util.GetFullURL(fmt.Sprintf("/user/personal_access_tokens/%s", pat.Id)),
		Name:      pat.Name,
		Revoked:   pat.Revoked,
		UserId:    pat.UserId,
		CreatedAt: pat.CreatedAt,
		ExpiresAt: pat.ExpiresAt,
	}
}

type PersonalAccessTokens []PersonalAccessToken

func (pats PersonalAccessTokens) Response() []*PATResponse {
	res := make([]*PATResponse, 0)
	for _, pat := range pats {
		res = append(res, pat.Response())
	}
	return res
}

type PATResponse struct {
	Id        string     `json:"id" bson:"id"`
	Href      string     `json:"href" bson:"href"`
	Name      string     `json:"name" bson:"name"`
	Revoked   bool       `json:"revoked" bson:"revoked"`
	UserId    string     `json:"user_id" bson:"user_id"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	ExpiresAt *time.Time `json:"expires_at" bson:"expires_at"`
}

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

type CreatePATRequest struct {
	Name      string `json:"name" bson:"name"`
	ExpiresAt string `json:"expires_at" bson:"expires_at"`
}

func (h *Handler) CreatePersonalAccessToken(c echo.Context) error {
	logger := c.Get("logger").(zerolog.Logger)
	token := c.Get("token").(jwt.Token)

	body := &CreatePATRequest{}
	if err := c.Bind(body); err != nil {
		logger.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"user_id", token.Subject()}, {"name", body.Name}}
	result, err := h.Mapper.WithCollection(PATCollection).FindOne(ctx, filter, &PersonalAccessToken{})
	if err != nil {
		if !errors.Is(err, data.ErrNoDocuments) {
			logger.Error().Err(err).Msg("failed getting personal access token")
			return err
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
		logger.Error().Err(err).Msg("failed generating personal access token")
		return err
	}

	decodedToken := newPAT.Token
	if err = newPAT.Encrypt(); err != nil {
		logger.Error().Err(err).Msg("failed encrypting personal access token")
		return err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true)
	upsert, err := h.Mapper.WithCollection(PATCollection).FindOneAndUpdate(ctx, filter, newPAT, &PersonalAccessToken{}, opts)
	if err != nil {
		logger.Error().Err(err).Msg("failed inserting personal access token")
		return err
	}

	pat := upsert.(*PersonalAccessToken)
	pat.Token = decodedToken

	return h.Validate(c, http.StatusOK, pat.Response())
}

type ListPATResponse struct {
	Tokens []*PATResponse `json:"personal_access_tokens"`
}

func (h *Handler) ListPersonalAccessTokens(c echo.Context) error {
	logger := c.Get("logger").(zerolog.Logger)
	token := c.Get("token").(jwt.Token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"user_id", token.Subject()}}
	result, err := h.Mapper.WithCollection(PATCollection).Find(ctx, filter, PersonalAccessTokens{})
	if err != nil {
		logger.Error().Err(err).Msg("failed getting personal access token")
		return err
	}

	return h.Validate(c, http.StatusOK, ListPATResponse{Tokens: result.(PersonalAccessTokens).Response()})
}

func (h *Handler) GetPersonalAccessToken(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pat, errResp := h.getPAT(ctx, c)
	if errResp != nil {
		return errResp()
	}

	return h.Validate(c, http.StatusOK, pat)
}

func (h *Handler) RevokePersonalAccessToken(c echo.Context) error {
	logger := c.Get("logger").(zerolog.Logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pat, errResp := h.getPAT(ctx, c)
	if errResp != nil {
		return errResp()
	}

	if pat.Revoked == true {
		return h.Validate(c, http.StatusConflict, echo.Map{"message": "personal access token already revoked"})
	}

	pat.Revoked = true
	_, err := h.Mapper.WithCollection(PATCollection).UpdateOneById(ctx, c.Param("id"), pat, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed deleting personal access token")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}

func (h *Handler) getPAT(ctx context.Context, c echo.Context) (*PATResponse, func() error) {
	id := c.Param("id")
	logger := c.Get("logger").(zerolog.Logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"$or", bson.A{
		bson.D{{"id", id}},
		bson.D{{"name", id}},
	}}}
	result, err := h.Mapper.WithCollection(PATCollection).FindOne(ctx, filter, &PersonalAccessToken{})
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return nil, util.WrapErr(h.Validate(c, http.StatusNotFound, echo.Map{"message": "personal access token not found"}))
		}
		logger.Error().Err(err).Msg("failed getting personal access token")
		return nil, util.WrapErr(err)
	}

	pat := result.(*PersonalAccessToken).Response()

	return pat, nil
}
