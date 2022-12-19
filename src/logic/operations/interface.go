// Package operations contains the task processes which can be triggered via this service's HTTP API
package operations

//go:generate mockgen -source=interface.go -package=mocks -destination=../../mocks/operationsinterface_mock.go

import (
	"PFleetManagement/logic/model"
	"context"
)

// The IOperations defines the interface of the operations provided to the controller.
// All blocking operations use the given context.
// Returned errors are either from logic/errors or internal errors from library calls.
type IOperations interface {
	// GetCarsInFleet Get an overview of all cars assigned to the given fleet
	GetCarsInFleet(ctx context.Context, fleetID model.FleetID) ([]model.CarBase, error)

	// RemoveCar Remove the given car from the given fleet
	RemoveCar(ctx context.Context, fleetID model.FleetID, vin model.Vin) error

	// GetCar Get data and status of the given car assigned to the given fleet
	GetCar(ctx context.Context, fleetID model.FleetID, vin model.Vin) (*model.Car, error)

	// AddCarToFleet Add (assign) the given car to the given fleet
	AddCarToFleet(ctx context.Context, fleetID model.FleetID, vin model.Vin) (*model.CarBase, error)
}
