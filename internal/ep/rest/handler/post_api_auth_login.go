package handler

import (
	"net/http"
)

func (cr *AppChiRouter) postApiAuthLogin(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postApiAuthLogin start")
	defer cr.log.Info("postApiAuthLogin finish")

	// ToDo: implement
}
