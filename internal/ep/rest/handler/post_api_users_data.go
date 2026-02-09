package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

func (cr *AppChiRouter) postAPIUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postAPIUsersData start")
	defer cr.log.Info("postAPIUsersData finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	income := dto.NewEmptyUserDataDto()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(income); err != nil {
		cr.renderError(rw, apperrs.NewEpMappingError("postAPIUsersData", "http req", "UserDataDto", "decode json", err))

		return
	}

	res, err := cr.userDataFacade.Create(r.Context(), income)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("postAPIUsersData", "create user data", err))

		return
	}

	cr.renderJSON(rw, http.StatusCreated, res)
}
