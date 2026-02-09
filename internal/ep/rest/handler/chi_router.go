package handler

import (
	"net/http"
	"time"

	"github.com/ElfAstAhe/goph-keeper/internal/app/server/config"
	"github.com/ElfAstAhe/goph-keeper/internal/ep/facade"
	appmware "github.com/ElfAstAhe/goph-keeper/internal/ep/rest/middleware"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AppChiRouter struct {
	router         *chi.Mux
	log            logger.Logger
	conf           *config.Config
	authHelper     *utils.AuthHelper
	jwtHTTPHelper  *utils.JWTHTTPHelper
	authFacade     facade.AuthFacade
	userFacade     facade.UserFacade
	userDataFacade facade.UserDataFacade
}

func NewAppChiRouter(
	authFacade facade.AuthFacade,
	userFacade facade.UserFacade,
	userDataFacade facade.UserDataFacade,
	authHelper *utils.AuthHelper,
	jwtHTTPHelper *utils.JWTHTTPHelper,
	conf *config.Config,
	logger logger.Logger,
) *AppChiRouter {
	res := &AppChiRouter{
		router:         chi.NewRouter(),
		log:            logger,
		conf:           conf,
		authHelper:     authHelper,
		jwtHTTPHelper:  jwtHTTPHelper,
		authFacade:     authFacade,
		userFacade:     userFacade,
		userDataFacade: userDataFacade,
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
	cr.router.Use(appmware.NewAuthExtractorMiddleware(cr.jwtHTTPHelper, logger).Handle)
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
			r.Post("/login", cr.postAPIAuthLogin)       // POST /api/auth/login
			r.Post("/register", cr.postAPIAuthRegister) // POST /api/auth/register
		})
		// users sub-router
		r.Route("/users", func(r chi.Router) {
			r.Get("/profile", cr.getAPIUsersProfile)
			r.Put("/keys", cr.putAPIUsersKeys)
			r.Put("/password", cr.putAPIUsersPassword)

			// data sub-router
			r.Route("/data", func(r chi.Router) {
				r.Get("/{id}", cr.getAPIUsersData)
				r.Get("/{dataKind}/{name}", cr.getAPIUsersDataKey)
				r.Post("/", cr.postAPIUsersData)
				r.Put("/{id}", cr.putAPIUsersData)
				r.Delete("/{id}", cr.deleteAPIUsersData)
			})
		})
	})
}
