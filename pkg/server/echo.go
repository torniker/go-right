package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/torniker/go-right/env"
)

func Echo() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(logLevel())
	e.HTTPErrorHandler = ErrorHandler
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10,
		LogLevel:  log.ERROR,
	}))
	return e
}

func logLevel() log.Lvl {
	switch env.String() {
	case env.Local, env.Staging, env.Testing:
		return log.DEBUG
	}
	return log.WARN
}

type errorResponse struct {
	Message string `json:"message"`
}


// Success wraps payoad in data field and responds with corresponding json
func Success(c echo.Context, payload interface{}) error {
	return c.JSON(http.StatusOK, struct {
		Data interface{} `json:"data"`
	}{Data: payload})
}
