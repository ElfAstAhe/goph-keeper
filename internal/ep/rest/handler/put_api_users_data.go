package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

func (cr *AppChiRouter) putAPIUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("putAPIUsersData start")
	defer cr.log.Info("putAPIUsersData finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	decoder := json.NewDecoder(r.Body)
	income := dto.NewEmptyUserDataDto()
	if err := decoder.Decode(income); err != nil {
		cr.renderError(rw, apperrs.NewEpMappingError("putAPIUsersData", "http req", "UserDataDto", "decode json", err))

		return
	}

	res, err := cr.userDataFacade.Change(r.Context(), income)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("putAPIUsersData", "change user data", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
