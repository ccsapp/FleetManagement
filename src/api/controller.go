// Package api contains the request handlers to be connected to the echo router
package api

import (
	"PFleetManagement/logic/errors"
	"PFleetManagement/logic/model"
	"PFleetManagement/logic/operations"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Controller The implementation of the request handlers. Delegates parsed (by ServerInterfaceWrapper)
// and validated requests to operations.IOperations implementation and writes the operation's return
// value as json response with the correct response code.
// Errors are simply returned to be later handled by the HTTPErrorHandler of the echo.Echo instance.
type Controller struct {
	operations operations.IOperations
}

func NewController(operations operations.IOperations) Controller {
	return Controller{
		operations,
	}
}

// extractRequestContext Get a context.Context from an echo.Context which is bound to the user request
// and thus cancelled when the HTTP request is closed before it returns
func extractRequestContext(ctx echo.Context) context.Context {
	return ctx.Request().Context()
}

func (c Controller) GetCarsInFleet(ctx echo.Context, fleetID model.FleetIDParam) error {
	if !model.IsFleetIDValid(fleetID) {
		return errors.ErrInvalidFleetId
	}

	cars, err := c.operations.GetCarsInFleet(extractRequestContext(ctx), fleetID)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, cars)
}

func (c Controller) RemoveCar(ctx echo.Context, fleetID model.FleetIDParam, vin model.VinParam) error {
	if !model.IsFleetIDValid(fleetID) {
		return errors.ErrInvalidFleetId
	}
	if !model.IsVinValid(vin) {
		return errors.ErrInvalidVin
	}

	err := c.operations.RemoveCar(extractRequestContext(ctx), fleetID, vin)

	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (c Controller) GetCar(ctx echo.Context, fleetID model.FleetIDParam, vin model.VinParam) error {
	if !model.IsFleetIDValid(fleetID) {
		return errors.ErrInvalidFleetId
	}
	if !model.IsVinValid(vin) {
		return errors.ErrInvalidVin
	}

	car, err := c.operations.GetCar(extractRequestContext(ctx), fleetID, vin)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, car)
}

func (c Controller) AddCarToFleet(ctx echo.Context, fleetID model.FleetIDParam, vin model.VinParam) error {
	if !model.IsFleetIDValid(fleetID) {
		return errors.ErrInvalidFleetId
	}
	if !model.IsVinValid(vin) {
		return errors.ErrInvalidVin
	}

	car, err := c.operations.AddCarToFleet(extractRequestContext(ctx), fleetID, vin)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, car)
}
