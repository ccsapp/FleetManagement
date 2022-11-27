package operations

import (
	"PFleetManagement/logic/model"
	openapiTypes "github.com/deepmap/oapi-codegen/pkg/types"
	"time"
)

type Operations struct {
	// TODO add client field here
}

func (o Operations) GetCarsInFleet(fleetID model.FleetID) (*[]model.CarBase, error) {
	// TODO implement me
	return &[]model.CarBase{
		{
			Brand: "Audi",
			Model: "A3",
			ProductionDate: openapiTypes.Date{
				Time: time.Date(2017, 7, 21, 12, 0, 0, 0, time.UTC),
			},
			Vin: "WDD1690071J236589",
		},
	}, nil
}

func (o Operations) RemoveCar(fleetID model.FleetID, vin model.Vin) error {
	// TODO implement me
	return nil
}

func (o Operations) GetCar(fleetID model.FleetID, vin model.Vin) (*model.Car, error) {
	// TODO implement me
	return &model.Car{
		Brand: "Audi",
		DynamicData: model.DynamicData{
			DoorsLockState:      model.LOCKED,
			EngineState:         model.ON,
			FuelLevelPercentage: 100,
			Position:            model.DynamicDataPosition{},
			TrunkLockState:      model.LOCKED,
		},
		Model: "A3",
		ProductionDate: openapiTypes.Date{
			Time: time.Date(2017, 7, 21, 12, 0, 0, 0, time.UTC),
		},
		TechnicalSpecification: model.TechnicalSpecification{
			Color:         "black",
			Consumption:   model.TechnicalSpecificationConsumption{},
			Emissions:     model.TechnicalSpecificationEmissions{},
			Engine:        model.TechnicalSpecificationEngine{},
			Fuel:          model.PETROL,
			FuelCapacity:  "54.0L;85.2kWh",
			NumberOfDoors: 3,
			NumberOfSeats: 5,
			Tire:          model.TechnicalSpecificationTire{},
			Transmission:  model.AUTOMATIC,
			TrunkVolume:   435,
			Weight:        1320,
		},
		Vin: "WDD1690071J236589",
	}, nil
}

func (o Operations) AddCarToFleet(fleetID model.FleetID, vin model.Vin) (*model.CarBase, error) {
	// TODO implement me
	return &model.CarBase{
		Brand: "Audi",
		Model: "A3",
		ProductionDate: openapiTypes.Date{
			Time: time.Date(2017, 7, 21, 12, 0, 0, 0, time.UTC),
		},
		Vin: "WDD1690071J236589",
	}, nil
}
