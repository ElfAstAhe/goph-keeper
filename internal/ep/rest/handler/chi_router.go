package handler

import (
	"net/http"
	"time"

	"github.com/ElfAstAhe/goph-keeper/internal/app/server/config"
	"github.com/ElfAstAhe/goph-keeper/internal/ep/facade"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AppChiRouter struct {
	router     *chi.Mux
	log        logger.Logger
	conf       *config.Config
	authFacade facade.AuthFacade
	userFacade facade.UserFacade
}

func NewAppChiRouter(authFacade facade.AuthFacade, userFacade facade.UserFacade, conf *config.Config, logger logger.Logger) *AppChiRouter {
	res := &AppChiRouter{
		router:     chi.NewRouter(),
		log:        logger.GetLogger("app router"),
		conf:       conf,
		authFacade: authFacade,
		userFacade: userFacade,
	}

	// setup middleware
	res.setupMiddleware(logger)

	// mount
	res.router.Mount("/debug", middleware.Profiler())

	// setup routes
	res.setupRoutes()

	return res
}

func (cr *AppChiRouter) GetRouter() http.Handler {
	return cr.router
}

func (cr *AppChiRouter) setupMiddleware(logger logger.Logger) {
	// jwt auth extractor - extract user info from token
	// ..
	// requestID
	cr.router.Use(middleware.RequestID)
	// realIP
	cr.router.Use(middleware.RealIP)
	// compress
	// ..
	// decompress
	// ..
	// income/outcome logger
	// ..
	// recoverer
	cr.router.Use(middleware.Recoverer)
	// timeout
	cr.router.Use(middleware.Timeout(20 * time.Second))
}

func (cr *AppChiRouter) setupRoutes() {
	// api
	cr.router.Route("/api", func(r chi.Router) {
		// auth sub-router
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", cr.postApiAuthLogin)       // POST /api/auth/login
			r.Post("/register", cr.postApiAuthRegister) // POST /api/auth/register
		})
		// users sub-router
		r.Route("/users", func(r chi.Router) {
			r.Get("/{username}/profile", cr.getApiUsersGetProfile)
			r.Put("/{username}/keys", cr.putApiUsersUpdateKeys)
			r.Put("/{username}/password", cr.putApiUsersChangePassword)
		})
	})
}
