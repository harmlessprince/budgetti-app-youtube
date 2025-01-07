package main

import (
	"errors"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/app_errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func (app *Application) customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)
	var validationErrors []*common.ValidationError
	if ve, ok := err.(*echo.HTTPError); ok {
		//,"message":"code=404, message=Not Found" route not found error
		isRouteNotFound := strings.Contains(err.Error(), "message=Not Found") && strings.Contains(err.Error(), "code=404")
		if isRouteNotFound {
			common.SendNotFoundResponse(c, "Route Not Found")
			return
		}
		if ve.Code == http.StatusUnprocessableEntity {
			// Assume the message of the error contains validation error details
			if errs, ok := ve.Message.([]*common.ValidationError); ok {
				validationErrors = errs
			} else {
				// Fallback for unexpected error format
				validationErrors = []*common.ValidationError{
					{
						Error:     "Unexpected error format",
						Key:       "unknown",
						Condition: "unknown",
					},
				}
			}
		}
		_ = common.SendFailedValidationResponse(c, validationErrors)
		return
	}

	if errors.Is(err, app_errors.NewNotFoundError(err.Error())) {
		common.SendNotFoundResponse(c, err.Error())
		return
	}
	common.SendInternalServerErrorResponse(c, err.Error())
	app.server.DefaultHTTPErrorHandler(err, c)
}
