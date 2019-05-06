package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/admiralobvious/echo-boilerplate/models"
)

// CreateUser creates a new user
func (h *Handler) CreateUser(c echo.Context) error {
	u := &models.User{}

	if err := c.Bind(u); err != nil {
		e := fmt.Sprint("Must provide an email and password")
		return c.JSON(http.StatusBadRequest, ErrorResponse{e})
	}

	user := models.NewUser(u.Email, u.Password)

	err := h.DB.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.DB.DeleteUser(id)

	if err != nil {
		e := fmt.Sprint("User not found")
		return c.JSON(http.StatusNotFound, ErrorResponse{e})
	}

	return c.JSON(http.StatusNoContent, nil)
}

// GetUser returns a single user
func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")
	u := &models.User{}
	err := h.DB.FindUserById(id, u)

	if u.Id == "" || err != nil {
		e := fmt.Sprint("User not found")
		return c.JSON(http.StatusNotFound, ErrorResponse{e})
	}

	return c.JSON(http.StatusOK, u)
}

// GetUsers returns all the users
func (h *Handler) GetUsers(c echo.Context) error {
	users := &[]models.User{}
	err := h.DB.GetAllUsers(users)

	if err != nil {
		e := fmt.Sprint("Users not found")
		return c.JSON(http.StatusNotFound, ErrorResponse{e})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"users": users})
}

// UpdateUser updates a single user
func (h *Handler) UpdateUser(c echo.Context) error {
	user := &models.User{}

	if err := c.Bind(user); err != nil {
		e := fmt.Sprint("Must provide an email and password")
		return c.JSON(http.StatusBadRequest, ErrorResponse{e})
	}

	err := h.DB.CreateUser(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}
