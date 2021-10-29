package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go-react-app/internal/app"
	"go-react-app/internal/config"
	"go-react-app/internal/handlers/account"
	"go-react-app/internal/handlers/session"
	"go-react-app/internal/middlewares"
	"go-react-app/internal/util/logger"

	"github.com/gorilla/mux"
)

func New(a *app.App) http.Handler {
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.Use(middlewares.NewRecovery(a.Logger).Middleware)
	router.Use(middlewares.NewRequestLogger(a.Logger).Middleware)
	if a.Config.Common.IsDebug {
		router.Use(middlewares.NewRequestBodyLogger(a.Logger).Middleware)
	}

	apiRouter := router.PathPrefix("/api").Subrouter()
	RouteApi(apiRouter, a)

	RouteFrontend(a.Config.Frontend, router, a.Logger)

	return router
}

func RouteFrontend(conf config.FrontendConfig, router *mux.Router, log *logger.Logger) {
	files, err := ioutil.ReadDir(conf.PathToDist)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}

	fileServer := http.FileServer(http.Dir(conf.PathToDist))

	for _, v := range files {
		if v.Name() == conf.IndexPath {
			continue
		}

		prefix := "/" + v.Name()
		if v.IsDir() {
			prefix += "/"
		}

		log.Debug().Msgf("Static route for %s", prefix)
		router.PathPrefix(prefix).Handler(http.StripPrefix("/", fileServer))
	}

	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.ServeFile(rw, r, conf.PathToDist+"/"+conf.IndexPath)
	})
}

func RouteApi(r *mux.Router, a *app.App) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	accountRouter := r.PathPrefix("/account").Subrouter()
	accountRouter.Handle("/login", account.NewLogin(a.Logger, a.DbContext, a.DbContext))
	accountRouter.Handle("/register", account.NewRegister(a.Logger, a.DbContext))

	sessionRouter := r.PathPrefix("/session").Subrouter()
	sessionRouter.Handle("/check", session.NewCheck(a.Logger, a.DbContext))

	securedRouter := r.PathPrefix("/").Subrouter()
	securedRouter.Use(middlewares.NewAuthentication(a.Logger, a.DbContext).Middleware)
	securedRouter.Handle("/hello", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) { rw.Write([]byte("hello")) }))
}
