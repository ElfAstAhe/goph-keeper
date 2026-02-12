package handler

import (
	"net/http"

	_ "github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/go-chi/chi/v5"
)

// deleteAPIUsersData godoc
// @Summary      Удаление секрета
// @Description  Удаляет запись с секретными данными по её ID (Soft Delete)
// @Tags         user-data
// @Security     BearerAuth
// @Param        id   path      string  true  "ID записи"
// @Success      204  {object}  nil     "Успешное удаление (нет контента)"
// @Failure      403  {object}  dto.ErrorDto
// @Failure      404  {object}  dto.ErrorDto
// @Failure      500  {object}  dto.ErrorDto
// @Router       /api/users/data/{id} [delete]
func (cr *AppChiRouter) deleteAPIUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("deleteAPIUsersData start")
	defer cr.log.Info("deleteAPIUsersData finish")

	id := chi.URLParam(r, "id")
	if id == "" {
		cr.renderError(rw, errs.NewAppInvalidArgumentError("id", "id is required"))

		return
	}

	err := cr.userDataFacade.Remove(r.Context(), id)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("deleteAPIUsersData", "remove user data", err))

		return
	}

	cr.renderEmpty(rw, http.StatusNoContent)
}
