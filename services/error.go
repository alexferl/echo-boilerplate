package services

import "fmt"

// Error represents an error that could be wrapping another error.
type Error struct {
	Internal error
	Kind     Kind
	Message  string
}

// Kind defines supported error types.
type Kind uint8

const (
	Other    Kind = iota + 1 // Unclassified error.
	Exist                    // Item already exist.
	NotExist                 // Item does not exist.
	Deleted                  // Item was deleted.
	Conflict
	Permission
)

func (k Kind) String() string {
	return [...]string{"other", "exist", "not_exist", "deleted", "conflict", "permission"}[k-1]
}

// NewError instantiates a new error.
func NewError(err error, code Kind, message string) error {
	e := &Error{
		Internal: err,
		Kind:     code,
		Message:  message,
	}
	return e
}

// Error returns the message, when wrapping errors the wrapped error is appended.
func (e *Error) Error() string {
	if e.Internal == nil {
		return fmt.Sprintf("kind=%s, message=%v", e.Kind, e.Message)
	}
	return fmt.Sprintf("kind=%s, message=%v, internal=%v", e.Kind, e.Message, e.Internal)
}

// Unwrap returns the wrapped error, if any.
func (e *Error) Unwrap() error {
	return e.Internal
}
