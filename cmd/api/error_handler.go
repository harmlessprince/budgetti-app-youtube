package main

import (
	"errors"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/app_errors"
	"github.com/labstack/echo/v4"
	"strings"
)

func (app *Application) customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)
	if _, ok := err.(*echo.HTTPError); ok {
		//,"message":"code=404, message=Not Found" route not found error
		isRouteNotFound := strings.Contains(err.Error(), "message=Not Found") && strings.Contains(err.Error(), "code=404")
		if isRouteNotFound {
			common.SendNotFoundResponse(c, "Route Not Found")
			return
		}
	}

	if errors.Is(err, app_errors.NewNotFoundError(err.Error())) {
		common.SendNotFoundResponse(c, err.Error())
		return
	}
	common.SendInternalServerErrorResponse(c, err.Error())
	app.server.DefaultHTTPErrorHandler(err, c)
}
