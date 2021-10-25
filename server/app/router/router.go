package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go-react-app/server/app"
	"go-react-app/server/config"
	"go-react-app/server/handlers/account"
	"go-react-app/server/handlers/session"
	"go-react-app/server/middlewares"
	"go-react-app/server/util/logger"

	"github.com/gorilla/mux"
)

func New(a *app.App) http.Handler {
	router := mux.NewRouter()
	router.StrictSlash(true)

	apiRouter := router.PathPrefix("/api").Subrouter()
	RouteApi(apiRouter, a)

	RouteFrontend(a.Config.Frontend, router, a.Logger)
	middlewares.Apply(router, a)

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
}
