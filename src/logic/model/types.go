// Package model provides the data model of this service
package model

import (
	openapiTypes "github.com/deepmap/oapi-codegen/pkg/types"
	"time"
)

// Defines values for DynamicDataEngineState.
const (
	OFF DynamicDataEngineState = "OFF"
	ON  DynamicDataEngineState = "ON"
)

// Defines values for LockState.
const (
	LOCKED   LockState = "LOCKED"
	UNLOCKED LockState = "UNLOCKED"
)

// Defines values for TechnicalSpecificationFuel.
const (
	DIESEL       TechnicalSpecificationFuel = "DIESEL"
	ELECTRIC     TechnicalSpecificationFuel = "ELECTRIC"
	HYBRIDDIESEL TechnicalSpecificationFuel = "HYBRID_DIESEL"
	HYBRIDPETROL TechnicalSpecificationFuel = "HYBRID_PETROL"
	PETROL       TechnicalSpecificationFuel = "PETROL"
)

// Defines values for TechnicalSpecificationTransmission.
const (
	AUTOMATIC TechnicalSpecificationTransmission = "AUTOMATIC"
	MANUAL    TechnicalSpecificationTransmission = "MANUAL"
)

// Car defines model for car.
type Car struct {
	// Brand Data that specifies the brand name of the Vehicle manufacturer
	Brand string `json:"brand"`

	// DynamicData Data that changes during a car's operation
	DynamicData DynamicData `json:"dynamicData"`

	// Model Data that specifies the particular type of Vehicle
	Model string `json:"model"`

	// ProductionDate Data that specifies the official date the vehicle was declared to have exited production by the manufacturer.
	ProductionDate         openapiTypes.Date      `json:"productionDate"`
	TechnicalSpecification TechnicalSpecification `json:"technicalSpecification"`

	// Vin A Vehicle Identification Number (VIN) which uniquely identifies a Vehicle
	Vin Vin `json:"vin"`

	// Rental Data that specifies the active or next upcoming rental of a car
	Rental *Rental `json:"rental,omitempty"`
}

// CarBase Overview of a car
type CarBase struct {
	// Brand Data that specifies the brand name of the Vehicle manufacturer
	Brand string `json:"brand"`

	// Model Data that specifies the particular type of Vehicle
	Model string `json:"model"`

	// ProductionDate Data that specifies the official date the vehicle was declared to have exited production by the manufacturer.
	ProductionDate openapiTypes.Date `json:"productionDate"`

	// Vin A Vehicle Identification Number (VIN) which uniquely identifies a Vehicle
	Vin Vin `json:"vin"`
}

// DynamicData Data that changes during a car's operation
type DynamicData struct {
	// DoorsLockState Data that specifies whether an object is locked or unlocked
	DoorsLockState LockState              `json:"doorsLockState"`
	EngineState    DynamicDataEngineState `json:"engineState"`

	// FuelLevelPercentage Data that specifies the relation of remaining fuelCapacity to the maximum fuelCapacity in percentage
	FuelLevelPercentage int `json:"fuelLevelPercentage"`

	// Position Data that specifies the GeoCoordinate of a car
	Position DynamicDataPosition `json:"position"`

	// TrunkLockState Data that specifies whether an object is locked or unlocked
	TrunkLockState LockState `json:"trunkLockState"`
}

// DynamicDataEngineState defines model for DynamicData.EngineState.
type DynamicDataEngineState string

// FleetID Unique identification of a car fleet
type FleetID = string

// LockState Data that specifies whether an object is locked or unlocked
type LockState string

