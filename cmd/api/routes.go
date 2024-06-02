package main

import (
	"github.com/harmlessprince/bougette-backend/cmd/api/handlers"
)

func (app *Application) routes(handler handlers.Handler) {
	apiGroup := app.server.Group("/api")
	publicAuthRoutes := apiGroup.Group("/auth")
	{
		publicAuthRoutes.POST("/register", handler.RegisterHandler)
		publicAuthRoutes.POST("/login", handler.LoginHandler)
	}

	profileRoutes := apiGroup.Group("/profile", app.appMiddleware.AuthenticationMiddleware)
	{
		profileRoutes.GET("/authenticated/user", handler.GetAuthenticatedUser)
		profileRoutes.PATCH("/change/password", handler.ChangeUserPassword)
	}

	app.server.GET("/", handler.HealthCheck)

}
