package handler

import (
	"net/http"
)

func (cr *AppChiRouter) deleteApiUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("deleteApiUsersData start")
	defer cr.log.Info("deleteApiUsersData finish")

	// ToDo: implement
}
