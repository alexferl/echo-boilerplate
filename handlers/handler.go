package handlers

import (
	"github.com/admiralobvious/echo-boilerplate/database"
)

type (
	// Handler represents the structure of our resource
	Handler struct {
		DB database.DB
	}
)

// ErrorResponse holds an error message
type ErrorResponse struct {
	Message string `json:"error"`
}
