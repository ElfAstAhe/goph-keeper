package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

// putAPIUsersData godoc
// @Summary      Обновление существующего секрета
// @Description  Перешифровывает и сохраняет изменения в секрете по его ID
// @Tags         user-data
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.UserDataDto  true  "Новые данные секрета (ID обязателен)"
// @Success      200    {object}  dto.UserDataDto
// @Failure      400    {object}  dto.ErrorDto
// @Failure      403    {object}  dto.ErrorDto
// @Failure      404    {object}  dto.ErrorDto
// @Failure      500    {object}  dto.ErrorDto
// @Router       /api/users/data [put]
func (cr *AppChiRouter) putAPIUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("putAPIUsersData start")
	defer cr.log.Info("putAPIUsersData finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	decoder := json.NewDecoder(r.Body)
	income := dto.NewEmptyUserDataDto()
	if err := decoder.Decode(income); err != nil {
		cr.renderError(rw, apperrs.NewEpMappingError("putAPIUsersData", "http req", "UserDataDto", "decode json", err))

		return
	}

	res, err := cr.userDataFacade.Change(r.Context(), income)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("putAPIUsersData", "change user data", err))

		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
