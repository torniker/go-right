package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/torniker/go-right/app/country"
	"github.com/torniker/go-right/pkg/server"
)


func New(e *echo.Echo) {
	countryGroup(e.Group("/country"))
}


func countryGroup(g *echo.Group) {
	service := country.New()
	// handlers
	g.GET("", countryList(service))
}


func countryList(service *country.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		countries, err := service.List()
		if err != nil {
			return err
		}
		return server.Success(c, countries)
	}
}
