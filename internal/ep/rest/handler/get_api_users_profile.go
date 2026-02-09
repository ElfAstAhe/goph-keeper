package handler

import "net/http"

func (cr *AppChiRouter) getApiUsersGetProfile(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getApiUsersGetProfile start")
	defer cr.log.Info("getApiUsersGetProfile finish")

	// ToDo: implement
}
