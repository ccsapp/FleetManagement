package api

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"

	"PFleetManagement/logic/model"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// GetCarsInFleet Get Overview of All Cars Assigned to the Given Fleet
	// (GET /fleets/{fleetID}/cars)
	GetCarsInFleet(ctx echo.Context, fleetID model.FleetIDParam) error
	// RemoveCar Remove Car From Fleet
	// (DELETE /fleets/{fleetID}/cars/{vin})
	RemoveCar(ctx echo.Context, fleetID model.FleetIDParam, vin model.VinParam) error
	// GetCar Get Status of the Car With the Given VIN Assigned to the Given Fleet
	// (GET /fleets/{fleetID}/cars/{vin})
	GetCar(ctx echo.Context, fleetID model.FleetIDParam, vin model.VinParam) error
	// AddCarToFleet Add a Car to the Fleet
	// (PUT /fleets/{fleetID}/cars/{vin})
	AddCarToFleet(ctx echo.Context, fleetID model.FleetIDParam, vin model.VinParam) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetCarsInFleet converts echo context to params.
func (w *ServerInterfaceWrapper) GetCarsInFleet(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "fleetID" -------------
	var fleetID model.FleetIDParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "fleetID", runtime.ParamLocationPath, ctx.Param("fleetID"), &fleetID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fleetID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCarsInFleet(ctx, fleetID)
	return err
}

// RemoveCar converts echo context to params.
func (w *ServerInterfaceWrapper) RemoveCar(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "fleetID" -------------
	var fleetID model.FleetIDParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "fleetID", runtime.ParamLocationPath, ctx.Param("fleetID"), &fleetID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fleetID: %s", err))
	}

	// ------------- Path parameter "vin" -------------
	var vin model.VinParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "vin", runtime.ParamLocationPath, ctx.Param("vin"), &vin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter vin: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RemoveCar(ctx, fleetID, vin)
	return err
}

// GetCar converts echo context to params.
func (w *ServerInterfaceWrapper) GetCar(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "fleetID" -------------
	var fleetID model.FleetIDParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "fleetID", runtime.ParamLocationPath, ctx.Param("fleetID"), &fleetID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fleetID: %s", err))
	}

	// ------------- Path parameter "vin" -------------
	var vin model.VinParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "vin", runtime.ParamLocationPath, ctx.Param("vin"), &vin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter vin: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCar(ctx, fleetID, vin)
	return err
}

// AddCarToFleet converts echo context to params.
func (w *ServerInterfaceWrapper) AddCarToFleet(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "fleetID" -------------
	var fleetID model.FleetIDParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "fleetID", runtime.ParamLocationPath, ctx.Param("fleetID"), &fleetID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fleetID: %s", err))
	}

	// ------------- Path parameter "vin" -------------
	var vin model.VinParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "vin", runtime.ParamLocationPath, ctx.Param("vin"), &vin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter vin: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AddCarToFleet(ctx, fleetID, vin)
	return err
}

// EchoRouter
// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// RegisterHandlersWithBaseURL
// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/fleets/:fleetID/cars", wrapper.GetCarsInFleet)
	router.DELETE(baseURL+"/fleets/:fleetID/cars/:vin", wrapper.RemoveCar)
	router.GET(baseURL+"/fleets/:fleetID/cars/:vin", wrapper.GetCar)
	router.PUT(baseURL+"/fleets/:fleetID/cars/:vin", wrapper.AddCarToFleet)

}
