package middlewares

import (
	"go-react-app/server/app"

	"github.com/gorilla/mux"
)

func Apply(router *mux.Router, a *app.App) {
	router.Use(
		NewRecovery(a.Logger).Middleware,
		NewRequestLogger(a.Logger).Middleware,
		NewAuthentication(a.Logger, a.DbContext).Middleware,
	)
}
