package middleware

import (
	"context"
	"net/http"

	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

type AuthExtractorMiddleware struct {
	authHelper    *utils.AuthHelper
	jwtHTTPHelper *utils.JWTHTTPHelper
	log           logger.Logger
}

func NewAuthExtractorMiddleware(authHelper *utils.AuthHelper, jwtHTTPHelper *utils.JWTHTTPHelper, logger logger.Logger) *AuthExtractorMiddleware {
	return &AuthExtractorMiddleware{
		authHelper:    authHelper,
		jwtHTTPHelper: jwtHTTPHelper,
		log:           logger.GetLogger("authExtractorMiddleware"),
	}
}

func (aem *AuthExtractorMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		aem.log.Info("NewAuthExtractorMiddleware.Handle start")
		defer aem.log.Info("NewAuthExtractorMiddleware.Handle finish")

		userInfo, err := aem.authHelper.UserInfoFromHTTPRequest(r)
		if err != nil {
			aem.log.Errorf("AUTH MW: [%s] [%s] no user info: %v", r.Method, r.RequestURI, err)

			next.ServeHTTP(rw, r)
		} else {
			aem.log.Info("NewAuthExtractorMiddleware.Handle userInfo placed into req context", userInfo)

			reqCtx := context.WithValue(r.Context(), aem.authHelper.GetUserInfoContextName(), userInfo)

			req := r.WithContext(reqCtx)

			next.ServeHTTP(rw, req)
		}
	})
}
