package main

import (
	"github.com/harmlessprince/bougette-backend/cmd/api/handlers"
	"github.com/harmlessprince/bougette-backend/cmd/api/middlewares"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/mailer"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	logger        echo.Logger
	server        *echo.Echo
	handler       handlers.Handler
	appMiddleware middlewares.AppMiddleware
}

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	db, err := common.NewMysql()

	if err != nil {
		e.Logger.Fatal("Error loading .env file")
	}

	appMailer := mailer.NewMailer(e.Logger)
	h := handlers.Handler{
		DB:     db,
		Logger: e.Logger,
		Mailer: appMailer,
	}
	appMiddleware := middlewares.AppMiddleware{
		DB:     db,
		Logger: e.Logger,
	}
	app := Application{
		logger:        e.Logger,
		server:        e,
		handler:       h,
		appMiddleware: appMiddleware,
	}
	e.Use(middleware.Logger())
	e.Use(middlewares.AnotherMiddleware)
	e.Use(middlewares.CustomMiddleware)
	err = app.serve()
	if err != nil {
		app.logger.Fatal(err)
	}
}

//go get github.com/labstack/echo/v4/middleware
//go get -u github.com/labstack/echo/v4/middleware
