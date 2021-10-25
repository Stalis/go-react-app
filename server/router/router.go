package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"go-react-app/server/config"
	"go-react-app/server/dal"
	"go-react-app/server/handlers/account"
	"go-react-app/server/middlewares"
	"go-react-app/server/util/logger"

	"github.com/gorilla/mux"
	"github.com/phuslu/log"
)

func New(conf *config.Config, log *logger.Logger) http.Handler {
	router := mux.NewRouter()
	router.StrictSlash(true)

	apiRouter := router.PathPrefix("/api").Subrouter()
	RouteApi(apiRouter, conf)

	RouteFrontend(conf.Frontend, router)
	middlewares.Apply(router, log)

	return router
}

func RouteFrontend(conf config.FrontendConfig, router *mux.Router) {
	files, err := ioutil.ReadDir(conf.PathToDist)
	if err != nil {
		log.Error().Err(err).Msg("")
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

func RouteApi(r *mux.Router, cfg *config.Config) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	db, err := dal.ConnectDB(cfg.Database.Url)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to DB")
		os.Exit(-1)
	}

	accountRouter := r.PathPrefix("/account").Subrouter()
	accountRouter.Handle("/login", account.NewLogin(db))
	accountRouter.Handle("/register", account.NewRegister(db))
}
