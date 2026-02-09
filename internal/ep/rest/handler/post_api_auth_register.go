package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

func (cr *AppChiRouter) postAPIAuthRegister(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postAPIAuthRegister start")
	defer cr.log.Info("postAPIAuthRegister finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	register := dto.NewEmptyRegisterDto()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(register); err != nil {
		cr.renderError(rw, apperrs.NewEpMappingError("postAPIAuthRegister", "http req", "RegisterDto", "decode json", err))

		return
	}

	res, err := cr.authFacade.Register(r.Context(), register)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("postAPIAuthRegister", "register user", err))

		return
	}

	cr.renderJSON(rw, http.StatusCreated, res)
}
