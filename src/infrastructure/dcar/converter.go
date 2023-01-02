package dcar

import (
	"PFleetManagement/logic/model"
	carTypes "git.scc.kit.edu/cm-tm/cm-team/projectwork/pse/domain/d-cargotypes.git"
)

func toModelFromTechnicalSpecification(t *carTypes.TechnicalSpecification) model.TechnicalSpecification {
	return model.TechnicalSpecification{
		Color: t.Color,
		Consumption: model.TechnicalSpecificationConsumption{
			City:     t.Consumption.City,
			Combined: t.Consumption.Combined,
			Overland: t.Consumption.Overland,
		},
		Emissions: model.TechnicalSpecificationEmissions{
			City:     t.Emissions.City,
			Combined: t.Emissions.Combined,
			Overland: t.Emissions.Overland,
		},
		Engine: model.TechnicalSpecificationEngine{
			Power: t.Engine.Power,
			Type:  t.Engine.Type,
		},
		Fuel:          model.TechnicalSpecificationFuel(t.Fuel),
		FuelCapacity:  t.FuelCapacity,
		NumberOfDoors: t.NumberOfDoors,
		NumberOfSeats: t.NumberOfSeats,
		Tire: model.TechnicalSpecificationTire{
			Manufacturer: t.Tire.Manufacturer,
			Type:         t.Tire.Type,
		},
		Transmission: model.TechnicalSpecificationTransmission(t.Transmission),
		TrunkVolume:  t.TrunkVolume,
		Weight:       t.Weight,
	}
}

func toModelFromDynamicData(d *carTypes.DynamicData) model.DynamicData {
	return model.DynamicData{
		DoorsLockState:      model.LockState(d.DoorsLockState),
		EngineState:         model.DynamicDataEngineState(d.EngineState),
		FuelLevelPercentage: d.FuelLevelPercentage,
		Position: model.DynamicDataPosition{
			Latitude:  d.Position.Latitude,
			Longitude: d.Position.Longitude,
		},
		TrunkLockState: model.LockState(d.TrunkLockState),
	}
}

// ToModelFromCar deep-copies a carTypes.Car to a model.Car
func ToModelFromCar(c *carTypes.Car) model.Car {
	return model.Car{
		Brand:                  c.Brand,
		DynamicData:            toModelFromDynamicData(&c.DynamicData),
		Model:                  c.Model,
		ProductionDate:         c.ProductionDate,
		TechnicalSpecification: toModelFromTechnicalSpecification(&c.TechnicalSpecification),
		Vin:                    c.Vin,
	}
}

// ToModelBaseFromCar deep-copies the relevant information from a carTypes.Car to a model.CarBase
func ToModelBaseFromCar(c *carTypes.Car) model.CarBase {
	return model.CarBase{
		Brand:          c.Brand,
		Model:          c.Model,
		ProductionDate: c.ProductionDate,
		Vin:            c.Vin,
	}
}

func newTechnicalSpecificationFromModel(t *model.TechnicalSpecification) carTypes.TechnicalSpecification {
	return carTypes.TechnicalSpecification{
		Color: t.Color,
		Consumption: carTypes.TechnicalSpecificationConsumption{
			City:     t.Consumption.City,
			Combined: t.Consumption.Combined,
			Overland: t.Consumption.Overland,
		},
		Emissions: carTypes.TechnicalSpecificationEmissions{
			City:     t.Emissions.City,
			Combined: t.Emissions.Combined,
			Overland: t.Emissions.Overland,
		},
		Engine: carTypes.TechnicalSpecificationEngine{
			Power: t.Engine.Power,
			Type:  t.Engine.Type,
		},
		Fuel:          carTypes.TechnicalSpecificationFuel(t.Fuel),
		FuelCapacity:  t.FuelCapacity,
		NumberOfDoors: t.NumberOfDoors,
		NumberOfSeats: t.NumberOfSeats,
		Tire: carTypes.TechnicalSpecificationTire{
			Manufacturer: t.Tire.Manufacturer,
			Type:         t.Tire.Type,
		},
		Transmission: carTypes.TechnicalSpecificationTransmission(t.Transmission),
		TrunkVolume:  t.TrunkVolume,
		Weight:       t.Weight,
	}
}

func newDynamicDataFromModel(d *model.DynamicData) carTypes.DynamicData {
	return carTypes.DynamicData{
		DoorsLockState:      carTypes.DynamicDataLockState(d.DoorsLockState),
		EngineState:         carTypes.DynamicDataEngineState(d.EngineState),
		FuelLevelPercentage: d.FuelLevelPercentage,
		Position: carTypes.DynamicDataPosition{
			Latitude:  d.Position.Latitude,
			Longitude: d.Position.Longitude,
		},
		TrunkLockState: carTypes.DynamicDataLockState(d.TrunkLockState),
	}
}

// NewCarFromModel deep-copies a model.Car to a carTypes.Car
func NewCarFromModel(c *model.Car) carTypes.Car {
	return carTypes.Car{
		Brand:                  c.Brand,
		DynamicData:            newDynamicDataFromModel(&c.DynamicData),
		Model:                  c.Model,
		ProductionDate:         c.ProductionDate,
		TechnicalSpecification: newTechnicalSpecificationFromModel(&c.TechnicalSpecification),
		Vin:                    c.Vin,
	}
}
