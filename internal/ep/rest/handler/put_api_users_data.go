package handler

import (
	"net/http"
)

func (cr *AppChiRouter) putApiUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("putApiUsersData start")
	defer cr.log.Info("putApiUsersData finish")

	// ToDo: implement
}
