package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

func (cr *AppChiRouter) postAPIAuthLogin(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postAPIAuthLogin start")
	defer cr.log.Info("postAPIAuthLogin finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	login := dto.NewEmptyLoginDto()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(login); err != nil {
		cr.renderError(rw, apperrs.NewEpMappingError("postAPIAuthLogin", "http req", "LoginDto", "decode json", err))

		return
	}

	res, err := cr.authFacade.Login(r.Context(), login)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("postAPIAuthLogin", "login user", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
