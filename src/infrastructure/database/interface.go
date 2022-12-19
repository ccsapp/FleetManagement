package database

//go:generate mockgen -source=interface.go -package=mocks -destination=../../mocks/database_mock.go

import (
	"PFleetManagement/logic/model"
	"context"
)

// FleetDB Abstraction over database backends to manage car-fleet assignment.
// Returns errors as defined in logic/operations
type FleetDB interface {
	// AddFleet creates a new empty fleet. This is necessary before a car can be assigned to it.
	AddFleet(ctx context.Context, fleetId model.FleetID) error

	// AddCarToFleet adds a reference to the given car (by its VIN) to the given fleet.
	// Fails on unknown fleet or duplicate entry but does not perform further checks on the VIN.
	AddCarToFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) error

	// RemoveCarFromFleet removes the reference to the given car (its VIN) from the given fleet if it is contained
	RemoveCarFromFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) error

	// GetCarsForFleet reads the VINs of the cars which are assigned to the given fleet
	GetCarsForFleet(ctx context.Context, fleetId model.FleetID) ([]model.Vin, error)

	// IsCarInFleet checks whether the given car (identified by its VIN) is assigned to the given fleet
	IsCarInFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) (bool, error)
}
