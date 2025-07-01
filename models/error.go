package models

import (
	"fmt"
	"net/http"
)

type JsonError struct {
	Status  int            `json:"code"`
	Message string         `json:"message"`
	Context map[string]any `json:"context,omitempty"`
}

func (e JsonError) Error() string        { return e.Message }
func (e JsonError) Code() int            { return e.Status }
func (e JsonError) JsonError() JsonError { return e }

type Error interface {
	error
	JsonError() JsonError
	Code() int
}

type InternalError struct{ Err error }

func (ie InternalError) Error() string { return ie.Err.Error() }
func (InternalError) JsonError() JsonError {
	return JsonError{
		Status:  http.StatusInternalServerError,
		Message: "internal error",
	}
}
func (InternalError) Code() int { return http.StatusInternalServerError }

func Internal(err error) InternalError {
	return InternalError{err}
}

type InvalidError[T any] struct {
	Field string
	Value T
}

func (i InvalidError[T]) Error() string {
	return fmt.Sprintf("invalid %s: %v", i.Field, i.Value)
}

func (InvalidError[T]) Code() int {
	return http.StatusBadRequest
}
func (i InvalidError[T]) JsonError() JsonError {
	return JsonError{
		Status:  http.StatusBadRequest,
		Message: fmt.Sprint("invalid ", i.Field),
		Context: map[string]any{i.Field: i.Value},
	}
}

func Invalid[T any](field string, value T) InvalidError[T] {
	return InvalidError[T]{field, value}
}

type NotAllowed struct {
	User, Id int64
	Thing    string
}

func (f NotAllowed) Error() string {
	return fmt.Sprintf("user %x is not allowed to %s %x", f.User, f.Thing, f.Id)
}
func (NotAllowed) Code() int { return http.StatusForbidden }
func (f NotAllowed) JsonError() JsonError {
	return JsonError{
		Status:  http.StatusForbidden,
		Message: fmt.Sprintf("user is not allowed to %s", f.Thing),
		Context: map[string]any{"id": f.Id, "user_id": f.User},
	}
}
