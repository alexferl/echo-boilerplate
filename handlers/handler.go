package handlers

type (
	// Handler represents the structure of our resource
	Handler struct {
	}
)

// ErrorResponse holds an error message
type ErrorResponse struct {
	Message string `json:"error"`
}