// TechnicalSpecification defines model for technicalSpecification.
type TechnicalSpecification struct {
	// Color Data on the description of the paint job of a car
	Color string `json:"color"`

	// Consumption Data that specifies the amount of fuel consumed during car operation in units per 100 kilometers
	Consumption TechnicalSpecificationConsumption `json:"consumption"`

	// Emissions Data that specifies the CO2 emitted by a car during operation in gram per kilometer
	Emissions TechnicalSpecificationEmissions `json:"emissions"`

	// Engine A physical unit that converts fuel into movement
	Engine TechnicalSpecificationEngine `json:"engine"`

	// Fuel Data that defines the source of energy that powers the car
	Fuel TechnicalSpecificationFuel `json:"fuel"`

	// FuelCapacity Data that specifies the amount of fuel that can be carried with the car
	FuelCapacity string `json:"fuelCapacity"`

	// NumberOfDoors Data that defines the number of doors that are built into a car
	NumberOfDoors int `json:"numberOfDoors"`

	// NumberOfSeats Data that defines the number of seats that are built into a car
	NumberOfSeats int `json:"numberOfSeats"`

	// Tire A physical unit that serves as the point of contact between a car and the ground
	Tire TechnicalSpecificationTire `json:"tire"`

	// Transmission A physical unit responsible for managing the conversion rate of the engine (can be automated or manually operated)
	Transmission TechnicalSpecificationTransmission `json:"transmission"`

	// TrunkVolume Data on the physical volume of the trunk in liters
	TrunkVolume int `json:"trunkVolume"`

	// Weight Data that specifies the total weight of a car when empty in kilograms (kg)
	Weight int `json:"weight"`
}

type DynamicDataPosition struct {
	// Latitude Data that specifies the distance from the equator
	Latitude float32 `json:"latitude"`

	// Longitude Data that specifies the distance east or west from a line (meridian) passing through Greenwich
	Longitude float32 `json:"longitude"`
}

type TechnicalSpecificationConsumption struct {
	// City Data that specifies the amount of fuel that is consumed when driving within the city in: kW/100km or l/100km
	City float32 `json:"city"`

	// Combined Data that specifies the combined amount of fuel that is consumed in: kW / 100 km or l / 100 km
	Combined float32 `json:"combined"`

	// Overland Data that specifies the amount of fuel that is consumed when driving outside a city in: kW/100km or l/100km
	Overland float32 `json:"overland"`
}

type TechnicalSpecificationEmissions struct {
	// City Data that specifies the amount of emissions when driving within the city in: g CO2 / km
	City float32 `json:"city"`

	// Combined Data that specifies the combined amount of emissions in: g CO2 / km. The combination is done by the manufacturer according to an industry-specific standard
	Combined float32 `json:"combined"`

	// Overland Data that specifies the amount of emissions when driving outside a city in: g CO2 / km
	Overland float32 `json:"overland"`
}

type TechnicalSpecificationEngine struct {
	// Power Data on the power the engine can provide in kW
	Power int `json:"power"`

	// Type Data that contains the manufacturer-given type description of the engine
	Type string `json:"type"`
}

type TechnicalSpecificationTire struct {
	// Manufacturer Data denoting the company responsible for the creation of a physical unit
	Manufacturer string `json:"manufacturer"`

	// Type Data that contains the manufacturer-given type description of the tire
	Type string `json:"type"`
}

// TechnicalSpecificationFuel Data that defines the source of energy that powers the car
type TechnicalSpecificationFuel string

// TechnicalSpecificationTransmission A physical unit responsible for managing the conversion rate of the engine (can be automated or manually operated)
type TechnicalSpecificationTransmission string

// Vin A Vehicle Identification Number (VIN) which uniquely identifies a Vehicle
type Vin = string

// FleetIDParam Unique identification of a car fleet
type FleetIDParam = FleetID

// VinParam A Vehicle Identification Number (VIN) which uniquely identifies a Vehicle
type VinParam = Vin

// Rental defines a model for rentals.
type Rental struct {
	// Active Describes whether this rental is active
	Active bool `json:"active"`

	// Id Unique identification of a rental
	Id RentalId `json:"id"`

	// Customer The renting customer
	Customer Customer `json:"customer"`

	// RentalPeriod A period of time
	RentalPeriod TimePeriod `json:"rentalPeriod"`
}

// TimePeriod A period of time
type TimePeriod struct {
	// StartDate start of the time period
	StartDate time.Time `json:"startDate"`

	// EndDate end of the time period
	EndDate time.Time `json:"endDate"`
}

// RentalId Unique identification of a rental
type RentalId = string

// Customer A customer
type Customer struct {
	// CustomerId Unique identification of a customer
	CustomerId CustomerId `json:"customerId"`
}

// CustomerId Unique identification of a customer
type CustomerId = string
