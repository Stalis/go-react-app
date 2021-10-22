package router

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Stalis/go-react-app/server/config"
	"github.com/Stalis/go-react-app/server/dal"
	"github.com/Stalis/go-react-app/server/handlers/account"
	"github.com/Stalis/go-react-app/server/middlewares"
	"github.com/gorilla/mux"
)

func CreateNew(conf *config.Config, l *log.Logger) http.Handler {
	router := mux.NewRouter()
	router.StrictSlash(true)

	apiRouter := router.PathPrefix("/api").Subrouter()
	RouteApi(apiRouter, conf, l)

	RouteFrontend(conf.Frontend, router, l)
	middlewares.Apply(router, l)

	return router
}

func RouteFrontend(conf config.FrontendConfig, router *mux.Router, l *log.Logger) {
	files, err := ioutil.ReadDir(conf.PathToDist)
	if err != nil {
		l.Println(err)
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

		l.Printf("Static route for %s\n", prefix)
		router.PathPrefix(prefix).Handler(http.StripPrefix("/", fileServer))
	}

	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.ServeFile(rw, r, conf.PathToDist+"/"+conf.IndexPath)
	})
}

func RouteApi(r *mux.Router, cfg *config.Config, l *log.Logger) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	db, err := dal.ConnectDB(cfg.Database.Url, l)
	if err != nil {
		l.Println(err.Error())
		os.Exit(-1)
	}

	accountRouter := r.PathPrefix("/account").Subrouter()
	accountRouter.Handle("/login", account.NewLogin(db))
	accountRouter.Handle("/register", account.NewRegister(db))
}
