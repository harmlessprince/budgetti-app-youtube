package main

import (
	"fmt"
	"os"
)

func (app *Application) serve() error {
	app.routes(app.handler)
	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf("localhost:%s", port)
	app.server.HTTPErrorHandler = app.customHTTPErrorHandler
	err := app.server.Start(appAddress)
	if err != nil {
		return err
	}
	return nil
}
