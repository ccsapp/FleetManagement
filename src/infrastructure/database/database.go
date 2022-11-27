package database

import (
	"PFleetManagement/logic/errors"
	"PFleetManagement/logic/model"
)

type Database struct {
	database map[model.FleetID][]model.Vin
}

func OpenDatabase() Database {
	return Database{database: make(map[model.FleetID][]model.Vin)}
}

func (db Database) AddFleet(fleetId model.FleetID) error {
	if db.database[fleetId] != nil {
		return errors.ErrFleetAlreadyExists
	}
	db.database[fleetId] = []model.Vin{}
	return nil
}

func (db Database) AddCarToFleet(fleetId model.FleetID, vin model.Vin) error {
	if db.database[fleetId] == nil {
		return errors.ErrFleetNotFound
	}
	index := db.getIndexOfCar(fleetId, vin)
	if index != -1 {
		return errors.ErrCarAlreadyInFleet
	}
	db.database[fleetId] = append(db.database[fleetId], vin)
	return nil
}

func (db Database) GetCarsForFleet(fleetId model.FleetID) ([]model.Vin, error) {
	if db.database[fleetId] == nil {
		return nil, errors.ErrFleetNotFound
	}
	return db.database[fleetId], nil
}

func (db Database) RemoveCarFromFleet(fleetId model.FleetID, vin model.Vin) error {
	if db.database[fleetId] == nil {
		return errors.ErrFleetNotFound
	}
	index := db.getIndexOfCar(fleetId, vin)
	if index == -1 {
		return errors.ErrCarNotInFleet
	}
	db.database[fleetId] = append(db.database[fleetId][:index], db.database[fleetId][index+1:]...)
	return nil
}

func (db Database) getIndexOfCar(fleetId model.FleetID, vin model.Vin) int {
	for i, vinSearch := range db.database[fleetId] {
		if vinSearch == vin {
			return i
		}
	}
	return -1
}
