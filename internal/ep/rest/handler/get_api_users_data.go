package handler

import (
	"net/http"
)

func (cr *AppChiRouter) getApiUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getApiUsersData start")
	defer cr.log.Info("getApiUsersData finish")

	// ToDo: implement
}
