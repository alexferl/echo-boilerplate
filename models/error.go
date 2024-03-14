package models

import "fmt"

// Error represents a model error
type Error struct {
	Kind    Kind
	Message string
}

// Kind defines supported error types.
type Kind uint8

const (
	Other Kind = iota + 1 // Unclassified error.
	Conflict
	Permission
)

func (k Kind) String() string {
	return [...]string{"other", "conflict", "permission"}[k-1]
}

// NewError instantiates a new error.
func NewError(err error, kind Kind) error {
	e := &Error{
		Kind:    kind,
		Message: err.Error(),
	}
	return e
}

// Error returns the message.
func (e *Error) Error() string {
	return fmt.Sprintf("kind=%s, message=%v", e.Kind, e.Message)
}
