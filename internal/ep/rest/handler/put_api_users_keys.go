package handler

import (
	"net/http"

	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

func (cr *AppChiRouter) putAPIUsersKeys(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("putAPIUsersKeys start")
	defer cr.log.Info("putAPIUsersKeys finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	res, err := cr.userFacade.ChangeKeys(r.Context())
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("putAPIUsersKeys", "change keys", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
