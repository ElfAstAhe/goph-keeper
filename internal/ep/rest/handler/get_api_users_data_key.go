package handler

import (
	"net/http"

	_ "github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/go-chi/chi/v5"
)

// getAPIUsersDataKey godoc
// @Summary      Получение данных по бизнес-ключу
// @Description  Возвращает расшифрованный секрет, используя тип данных и его имя (например, credential и Gmail)
// @Tags         user-data
// @Security     BearerAuth
// @Param        dataKind  path      string  true  "Тип данных (credential, plaintext, binary, bankcard)"
// @Param        name      path      string  true  "Имя секрета"
// @Success      200       {object}  dto.UserDataDto
// @Failure      400       {object}  dto.ErrorDto
// @Failure      403       {object}  dto.ErrorDto
// @Failure      404       {object}  dto.ErrorDto
// @Failure      500       {object}  dto.ErrorDto
// @Router       /api/users/data/{dataKind}/{name} [get]
func (cr *AppChiRouter) getAPIUsersDataKey(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getAPIUsersDataKey start")
	defer cr.log.Info("getAPIUsersDataKey finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	dataKind := chi.URLParam(r, "dataKind")
	if dataKind == "" {
		cr.renderError(rw, errs.NewAppInvalidArgumentError("dataKind", "must be specified"))

		return
	}
	name := chi.URLParam(r, "name")
	if name == "" {
		cr.renderError(rw, errs.NewAppInvalidArgumentError("name", "must be specified"))

		return
	}

	res, err := cr.userDataFacade.GetByKey(r.Context(), dataKind, name)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("getAPIUsersDataKey", "get user data by key", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
