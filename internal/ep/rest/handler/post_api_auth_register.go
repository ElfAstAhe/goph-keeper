package handler

import (
	"net/http"
)

func (cr *AppChiRouter) postApiAuthRegister(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postApiAuthRegister start")
	defer cr.log.Info("postApiAuthRegister finish")

	// ToDo: implement
}
