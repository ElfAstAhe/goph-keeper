package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

// postAPIAuthRegister godoc
// @Summary      Регистрация пользователя
// @Description  Создает новый аккаунт, генерирует пару RSA ключей и возвращает данные созданного пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.RegisterDto  true  "Данные для регистрации"
// @Success      201    {object}  dto.RegisterResultDto
// @Failure      400    {object}  dto.ErrorDto
// @Failure      409    {object}  dto.ErrorDto
// @Failure      500    {object}  dto.ErrorDto
// @Router       /api/auth/register [post]
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
