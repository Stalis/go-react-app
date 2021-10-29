package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/pkg/errors"

	"go-react-app/internal/app"
	"go-react-app/internal/app/router"
	"go-react-app/internal/config"
	"go-react-app/internal/dal"
	"go-react-app/internal/util/logger"
)

func main() {
	configPath := flag.String("config", ".env", "path to '.env' format config file")
	flag.Parse()

	conf := config.New(configPath)
	logger := logger.New(conf.Common.IsDebug, &conf.Log)

	db, err := dal.ConnectDB(logger, &conf.Database)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("")
		return
	}
	defer db.Close()

	application := app.New(logger, conf, db)
	appRouter := router.New(application)

	if conf.Common.IsDebug {
		fmt.Println("================================================================")
		fmt.Println("=            DEBUG MODE!!! DON'T USE FOR PRODUCTION            =")
		fmt.Println("================================================================")
	}

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", conf.HttpServer.Host, conf.HttpServer.Port),

		WriteTimeout: conf.HttpServer.WriteTimeout,
		ReadTimeout:  conf.HttpServer.ReadTimeout,
		IdleTimeout:  conf.HttpServer.IdleTimeout,
		Handler:      appRouter,
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

func foo() error {
	err := bar()
	if err != nil {
		return err
	}
	return nil
}

func bar() error {
	err := func() error {
		return errors.New("test error")
	}()
	if err != nil {
		return err
	}
	return nil
}
