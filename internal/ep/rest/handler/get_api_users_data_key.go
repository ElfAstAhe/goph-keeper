package handler

import (
	"net/http"
)

func (cr *AppChiRouter) getApiUsersDataKey(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getApiUsersDataKey start")
	defer cr.log.Info("getApiUsersDataKey finish")

	// ToDo: implement
}
