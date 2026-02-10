package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

// postAPIAuthLogin godoc
// @Summary      Аутентификация пользователя
// @Description  Проверяет учетные данные и возвращает JWT токен доступа
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.LoginDto  true  "Данные для входа"
// @Success      200    {object}  dto.LoginResultDto
// @Failure      400    {object}  dto.ErrorDto
// @Failure      401    {object}  dto.ErrorDto
// @Failure      500    {object}  dto.ErrorDto
// @Router       /api/auth/login [post]
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
