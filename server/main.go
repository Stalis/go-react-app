package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"go-react-app/server/config"
	"go-react-app/server/router"
	"go-react-app/server/util/logger"
)

func main() {
	conf := config.New()

	if conf.Common.IsDebug {
		fmt.Println("================================================================")
		fmt.Println("=            DEBUG MODE!!! DON'T USE FOR PRODUCTION            =")
		fmt.Println("================================================================")
	}

	logger := logger.New(conf.Common.IsDebug)
	r := router.New(conf, logger)

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", conf.HttpServer.Host, conf.HttpServer.Port),

		WriteTimeout: conf.HttpServer.WriteTimeout,
		ReadTimeout:  conf.HttpServer.ReadTimeout,
		IdleTimeout:  conf.HttpServer.IdleTimeout,
		Handler:      r,
	}

	// запускаем сервер в горутине, чтобы не блокировать его
	go func() {
		logger.Info().Msgf("Server started at: http://%s", srv.Addr)
		logger.Fatal().Err(srv.ListenAndServe()).Msg("listen failed")
	}()

	go func() {
		response, err := http.Get("https://google.com")
		if err != nil {
			logger.Error().Err(err).Msg("Test HTTPS Request Failed")
			return
		}
		logger.Debug().Str("httpStatus", response.Status).Msg("Test HTTPS Request passed!")
	}()

	sigChan := make(chan os.Signal, 1)
	// будем ждать любой сигнал на прерывание процесса
	signal.Notify(sigChan, os.Interrupt)

	// блокируем поток, пока не придет нужный сигнал
	sig := <-sigChan
	logger.Debug().Str("signal", sig.String()).Msg("Interrupted with signal")

	// завершаем все рабочие процессы контекста
	ctx, cancel := context.WithTimeout(context.Background(), conf.HttpServer.ShutdownWait)
	defer cancel()

	srv.Shutdown(ctx)

	logger.Debug().Msg("shutting down")
	os.Exit(0)
}
