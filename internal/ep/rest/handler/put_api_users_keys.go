package handler

import (
	"net/http"

	_ "github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

// putAPIUsersKeys godoc
// @Summary      Обновление мастер-ключей пользователя
// @Description  Генерирует новую пару RSA ключей (приватный перешифровывается мастер-паролем) и возвращает новый публичный ключ
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dto.ChangeKeysResultDto
// @Failure      403  {object}  dto.ErrorDto
// @Failure      500  {object}  dto.ErrorDto
// @Router       /api/users/keys [put]
func (cr *AppChiRouter) putAPIUsersKeys(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("putAPIUsersKeys start")
	defer cr.log.Info("putAPIUsersKeys finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	res, err := cr.userFacade.ChangeKeys(r.Context())
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("putAPIUsersKeys", "change keys", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
