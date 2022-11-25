package main

import (
	"github.com/labstack/echo/v4"

	"PFleetManagement/api"
)

func main() {
	e := echo.New()

	// TODO insert handler implementation
	api.RegisterHandlers(e, nil)

	e.Logger.Fatal(e.Start(":80"))
}
