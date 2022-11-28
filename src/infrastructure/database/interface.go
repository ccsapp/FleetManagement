package database

import "PFleetManagement/logic/model"

// FleetDB Abstraction over database backends to manage car-fleet assignment.
// Returns errors as defined in logic/operations
type FleetDB interface {
	AddCarToFleet(fleetId model.FleetID, vin model.Vin) error
	RemoveCarFromFleet(fleetId model.FleetID, vin model.Vin) error
	GetCarsForFleet(fleetId model.FleetID) ([]model.Vin, error)
	IsCarInFleet(fleetId model.FleetID, vin model.Vin) (bool, error)
}
