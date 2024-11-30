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
		publicAuthRoutes.POST("/forgot/password", handler.ForgotPasswordHandler)
		publicAuthRoutes.POST("/reset/password", handler.ResetPasswordHandler)
	}

	profileRoutes := apiGroup.Group("/profile", app.appMiddleware.AuthenticationMiddleware)
	{
		profileRoutes.GET("/authenticated/user", handler.GetAuthenticatedUser)
		profileRoutes.PATCH("/change/password", handler.ChangeUserPassword)
	}

	categoryRoutes := apiGroup.Group("/categories", app.appMiddleware.AuthenticationMiddleware)
	{
		categoryRoutes.GET("/user/all", handler.ListUserCategories)
		categoryRoutes.POST("/custom/store", handler.CreateCustomUserCategory)
		categoryRoutes.GET("/all", handler.ListCategories)
		categoryRoutes.POST("/store", handler.CreateCategory)
		categoryRoutes.DELETE("/delete/:id", handler.DeleteCategory)
		categoryRoutes.POST("/associate/user/to/categories", handler.AssociateUserToCategories)
	}

	app.server.GET("/", handler.HealthCheck)

}
