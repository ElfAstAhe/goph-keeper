package handler

import (
	"net/http"

	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

func (cr *AppChiRouter) getAPIUsersDataList(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getAPIUsersDataList start")
	defer cr.log.Info("getAPIUsersDataList finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	res, err := cr.userDataFacade.ListAll(r.Context())
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("getAPIUsersDataList", "get all user data", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
