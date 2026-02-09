package handler

import "net/http"

func (cr *AppChiRouter) getApiUsersProfile(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getApiUsersGetProfile start")
	defer cr.log.Info("getApiUsersGetProfile finish")

	// ToDo: implement
}
