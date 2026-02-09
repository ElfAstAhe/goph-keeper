package handler

import (
	"net/http"

	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/go-chi/chi/v5"
)

func (cr *AppChiRouter) getAPIUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getAPIUsersData start")
	defer cr.log.Info("getAPIUsersData finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		cr.renderError(rw, errs.NewAppInvalidArgumentError("id", "id is required"))

		return
	}

	res, err := cr.userDataFacade.Get(r.Context(), id)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("getAPIUsersData", "get user data", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
