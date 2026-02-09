package handler

import (
	"net/http"

	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/go-chi/chi/v5"
)

func (cr *AppChiRouter) deleteAPIUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("deleteAPIUsersData start")
	defer cr.log.Info("deleteAPIUsersData finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		cr.renderError(rw, errs.NewAppInvalidArgumentError("id", "id is required"))

		return
	}

	err := cr.userDataFacade.Remove(r.Context(), id)
	if err.Error != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("deleteAPIUsersData", "remove user data", err))

		return
	}

	cr.renderEmpty(rw, http.StatusNoContent)
}
