package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	r := CreateRouter(*conf)

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", conf.HttpServer.Host, conf.HttpServer.Port),

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// запускаем сервер в горутине, чтобы не блокировать его
	go func() {
		fmt.Printf("Server started at: http://%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
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

func CreateRouter(conf config.Config) http.Handler {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	RouteApi(apiRouter)

	RouteFrontend(conf, router)

	router.Use(middlewares.LoggingMiddleware)

	return router
}

func RouteFrontend(conf config.Config, router *mux.Router) {
	fileServer := http.FileServer(http.Dir(conf.Frontend.PathToDist))
	router.PathPrefix("/asset-manifest.json").Handler(http.StripPrefix("/", fileServer))
	router.PathPrefix("/favicon.ico").Handler(http.StripPrefix("/", fileServer))
	router.PathPrefix("/logo192.png").Handler(http.StripPrefix("/", fileServer))
	router.PathPrefix("/logo512.png").Handler(http.StripPrefix("/", fileServer))
	router.PathPrefix("/manifest.json").Handler(http.StripPrefix("/", fileServer))
	router.PathPrefix("/robots.txt").Handler(http.StripPrefix("/", fileServer))

	staticRouter := router.PathPrefix("/static").Subrouter()
	staticRouter.PathPrefix("/js").Handler(http.StripPrefix("/", fileServer))
	staticRouter.PathPrefix("/css").Handler(http.StripPrefix("/", fileServer))
	staticRouter.PathPrefix("/media").Handler(http.StripPrefix("/", fileServer))

	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.ServeFile(rw, r, conf.Frontend.PathToDist+"/"+conf.Frontend.IndexPath)
	})
}

func RouteApi(r *mux.Router) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/hello", HelloApiHandler)
}

func HelloApiHandler(w http.ResponseWriter, r *http.Request) {
	type HelloResponse struct {
		Message string `json:"message"`
	}

	response := &HelloResponse{
		Message: "Hello, World!",
	}

	formatted, _ := json.Marshal(response)
	w.Header().Add("Content-Type", "application/json")
	w.Write(formatted)
}
