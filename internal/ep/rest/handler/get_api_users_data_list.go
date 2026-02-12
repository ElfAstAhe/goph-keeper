package handler

import (
	"net/http"

	_ "github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

// getAPIUsersDataList godoc
// @Summary      Получение списка всех секретов
// @Description  Возвращает список метаданных всех секретов текущего пользователя (без зашифрованных данных)
// @Tags         user-data
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   dto.UserDataDto
// @Failure      403  {object}  dto.ErrorDto
// @Failure      500  {object}  dto.ErrorDto
// @Router       /api/users/data [get]
func (cr *AppChiRouter) getAPIUsersDataList(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getAPIUsersDataList start")
	defer cr.log.Info("getAPIUsersDataList finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	res, err := cr.userDataFacade.ListAll(r.Context())
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("getAPIUsersDataList", "get all user data", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
