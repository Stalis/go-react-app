package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Stalis/go-react-app/server/config"
	"github.com/Stalis/go-react-app/server/router"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found!")
	}
}

func main() {
	conf := config.New()

	if conf.Common.IsDebug {
		fmt.Println("================================================================")
		fmt.Println("=            DEBUG MODE!!! DON'T USE FOR PRODUCTION            =")
		fmt.Println("================================================================")
	}

	l := log.New(os.Stdout, conf.Common.LoggerPrefix+" ", log.LstdFlags)

	r := router.CreateNew(conf, l)

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", conf.HttpServer.Host, conf.HttpServer.Port),

		WriteTimeout: conf.HttpServer.WriteTimeout,
		ReadTimeout:  conf.HttpServer.ReadTimeout,
		IdleTimeout:  conf.HttpServer.IdleTimeout,
		Handler:      r,
	}

	// запускаем сервер в горутине, чтобы не блокировать его
	go func() {
		l.Printf("Server started at: http://%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			l.Println(err)
		}
	}()

	go func() {
		response, err := http.Get("https://google.com")
		if err != nil {
			l.Printf("Test HTTPS Request Failed: %v\n", err.Error())
			return
		}
		l.Printf("Test HTTPS Request passed!: %v\n", response.Status)
	}()

	sigChan := make(chan os.Signal, 1)
	// будем ждать любой сигнал на прерывание процесса
	signal.Notify(sigChan, os.Interrupt)

	// блокируем поток, пока не придет нужный сигнал
	sig := <-sigChan
	l.Printf("Interrupted with signal: %v", sig)

	// завершаем все рабочие процессы контекста
	ctx, cancel := context.WithTimeout(context.Background(), conf.HttpServer.ShutdownWait)
	defer cancel()

	srv.Shutdown(ctx)

	l.Println("shutting down")
	os.Exit(0)
}
