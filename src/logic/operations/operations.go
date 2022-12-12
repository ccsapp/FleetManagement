package operations

import (
	"PFleetManagement/infrastructure/database"
	"PFleetManagement/infrastructure/dcar"
	"PFleetManagement/logic/errors"
	"PFleetManagement/logic/model"
	"context"
	"fmt"
	"net/http"
)

type operations struct {
	database database.FleetDB
	dcar     dcar.ClientWithResponsesInterface
}

// NewOperations creates an implementation of IOperations from its dependencies.
//
// The database.FleetDB is queried for/updated with the fleet-car assignment.
//
// The dcar.ClientWithResponsesInterface is queried for resolving VINs to full car data.
func NewOperations(fleetDB database.FleetDB, dcarClient dcar.ClientWithResponsesInterface) IOperations {
	return operations{
		database: fleetDB,
		dcar:     dcarClient,
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
		carResponse, err := o.dcar.GetCarWithResponse(ctx, vin)
		if err != nil {
			return nil, err
		}

		if carResponse.JSON200 != nil {
			// if the car data could be retrieved -> add the data to the array
			// remark: index just increases in every loop, so it's writing successive values to the array
			cars[index] = carResponse.JSON200.ToModelBase()
		} else {
			// if the retrieval fails for any car, the whole operation fails

			statusCode := carResponse.StatusCode()
			if statusCode == http.StatusNotFound {
				// this error by the domain is known but results from an inconsistency
				// (car was deleted since it was added to the fleet -> fleet database references unknown data)
				return nil, fmt.Errorf("%w: car %s from fleet %s not in domain", errors.ErrDomainAssertion, vin, fleetID)
			} else {
				return nil, fmt.Errorf("%w: unknown error (domain code %d)", errors.ErrDomainAssertion, statusCode)
			}
		}
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
		return nil, errors.ErrCarNotInFleet
	}

	// --- Car service interaction ---
	// get the car data itself
	response, err := o.dcar.GetCarWithResponse(ctx, vin)
	if err != nil {
		return nil, err
	}
	if response.JSON200 != nil {
		carData := response.JSON200.ToModel()
		return &carData, nil
	} else {
		return nil, fmt.Errorf("%w: error code %d", errors.ErrDomainAssertion, response.StatusCode())
	}
}

func (o operations) AddCarToFleet(ctx context.Context, fleetID model.FleetID, vin model.Vin) (*model.CarBase, error) {
	// --- Car service interaction ---
	// check for the car first to prevent VINs which no cars are known for to be stored in the database
	carResponse, err := o.dcar.GetCarWithResponse(ctx, vin)
	if err != nil {
		return nil, err
	}

	if carResponse.JSON200 == nil {
		statusCode := carResponse.StatusCode()
		if statusCode == http.StatusNotFound {
			// it is a defined error of this operation that the car does not exist
			return nil, errors.ErrCarNotFound
		} else {
			return nil, fmt.Errorf("%w: unknown error (domain code %d)", errors.ErrDomainAssertion, statusCode)
		}
	}

	// --- database interaction ---
	err = o.database.AddCarToFleet(ctx, fleetID, vin)
	if err != nil {
		return nil, err
	}

	// if this line is executed, carResponse.JSON200 is not nil -> return the necessary information
	baseData := carResponse.JSON200.ToModelBase()
	return &baseData, nil
}
