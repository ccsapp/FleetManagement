package dcar

import (
	"PFleetManagement/logic/model"
)

// TODO validate enum values

func (t *TechnicalSpecification) toModel() model.TechnicalSpecification {
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

func (d *DynamicData) toModel() model.DynamicData {
	return model.DynamicData{
		DoorsLockState:      model.LockState(d.DoorsLockState),
		EngineState:         model.DynamicDataEngineState(d.EngineState),
		FuelLevelPercentage: d.FuelLevelPercentage,
		Position: model.DynamicDataPosition{
			Latitude:  d.Position.Latitude,
			Longitude: d.Position.Latitude,
		},
		TrunkLockState: model.LockState(d.TrunkLockState),
	}
}

func (c *Car) ToModel() model.Car {
	return model.Car{
		Brand:                  c.Brand,
		DynamicData:            c.DynamicData.toModel(),
		Model:                  c.Model,
		ProductionDate:         c.ProductionDate,
		TechnicalSpecification: c.TechnicalSpecification.toModel(),
		Vin:                    c.Vin,
	}
}

func (c *Car) ToModelBase() model.CarBase {
	return model.CarBase{
		Brand:          c.Brand,
		Model:          c.Model,
		ProductionDate: c.ProductionDate,
		Vin:            c.Vin,
	}
}

func newTechnicalSpecificationFromModel(t *model.TechnicalSpecification) TechnicalSpecification {
	return TechnicalSpecification{
		Color: t.Color,
		Consumption: TechnicalSpecificationConsumption{
			City:     t.Consumption.City,
			Combined: t.Consumption.Combined,
			Overland: t.Consumption.Overland,
		},
		Emissions: TechnicalSpecificationEmissions{
			City:     t.Emissions.City,
			Combined: t.Emissions.Combined,
			Overland: t.Emissions.Overland,
		},
		Engine: TechnicalSpecificationEngine{
			Power: t.Engine.Power,
			Type:  t.Engine.Type,
		},
		Fuel:          TechnicalSpecificationFuel(t.Fuel),
		FuelCapacity:  t.FuelCapacity,
		NumberOfDoors: t.NumberOfDoors,
		NumberOfSeats: t.NumberOfSeats,
		Tire: TechnicalSpecificationTire{
			Manufacturer: t.Tire.Manufacturer,
			Type:         t.Tire.Type,
		},
		Transmission: TechnicalSpecificationTransmission(t.Transmission),
		TrunkVolume:  t.TrunkVolume,
		Weight:       t.Weight,
	}
}

func newDynamicDataFromModel(d *model.DynamicData) DynamicData {
	return DynamicData{
		DoorsLockState:      DynamicDataLockState(d.DoorsLockState),
		EngineState:         DynamicDataEngineState(d.EngineState),
		FuelLevelPercentage: d.FuelLevelPercentage,
		Position: DynamicDataPosition{
			Latitude:  d.Position.Latitude,
			Longitude: d.Position.Latitude,
		},
		TrunkLockState: DynamicDataLockState(d.TrunkLockState),
	}
}

func NewCarFromModel(c *model.Car) Car {
	return Car{
		Brand:                  c.Brand,
		DynamicData:            newDynamicDataFromModel(&c.DynamicData),
		Model:                  c.Model,
		ProductionDate:         c.ProductionDate,
		TechnicalSpecification: newTechnicalSpecificationFromModel(&c.TechnicalSpecification),
		Vin:                    c.Vin,
	}
}
