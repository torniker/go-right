package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/torniker/go-right/api/country/handler"
	"github.com/torniker/go-right/env"
	"github.com/torniker/go-right/pkg/server"
)

func main() {
	// setup environment
	env.New("country")
	basePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	env.SetForLocal(env.BasePathKey, basePath)
	env.SetForLocal(env.PortKey, ":5656")
	env.SetForLocal(env.MainGoPathKey, "api/country")

	// setup echo
	e := server.Echo()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/bad-request", func(c echo.Context) error {
		return server.ErrBadRequest("Bad request")
	})
	e.GET("/internal", func(c echo.Context) error {
		return fmt.Errorf("internal error happened")
	})

	// initialize handlers
	handler.New(e)

	// start the web server
	server.Start(e, env.IsLocal())
}
