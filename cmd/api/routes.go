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

	budgetRoutes := apiGroup.Group("/budgets", app.appMiddleware.AuthenticationMiddleware)
	{
		budgetRoutes.POST("/store", handler.CreateBudget)
		budgetRoutes.GET("/all", handler.ListBudget)
		budgetRoutes.PATCH("/update/:id", handler.UpdateBudget)
		budgetRoutes.DELETE("/delete/:id", handler.DeleteBudget)
	}

	walletRoutes := apiGroup.Group("/wallets", app.appMiddleware.AuthenticationMiddleware)
	{
		walletRoutes.POST("/store", handler.CreateWallet)
		walletRoutes.GET("/all", handler.ListWallet)
		walletRoutes.GET("/generate/defaults", handler.GenerateDefaultWallets)
		//budgetRoutes.GET("/all", handler.ListBudget)
		//budgetRoutes.PATCH("/update/:id", handler.UpdateBudget)
		//budgetRoutes.DELETE("/delete/:id", handler.DeleteBudget)
	}

	transactionRoutes := apiGroup.Group("/transactions", app.appMiddleware.AuthenticationMiddleware)
	{
		transactionRoutes.POST("/store", handler.StoreTransaction)
		transactionRoutes.PATCH("/reverse/:id", handler.ReverseTransaction)
		transactionRoutes.GET("/all", handler.ListTransactions)
		//budgetRoutes.DELETE("/delete/:id", handler.DeleteBudget)
	}
	apiGroup.POST("/transfer", handler.Transfer, app.appMiddleware.AuthenticationMiddleware)
	app.server.GET("/", handler.HealthCheck)

}
