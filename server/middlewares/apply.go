package middlewares

import (
	"go-react-app/server/middlewares/requestlog"
	"go-react-app/server/util/logger"

	"github.com/gorilla/mux"
)

func Apply(router *mux.Router, log *logger.Logger) {

	recovery := &recovery{log}
	router.Use(recovery.Middleware)

	logging := &logging{log}
	router.Use(logging.Middleware)

	requestlogger := requestlog.New(log)
	router.Use(requestlogger.Middleware)
}
