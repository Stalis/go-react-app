package app

import (
	"go-react-app/internal/config"
	"go-react-app/internal/dal"
	"go-react-app/internal/util/logger"
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
