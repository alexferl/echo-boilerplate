package auth

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog/log"
)

func (h *Handler) AuthToken(c echo.Context) error {
	token := c.Get("token").(jwt.Token)
	m, err := token.AsMap(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed converting token to map")
		return err
	}

	return h.Validate(c, http.StatusOK, m)
}
