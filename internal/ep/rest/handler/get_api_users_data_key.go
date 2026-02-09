package handler

import (
	"net/http"

	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/go-chi/chi/v5"
)

func (cr *AppChiRouter) getAPIUsersDataKey(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getAPIUsersDataKey start")
	defer cr.log.Info("getAPIUsersDataKey finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	dataKind := chi.URLParam(r, "dataKind")
	if dataKind == "" {
		cr.renderError(rw, errs.NewAppInvalidArgumentError("dataKind", "must be specified"))

		return
	}
	name := chi.URLParam(r, "name")
	if name == "" {
		cr.renderError(rw, errs.NewAppInvalidArgumentError("name", "must be specified"))

		return
	}

	res, err := cr.userDataFacade.GetByKey(r.Context(), dataKind, name)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("getAPIUsersDataKey", "get user data by key", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
