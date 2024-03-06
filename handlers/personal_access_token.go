package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog/log"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/models"
)

type PersonalAccessTokenService interface {
	Create(ctx context.Context, model *models.PersonalAccessToken) (*models.PersonalAccessToken, error)
	Read(ctx context.Context, id string) (*models.PersonalAccessToken, error)
	Delete(ctx context.Context, model *models.PersonalAccessToken) error

	Find(ctx context.Context, userId string) (models.PersonalAccessTokens, error)
	FindOne(ctx context.Context, userId string, name string) (*models.PersonalAccessToken, error)
}

type PersonalAccessTokenHandler struct {
	*openapi.Handler
	svc PersonalAccessTokenService
}

func (h *PersonalAccessTokenHandler) Register(s *server.Server) {
	s.Add(http.MethodPost, "/personal_access_tokens", h.create)
	s.Add(http.MethodGet, "/personal_access_tokens", h.list)
	s.Add(http.MethodGet, "/personal_access_tokens/:id", h.get)
	s.Add(http.MethodDelete, "/personal_access_tokens/:id", h.revoke)
}

func NewPersonalAccessTokenHandler(openapi *openapi.Handler, svc PersonalAccessTokenService) *PersonalAccessTokenHandler {
	return &PersonalAccessTokenHandler{
		Handler: openapi,
		svc:     svc,
	}
}

type CreatePersonalAccessTokenRequest struct {
	Name      string `json:"name" bson:"name"`
	ExpiresAt string `json:"expires_at" bson:"expires_at"`
}

func (h *PersonalAccessTokenHandler) create(c echo.Context) error {
	token := c.Get("token").(jwt.Token)

	body := &CreatePersonalAccessTokenRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	res, err := h.svc.FindOne(ctx, token.Subject(), body.Name)
	if err != nil {
		if !errors.Is(err, data.ErrNoDocuments) {
			log.Error().Err(err).Msg("failed getting personal access token")
			return err
		}
	}

	if res != nil {
		return h.Validate(c, http.StatusConflict, echo.Map{"message": "token name already in-use"})
	}

	newPAT, err := models.NewPersonalAccessToken(token, body.Name, body.ExpiresAt)
	if err != nil {
		if errors.Is(err, models.ErrExpiresAtPast) {
			m := echo.Map{
				"message": "Validation error",
				"errors":  []string{models.ErrExpiresAtPast.Error()},
			}
			return h.Validate(c, http.StatusUnprocessableEntity, m)
		}
		log.Error().Err(err).Msg("failed generating personal access token")
		return err
	}

	decodedToken := newPAT.Token
	if err = newPAT.Encrypt(); err != nil {
		log.Error().Err(err).Msg("failed encrypting personal access token")
		return err
	}

	pat, err := h.svc.Create(ctx, newPAT)
	if err != nil {
		log.Error().Err(err).Msg("failed inserting personal access token")
		return err
	}

	pat.Token = decodedToken

	return h.Validate(c, http.StatusOK, pat.Response())
}

type ListPersonalAccessTokeResponse struct {
	Tokens []*models.PersonalAccessTokenResponse `json:"personal_access_tokens"`
}

func (h *PersonalAccessTokenHandler) list(c echo.Context) error {
	token := c.Get("token").(jwt.Token)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	pats, err := h.svc.Find(ctx, token.Subject())
	if err != nil {
		log.Error().Err(err).Msg("failed getting personal access token")
		return err
	}

	return h.Validate(c, http.StatusOK, ListPersonalAccessTokeResponse{Tokens: pats.Response()})
}

func (h *PersonalAccessTokenHandler) get(c echo.Context) error {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	pat, err := h.svc.Read(ctx, id)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return h.Validate(c, http.StatusNotFound, echo.Map{"message": "personal access token not found"})
		}
		log.Error().Err(err).Msg("failed getting user")
		return err
	}

	return h.Validate(c, http.StatusOK, pat.Response())
}

func (h *PersonalAccessTokenHandler) revoke(c echo.Context) error {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	pat, err := h.svc.Read(ctx, id)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return h.Validate(c, http.StatusNotFound, echo.Map{"message": "personal access token not found"})
		}
		log.Error().Err(err).Msg("failed getting user")
		return err
	}

	if pat.Revoked == true {
		return h.Validate(c, http.StatusConflict, echo.Map{"message": "personal access token already revoked"})
	}

	err = h.svc.Delete(ctx, pat)
	if err != nil {
		log.Error().Err(err).Msg("failed deleting personal access token")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}
