package handler

import (
	"net/http"
)

func (cr *AppChiRouter) postApiUsersData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postApiUsersData start")
	defer cr.log.Info("postApiUsersData finish")

	// ToDo: implement
}
