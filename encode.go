package errorhelper

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrUnknown        = errors.New("unknown error")
	ErrDuplicate      = errors.New("duplicate")
	ErrNotImplemented = errors.New("not implemented")
	ErrNoContent      = errors.New("no content")
	ErrTooLarge       = errors.New("payload too large")
)

type WebError struct {
	Error struct {
		Code   int    `json:"code"`
		String string `json:"string"`
	} `json:"error"`
}

type Encoder interface {
	Encode(ctx context.Context, err error, w http.ResponseWriter)
}

type DefaultEncoder struct {
	Mapper func(error) int
}

func NewDefaultEncoder(mapper func(error) int) *DefaultEncoder {
	return &DefaultEncoder{Mapper: mapper}
}

func (d *DefaultEncoder) Encode(ctx context.Context, err error, w http.ResponseWriter) {
	var response WebError
	code := d.statusCode(err)
	response.Error.Code = code
	response.Error.String = strings.ToTitle(err.Error())

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(response)
}

func (d *DefaultEncoder) statusCode(err error) int {
	if d.Mapper != nil {
		return d.Mapper(err)
	}

	return DefaultMapper(err)
}

// DefaultMapper provides the default error to status code mapping
func DefaultMapper(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrDuplicate:
		return http.StatusConflict
	case ErrNotImplemented:
		return http.StatusNotImplemented
	case ErrNoContent:
		return http.StatusNoContent
	case ErrTooLarge:
		return http.StatusRequestEntityTooLarge
	default:
		log.Println(err)
		return http.StatusInternalServerError
	}
}
