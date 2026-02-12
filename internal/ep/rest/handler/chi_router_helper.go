package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func (cr *AppChiRouter) renderError(rw http.ResponseWriter, err error) {
	// 500 Internal Server Error
	status := http.StatusInternalServerError

	// 400 Bad Request
	if errors.As(err, &apperrs.ErrBllValidate) ||
		errors.As(err, &apperrs.ErrDalValidate) ||
		errors.As(err, &apperrs.ErrEpMapping) ||
		errors.As(err, &apperrs.ErrEpValidate) ||
		errors.As(err, &errs.ErrAppInvalidArgument) {
		status = http.StatusBadRequest
	} else
	// 401 Unauthorized
	if errors.As(err, &errs.ErrAuthUnauthorized) {
		status = http.StatusUnauthorized
	} else
	// 403 Forbidden
	if errors.As(err, &errs.ErrAuthForbidden) {
		status = http.StatusForbidden
	} else
	// 404 Not Found
	if errors.As(err, &apperrs.ErrDalNotFound) {
		status = http.StatusNotFound
	} else
	// 409 Conflict
	if errors.As(err, &apperrs.ErrDalAlreadyExists) {
		status = http.StatusConflict
	} else
	// 410 Gone
	if errors.As(err, &apperrs.ErrDalSoftDeleted) {
		status = http.StatusGone
	}

	// отправляем ошибку
	cr.renderJSON(rw, status, dto.NewErrorDtoFromError(status, err))
}

func (cr *AppChiRouter) renderJSON(rw http.ResponseWriter, status int, data any) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(rw).Encode(data)
	}
}

func (cr *AppChiRouter) renderEmpty(rw http.ResponseWriter, status int) {
	rw.Header().Set("Content-Type", "*/*")
	rw.WriteHeader(status)
}
