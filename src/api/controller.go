package api

import (
	"PFleetManagement/logic/errors"
	"PFleetManagement/logic/model"
	"PFleetManagement/logic/operations"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	operations operations.Interface
}

func NewController(operations operations.Interface) Controller {
	return Controller{
		operations,
	}
}

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
	} else if !model.IsVinValid(vin) {
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
	} else if !model.IsVinValid(vin) {
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
	} else if !model.IsVinValid(vin) {
		return errors.ErrInvalidVin
	}

	car, err := c.operations.AddCarToFleet(extractRequestContext(ctx), fleetID, vin)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, car)
}
