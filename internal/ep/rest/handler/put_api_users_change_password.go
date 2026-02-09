package handler

import (
	"net/http"
)

func (cr *AppChiRouter) putApiUsersPassword(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("putApiUsersChangePassword start")
	defer cr.log.Info("putApiUsersChangePassword finish")

	// ToDo: implement
}
