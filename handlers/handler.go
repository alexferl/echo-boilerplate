package handlers

// Handler represents the structure of our resource
type Handler struct{}

// ErrorResponse holds an error message
type ErrorResponse struct {
	Message string `json:"error"`
}

type Response struct {
	Message string `json:"message"`
}
