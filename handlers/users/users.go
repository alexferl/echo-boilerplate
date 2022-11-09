package users

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexferl/echo-boilerplate/util"
)

type ShortUser struct {
	Id       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
}

type UsersResponse struct {
	Users []*ShortUser `json:"users"`
}

func (h *Handler) Users(c echo.Context) error {
	var page int
	pageQuery := c.QueryParam("page")
	page, _ = strconv.Atoi(pageQuery)

	var perPage int
	perPageQuery := c.QueryParam("per_page")
	perPage, _ = strconv.Atoi(perPageQuery)

	limit := perPage
	skip := 0
	if page > 1 {
		skip = (page * perPage) - perPage
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := h.Mapper.Count(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed counting users: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	result, err := h.Mapper.Find(ctx, nil, []*ShortUser{}, opts)
	if err != nil {
		return fmt.Errorf("failed getting users: %v", err)
	}

	uri := fmt.Sprintf("http://%s%s", c.Request().Host, c.Request().URL.Path)
	util.Paginate(c.Response().Header(), int(count), page, perPage, uri)

	resp := &UsersResponse{Users: result.([]*ShortUser)}

	return h.Validate(c, http.StatusOK, resp)
}
