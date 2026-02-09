package handler

import (
	"net/http"

	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

func (cr *AppChiRouter) getAPIUsersProfile(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getApiUsersGetProfile start")
	defer cr.log.Info("getApiUsersGetProfile finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	res, err := cr.userFacade.GetProfile(r.Context())
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("getAPIUsersProfile", "user profile", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
