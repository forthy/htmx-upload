package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// echo.Context -> error
// Provides a health check
func HealthCheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
