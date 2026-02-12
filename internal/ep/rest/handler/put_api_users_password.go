package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

// putAPIUsersPassword godoc
// @Summary      Смена пароля пользователя
// @Description  Проверяет старый пароль и устанавливает новый хеш (может потребоваться перешифрование RSA-ключей на клиенте)
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.UpdatePasswordDto  true  "Старый и новый пароли"
// @Success      200    {object}  nil  "Пароль успешно изменен"
// @Failure      400    {object}  dto.ErrorDto
// @Failure      403    {object}  dto.ErrorDto
// @Failure      500    {object}  dto.ErrorDto
// @Router       /api/users/password [put]
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
