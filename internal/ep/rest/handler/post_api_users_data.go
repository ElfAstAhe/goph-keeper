package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

// postAPIUsersData godoc
// @Summary      Создание нового секрета
// @Description  Шифрует и сохраняет новый секрет (пароль, текст, карту или файл) для текущего пользователя
// @Tags         user-data
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.UserDataDto  true  "Данные секрета (Kind, Name, TextData/BinaryData)"
// @Success      201    {object}  dto.UserDataDto
// @Failure      400    {object}  dto.ErrorDto
// @Failure      403    {object}  dto.ErrorDto
// @Failure      409    {object}  dto.ErrorDto
// @Failure      500    {object}  dto.ErrorDto
// @Router       /api/users/data [post]
func (cr *AppChiRouter) postAPIUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postAPIUsersData start")
	defer cr.log.Info("postAPIUsersData finish")

	if r.Body != nil {
		defer r.Body.Close()
	}

	income := dto.NewEmptyUserDataDto()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(income); err != nil {
		cr.renderError(rw, apperrs.NewEpMappingError("postAPIUsersData", "http req", "UserDataDto", "decode json", err))

		return
	}

	res, err := cr.userDataFacade.Create(r.Context(), income)
	if err != nil {
		cr.renderError(rw, apperrs.NewEpCommonError("postAPIUsersData", "create user data", err))

		return
	}

	cr.renderJSON(rw, http.StatusCreated, res)
}
