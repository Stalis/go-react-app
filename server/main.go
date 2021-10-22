package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/Stalis/go-react-app/server/config"
	"github.com/Stalis/go-react-app/server/router"
	"github.com/phuslu/log"
)

func main() {
	conf := config.New()

	if conf.Common.IsDebug {
		fmt.Println("================================================================")
		fmt.Println("=            DEBUG MODE!!! DON'T USE FOR PRODUCTION            =")
		fmt.Println("================================================================")
	}

	r := router.CreateNew(conf)

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", conf.HttpServer.Host, conf.HttpServer.Port),

		WriteTimeout: conf.HttpServer.WriteTimeout,
		ReadTimeout:  conf.HttpServer.ReadTimeout,
		IdleTimeout:  conf.HttpServer.IdleTimeout,
		Handler:      r,

		ErrorLog: log.DefaultLogger.Std(log.ErrorLevel, nil, "", 0),
	}

	// запускаем сервер в горутине, чтобы не блокировать его
	go func() {
		log.Debug().Msgf("Server started at: http://%s", srv.Addr)
		log.Fatal().Err(srv.ListenAndServe()).Msg("listen failed")
	}()

	go func() {
		response, err := http.Get("https://google.com")
		if err != nil {
			log.Error().Err(err).Msg("Test HTTPS Request Failed")
			return
		}
		log.Debug().Str("httpStatus", response.Status).Msg("Test HTTPS Request passed!")
	}()

	sigChan := make(chan os.Signal, 1)
	// будем ждать любой сигнал на прерывание процесса
	signal.Notify(sigChan, os.Interrupt)

	// блокируем поток, пока не придет нужный сигнал
	sig := <-sigChan
	log.Debug().Str("signal", sig.String()).Msg("Interrupted with signal")

	// завершаем все рабочие процессы контекста
	ctx, cancel := context.WithTimeout(context.Background(), conf.HttpServer.ShutdownWait)
	defer cancel()

	srv.Shutdown(ctx)

	log.Debug().Msg("shutting down")
	os.Exit(0)
}

func ConfigureLogger(cfg *config.Config) {
	if log.IsTerminal(os.Stderr.Fd()) {
		log.DefaultLogger = log.Logger{
			Level:      log.ParseLevel(cfg.Log.Level),
			Caller:     1,
			TimeFormat: "15:04:05",
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				EndWithMessage: true,
			},
		}
	} else {
		log.DefaultLogger = log.Logger{
			Level: log.ParseLevel(cfg.Log.Level),
			Writer: &log.FileWriter{
				Filename:   "logs/main.log",
				MaxSize:    int64(cfg.Log.MaxFileSize),
				MaxBackups: cfg.Log.MaxBackups,
				LocalTime:  true,
			},
		}
	}
}
