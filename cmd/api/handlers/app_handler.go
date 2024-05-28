package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type healthCheck struct {
	Health bool `json:"health"`
}

func (h *Handler) HealthCheck(c echo.Context) error {
	healthcheckStruct := healthCheck{
		Health: true,
	}
	return c.JSON(http.StatusOK, healthcheckStruct)
}
