package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/admiralobvious/echo-boilerplate/models"
)

// Login finds a user and returns a JWT if it exists
func (h *Handler) Login(c echo.Context) error {
	user := &models.User{}

	if err := c.Bind(user); err != nil {
		e := fmt.Sprint("Must provide an email and password")
		return c.JSON(http.StatusBadRequest, ErrorResponse{e})
	}

	err := h.DB.FindUserByEmail(user.Email, user)
	if err != nil {
		e := fmt.Sprint("User not found")
		return c.JSON(http.StatusNotFound, ErrorResponse{e})
	}

	t, err := user.GenerateJWT()
	if err != nil {
		e := fmt.Sprint("Error generating token")
		return c.JSON(http.StatusInternalServerError, ErrorResponse{e})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": t})
}
