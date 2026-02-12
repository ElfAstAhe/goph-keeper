package handler

import (
	"net/http"

	_ "github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/go-chi/chi/v5"
)

// getAPIUsersData godoc
// @Summary      Получение данных по ID
// @Description  Возвращает расшифрованный секрет (пароль, карту или файл) по его идентификатору
// @Tags         user-data
// @Security     BearerAuth
// @Param        id   path      string  true  "Идентификатор записи"
// @Success      200  {object}  dto.UserDataDto
// @Failure      403  {object}  dto.ErrorDto
// @Failure      404  {object}  dto.ErrorDto
// @Failure      500  {object}  dto.ErrorDto
// @Router       /api/users/data/{id} [get]
func (cr *AppChiRouter) getAPIUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getAPIUsersData start")
	defer cr.log.Info("getAPIUsersData finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		cr.renderError(rw, errs.NewAppInvalidArgumentError("id", "id is required"))

		return
	}

	res, err := cr.userDataFacade.Get(r.Context(), id)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("getAPIUsersData", "get user data", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
