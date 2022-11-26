package api

import (
	"PFleetManagement/logic/model"
	"PFleetManagement/logic/operations"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	operations operations.Operations
}

func NewController(operations operations.Operations) Controller {
	return Controller{
		operations,
	}
}

func (c Controller) GetCarsInFleet(ctx echo.Context, fleetID model.FleetIDParam) error {
	cars, err := c.operations.GetCarsInFleet(fleetID)

	if err != nil { // TODO do real error handling/forwarding
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to fetch data")
	}

	return ctx.JSON(http.StatusOK, cars)
}

func (c Controller) RemoveCar(ctx echo.Context, fleetID model.FleetIDParam, vin model.VinParam) error {
	err := c.operations.RemoveCar(fleetID, vin)

	if err != nil { // TODO do real error handling/forwarding
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to perform remove")
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (c Controller) GetCar(ctx echo.Context, fleetID model.FleetIDParam, vin model.VinParam) error {
	car, err := c.operations.GetCar(fleetID, vin)

	if err != nil { // TODO do real error handling/forwarding
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to fetch car")
	}

	return ctx.JSON(http.StatusOK, car)
}

func (c Controller) AddCarToFleet(ctx echo.Context, fleetID model.FleetIDParam, vin model.VinParam) error {
	car, err := c.operations.AddCarToFleet(fleetID, vin)

	if err != nil { // TODO do real error handling/forwarding, remember 204 if already assigned
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to perform add")
	}

	return ctx.JSON(http.StatusOK, car)
}
