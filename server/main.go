package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Stalis/go-react-app/api"
	"github.com/Stalis/go-react-app/config"
	"github.com/Stalis/go-react-app/middlewares"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found!")
	}
}

func main() {
	conf := config.New()

	if conf.IsDebug {
		fmt.Println("================================================================")
		fmt.Println("=            DEBUG MODE!!! DON'T USE FOR PRODUCTION            =")
		fmt.Println("================================================================")
	}

	r := CreateRouter(conf)

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", conf.HttpServer.Host, conf.HttpServer.Port),

		WriteTimeout: conf.HttpServer.WriteTimeout,
		ReadTimeout:  conf.HttpServer.ReadTimeout,
		IdleTimeout:  conf.HttpServer.IdleTimeout,
		Handler:      r,
	}

	// запускаем сервер в горутине, чтобы не блокировать его
	go func() {
		log.Printf("Server started at: http://%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		response, err := http.Get("https://google.com")
		if err != nil {
			log.Printf("Test HTTPS Request Failed: %v\n", err.Error())
			return
		}
		log.Printf("Test HTTPS Request passed!: %v\n", response.Status)
	}()

	c := make(chan os.Signal, 1)
	// будем ждать любой сигнал на прерывание процесса
	signal.Notify(c, os.Interrupt)

	// блокируем поток, пока не придет нужный сигнал
	<-c

	// завершаем все рабочие процессы контекста
	ctx, cancel := context.WithTimeout(context.Background(), conf.HttpServer.ShutdownWait)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

func CreateRouter(conf *config.Config) http.Handler {
	router := mux.NewRouter()
	router.StrictSlash(true)
	apiRouter := router.PathPrefix("/api").Subrouter()
	api.Route(apiRouter, conf)

	RouteFrontend(conf.Frontend, router)
	middlewares.Apply(router)

	return router
}

func RouteFrontend(conf config.FrontendConfig, router *mux.Router) {
	files, err := ioutil.ReadDir(conf.PathToDist)
	if err != nil {
		log.Println(err)
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

		log.Printf("Static route for %s\n", prefix)
		router.PathPrefix(prefix).Handler(http.StripPrefix("/", fileServer))
	}

	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.ServeFile(rw, r, conf.PathToDist+"/"+conf.IndexPath)
	})
}
