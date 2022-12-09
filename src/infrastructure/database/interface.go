package database

import (
	"PFleetManagement/logic/model"
	"context"
)

// FleetDB Abstraction over database backends to manage car-fleet assignment.
// Returns errors as defined in logic/operations
type FleetDB interface {
	AddFleet(ctx context.Context, fleetId model.FleetID) error
	AddCarToFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) error
	RemoveCarFromFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) error
	GetCarsForFleet(ctx context.Context, fleetId model.FleetID) ([]model.Vin, error)
	IsCarInFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) (bool, error)
}
