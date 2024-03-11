package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	testCases := []struct {
		name string
		msg  string
		kind Kind
	}{
		{"Other", "other", Other},
		{"Exist", "exist", Exist},
		{"NotExist", "not_exist", NotExist},
		{"Deleted", "deleted", Deleted},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := NewError(nil, tc.kind, tc.msg)
			var se *Error
			if errors.As(e, &se) {
				assert.Equal(t, tc.msg, se.Message)
				assert.Equal(t, tc.kind, se.Kind)
			}
		})
	}
}

func TestKind(t *testing.T) {
	testCases := []struct {
		name string
		msg  string
		kind Kind
	}{
		{"Other", "other", Other},
		{"Exist", "exist", Exist},
		{"NotExist", "not_exist", NotExist},
		{"Deleted", "deleted", Deleted},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.msg, tc.kind.String())
		})
	}
}

func TestError_Error(t *testing.T) {
	errMsg := "my error"
	err := errors.New(errMsg)
	msg := "my msg"
	e := NewError(err, Other, msg)
	assert.Equal(t, fmt.Sprintf("kind=other, message=%s, internal=%s", msg, errMsg), e.Error())
}

func TestError_Error_No_Internal(t *testing.T) {
	msg := "my msg"
	e := NewError(nil, Other, msg)
	assert.Equal(t, fmt.Sprintf("kind=other, message=%s", msg), e.Error())
}

func TestError_Unwrap(t *testing.T) {
	errMsg := "my error"
	err := errors.New(errMsg)
	e := NewError(err, Other, "")
	assert.Equal(t, errMsg, errors.Unwrap(e).Error())
}
