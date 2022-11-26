package main

import (
	"PFleetManagement/logic/operations"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"PFleetManagement/api"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*.cloud.iai.kit.edu"},
		AllowMethods: []string{
			http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete,
		},
	}))

	operationsInstance := operations.Operations{}
	controllerInstance := api.NewController(operationsInstance)

	api.RegisterHandlers(e, controllerInstance)

	e.Logger.Fatal(e.Start(":80"))
}
