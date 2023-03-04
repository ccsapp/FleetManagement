package main

import (
	"PFleetManagement/api"
	"PFleetManagement/environment"
	"PFleetManagement/infrastructure/database"
	"PFleetManagement/infrastructure/dcar"
	rentalManagement "PFleetManagement/infrastructure/rentalmanagement"
	"PFleetManagement/logic/fleetErrors"
	"PFleetManagement/logic/operations"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

// newApp allows production as well as testing to create a new Echo instance for the API.
// Configuration values are read from the environment with environment.GetEnvironment().
func newApp(fleetDb database.FleetDB) (*echo.Echo, error) {
	e := echo.New()
	e.HTTPErrorHandler = api.FleetErrorHandler

	allowOrigins := environment.GetEnvironment().GetAllowOrigins()

	if len(allowOrigins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: allowOrigins,
		}))
	}

	// validate incoming requests against the OpenAPI spec
	err := api.AddOpenApiValidationMiddleware(e)
	if err != nil {
		return nil, err
	}

	err = fleetDb.AddFleet(context.TODO(), "xk48jpgz") // TODO manage fleets correctly
	if err != nil && !errors.Is(err, fleetErrors.ErrFleetAlreadyExists) {
		return nil, err
	}

	requestTimeout := environment.GetEnvironment().GetRequestTimeout()

	dcarClient, err := dcar.NewClientWithResponses(
		environment.GetEnvironment().GetCarServerUrl(),
		dcar.WithHTTPClient(&http.Client{Timeout: requestTimeout}),
	)

	rmClient, err := rentalManagement.NewClientWithResponses(
		environment.GetEnvironment().GetRentalServerUrl(),
		rentalManagement.WithHTTPClient(&http.Client{Timeout: requestTimeout}),
	)

	if err != nil {
		return nil, err
	}

	operationsInstance := operations.NewOperations(fleetDb, dcarClient, rmClient)
	controllerInstance := api.NewController(operationsInstance)

	api.RegisterHandlers(e, controllerInstance)

	return e, nil
}

func main() {
	var fleetDb database.FleetDB
	fleetDb, err := database.OpenDatabase(environment.GetEnvironment())
	if err != nil {
		log.Fatal(err)
	}

	var e *echo.Echo
	e, err = newApp(fleetDb)
	if err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", environment.GetEnvironment().GetAppExposePort())))
}
