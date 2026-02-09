package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

func (cr *AppChiRouter) putAPIUsersPassword(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("putAPIUsersPassword start")
	defer cr.log.Info("putAPIUsersPassword finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	income := dto.NewEmptyUpdatePasswordDto()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(income); err != nil {
		cr.renderError(rw, apperrs.NewEpMappingError("putAPIUsersPassword", "http req", "UpdatePasswordDto", "decode json: %w", err))

		return
	}

	err := cr.userFacade.UpdatePassword(r.Context(), income)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("putAPIUsersPassword", "update password", err))

		return
	}

	cr.renderEmpty(rw, http.StatusOK)
}
