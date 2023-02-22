package operations

import (
	"PFleetManagement/infrastructure/database"
	"PFleetManagement/infrastructure/dcar"
	rentalManagement "PFleetManagement/infrastructure/rentalmanagement"
	"PFleetManagement/logic/fleetErrors"
	"PFleetManagement/logic/model"
	"context"
	"fmt"
	"net/http"
)

type operations struct {
	database               database.FleetDB
	carClient              dcar.ClientWithResponsesInterface
	rentalManagementClient rentalManagement.ClientWithResponsesInterface
}

// NewOperations creates an implementation of IOperations from its dependencies.
//
// The database.FleetDB is queried for/updated with the fleet-car assignment.
//
// The dcar.ClientWithResponsesInterface is queried for resolving VINs to full car data.
func NewOperations(fleetDB database.FleetDB, carClient dcar.ClientWithResponsesInterface,
	rentalManagementClient rentalManagement.ClientWithResponsesInterface) IOperations {

	return operations{
		database:               fleetDB,
		carClient:              carClient,
		rentalManagementClient: rentalManagementClient,
	}
}

func (o operations) GetCarsInFleet(ctx context.Context, fleetID model.FleetID) ([]model.CarBase, error) {
	// --- database interaction ---
	vins, err := o.database.GetCarsForFleet(ctx, fleetID)
	if err != nil {
		return nil, err
	}

	// --- Car service interaction ---

	// create an array to hold the car (base) objects for all those VINs returned by the database
	cars := make([]model.CarBase, len(vins))

	// iterate over the VINs and query for the cars respectively
	for index, vin := range vins {
		carResponse, err := o.carClient.GetCarWithResponse(ctx, vin)
		if err != nil {
			return nil, err
		}

		if carResponse.JSON200 == nil {
			// if the retrieval fails for any car, the whole operation fails
			statusCode := carResponse.StatusCode()
			if statusCode == http.StatusNotFound {
				// (car was deleted since it was added to the fleet -> fleet database references unknown data)
				// this error by the domain is known but results from an inconsistency
				return nil, fmt.Errorf("%w: car %s from fleet %s not in domain", fleetErrors.ErrDomainAssertion, vin, fleetID)
			}
			return nil, fmt.Errorf("%w: unknown error (domain code %d)", fleetErrors.ErrDomainAssertion, statusCode)
		}

		// if the car data could be retrieved -> add the data to the array
		// remark: index just increases in every loop, so it's writing successive values to the array
		cars[index] = dcar.ToModelBaseFromCar(carResponse.JSON200)
	}

	return cars, nil
}

func (o operations) RemoveCar(ctx context.Context, fleetID model.FleetID, vin model.Vin) error {
	// --- database interaction ---
	return o.database.RemoveCarFromFleet(ctx, fleetID, vin)
}

func (o operations) GetCar(ctx context.Context, fleetID model.FleetID, vin model.Vin) (*model.Car, error) {
	// --- database interaction ---
	// while it would be possible to only query the Car service and not the database at all,
	// there would be no way of knowing whether this operation is actually valid -> check whether
	// car assigned to fleet at all, first
	carInFleet, err := o.database.IsCarInFleet(ctx, fleetID, vin)
	if err != nil {
		return nil, err
	}
	if !carInFleet {
		return nil, fleetErrors.ErrCarNotInFleet
	}

	// --- Car service interaction ---
	// get the car data itself
	response, err := o.carClient.GetCarWithResponse(ctx, vin)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("%w: status code %d", fleetErrors.ErrDomainAssertion, response.StatusCode())
	}

	// --- Rental management service interaction ---
	// get the rental data for the car
	rentalResponse, err := o.rentalManagementClient.GetNextRentalWithResponse(ctx, vin)
	if err != nil {
		return nil, err
	}
	if rentalResponse.JSON200 == nil && rentalResponse.StatusCode() != http.StatusNoContent {
		return nil, fmt.Errorf("%w: error code %d", fleetErrors.ErrRentalManagementAssertion, rentalResponse.StatusCode())
	}

	carData := dcar.ToModelFromCar(response.JSON200)
	carData.Rental = rentalResponse.JSON200
	return &carData, nil
}

func (o operations) AddCarToFleet(ctx context.Context, fleetID model.FleetID, vin model.Vin) (*model.CarBase, error) {
	// --- Car service interaction ---
	// check for the car first to prevent VINs which no cars are known for to be stored in the database
	carResponse, err := o.carClient.GetCarWithResponse(ctx, vin)
	if err != nil {
		return nil, err
	}

	if carResponse.JSON200 == nil {
		statusCode := carResponse.StatusCode()
		if statusCode == http.StatusNotFound {
			// it is a defined error of this operation that the car does not exist
			return nil, fleetErrors.ErrCarNotFound
		} else {
			return nil, fmt.Errorf("%w: unknown error (domain code %d)", fleetErrors.ErrDomainAssertion, statusCode)
		}
	}

	// --- database interaction ---
	err = o.database.AddCarToFleet(ctx, fleetID, vin)
	if err != nil {
		return nil, err
	}

	// if this line is executed, carResponse.JSON200 is not nil -> return the necessary information
	baseData := dcar.ToModelBaseFromCar(carResponse.JSON200)
	return &baseData, nil
}
