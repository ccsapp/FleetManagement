package dcar

import (
	"PFleetManagement/logic/model"
	carTypes "github.com/ccsapp/cargotypes"
	openapiTypes "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var car1 = carTypes.Car{
	Brand: "Tesla",
	DynamicData: carTypes.DynamicData{
		DoorsLockState:      carTypes.LOCKED,
		EngineState:         carTypes.ON,
		FuelLevelPercentage: 87,
		Position: carTypes.DynamicDataPosition{
			Latitude:  49.0069,
			Longitude: 8.4037,
		},
		TrunkLockState: carTypes.LOCKED,
	},
	Model: "Model X",
	ProductionDate: openapiTypes.Date{
		Time: time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
	},
	TechnicalSpecification: carTypes.TechnicalSpecification{
		Color: "red",
		Consumption: carTypes.TechnicalSpecificationConsumption{
			City:     2.1,
			Combined: 5.2,
			Overland: 7.3,
		},
		Emissions: carTypes.TechnicalSpecificationEmissions{
			City:     0.2,
			Combined: 1.2,
			Overland: 2.4,
		},
		Engine: carTypes.TechnicalSpecificationEngine{
			Power: 10,
			Type:  "110 CDM",
		},
		Fuel:          "HYBRID_DIESEL",
		FuelCapacity:  "54.0L;85.2kWh",
		NumberOfDoors: 3,
		NumberOfSeats: 7,
		Tire: carTypes.TechnicalSpecificationTire{
			Manufacturer: "GOODYEAR",
			Type:         "195/65R15",
		},
		Transmission: "manual",
		TrunkVolume:  90,
		Weight:       1002,
	},
	Vin: "3B7HF13Y81G193584",
}

var modelCar1 = model.Car{
	Brand: "Tesla",
	DynamicData: model.DynamicData{
		DoorsLockState:      model.LOCKED,
		EngineState:         model.ON,
		FuelLevelPercentage: 87,
		Position: model.DynamicDataPosition{
			Latitude:  49.0069,
			Longitude: 8.4037,
		},
		TrunkLockState: model.LOCKED,
	},
	Model: "Model X",
	ProductionDate: openapiTypes.Date{
		Time: time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
	},
	TechnicalSpecification: model.TechnicalSpecification{
		Color: "red",
		Consumption: model.TechnicalSpecificationConsumption{
			City:     2.1,
			Combined: 5.2,
			Overland: 7.3,
		},
		Emissions: model.TechnicalSpecificationEmissions{
			City:     0.2,
			Combined: 1.2,
			Overland: 2.4,
		},
		Engine: model.TechnicalSpecificationEngine{
			Power: 10,
			Type:  "110 CDM",
		},
		Fuel:          "HYBRID_DIESEL",
		FuelCapacity:  "54.0L;85.2kWh",
		NumberOfDoors: 3,
		NumberOfSeats: 7,
		Tire: model.TechnicalSpecificationTire{
			Manufacturer: "GOODYEAR",
			Type:         "195/65R15",
		},
		Transmission: "manual",
		TrunkVolume:  90,
		Weight:       1002,
	},
	Vin: "3B7HF13Y81G193584",
}

var modelBase1 = model.CarBase{
	Brand: "Tesla",
	Model: "Model X",
	ProductionDate: openapiTypes.Date{
		Time: time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
	},
	Vin: "3B7HF13Y81G193584",
}

func TestCar_ToModel(t *testing.T) {
	modelResult := ToModelFromCar(&car1)

	assert.Equal(t, modelCar1, modelResult)
}

func TestCar_ToModelBase(t *testing.T) {
	modelResult := ToModelBaseFromCar(&car1)

	assert.Equal(t, modelBase1, modelResult)
}

func TestNewCarFromModel(t *testing.T) {
	result := NewCarFromModel(&modelCar1)

	assert.Equal(t, car1, result)
}
