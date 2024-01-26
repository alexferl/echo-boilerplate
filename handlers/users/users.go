package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexferl/echo-boilerplate/util"
)

type ListUsersResponse struct {
	Users []*Public `json:"users"`
}

func (h *Handler) ListUsers(c echo.Context) error {
	logger := c.Get("logger").(zerolog.Logger)
	page, perPage, limit, skip := util.ParsePaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := h.Mapper.Count(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed counting tasks: %v", err)
	}

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	res, err := h.Mapper.Find(ctx, nil, Users{}, opts)
	if err != nil {
		logger.Error().Err(err).Msg("failed getting users")
		return err
	}

	users := res.(Users)
	util.SetPaginationHeaders(c.Request(), c.Response().Header(), int(count), page, perPage)

	return h.Validate(c, http.StatusOK, &ListUsersResponse{Users: users.Public()})
}
