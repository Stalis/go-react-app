package app

import (
	"go-react-app/server/config"
	"go-react-app/server/dal"
	"go-react-app/server/util/logger"
)

type App struct {
	Logger    *logger.Logger
	Config    *config.Config
	DbContext *dal.DB
}

func New(
	log *logger.Logger,
	conf *config.Config,
	db *dal.DB) *App {
	return &App{log, conf, db}
}
