package handler

import (
	"net/http"

	_ "github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

// getAPIUsersProfile godoc
// @Summary      Получение профиля текущего пользователя
// @Description  Возвращает публичную информацию о пользователе (ID, логин, роль, ФИО) из JWT-контекста
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dto.UserDto
// @Failure      400  {object}  dto.ErrorDto
// @Failure      403  {object}  dto.ErrorDto
// @Failure      500  {object}  dto.ErrorDto
// @Router       /api/users/profile [get]
func (cr *AppChiRouter) getAPIUsersProfile(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getApiUsersGetProfile start")
	defer cr.log.Info("getApiUsersGetProfile finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	res, err := cr.userFacade.GetProfile(r.Context())
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("getAPIUsersProfile", "user profile", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
