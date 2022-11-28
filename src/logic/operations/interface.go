package operations

import (
	"PFleetManagement/logic/model"
	"context"
)

type Interface interface {
	GetCarsInFleet(ctx context.Context, fleetID model.FleetID) ([]model.CarBase, error)
	RemoveCar(ctx context.Context, fleetID model.FleetID, vin model.Vin) error
	GetCar(ctx context.Context, fleetID model.FleetID, vin model.Vin) (*model.Car, error)
	AddCarToFleet(ctx context.Context, fleetID model.FleetID, vin model.Vin) (*model.CarBase, error)
}
