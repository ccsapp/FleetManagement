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

type Operations struct {
	database database.FleetDB
	dcar     dcar.ClientWithResponsesInterface
}

func NewOperations(fleetDB database.FleetDB, dcarClient dcar.ClientWithResponsesInterface) Operations {
	return Operations{
		database: fleetDB,
		dcar:     dcarClient,
	}
}

func (o Operations) GetCarsInFleet(ctx context.Context, fleetID model.FleetID) ([]model.CarBase, error) {
	vins, err := o.database.GetCarsForFleet(fleetID)
	if err != nil {
		return nil, err
	}

	cars := make([]model.CarBase, len(vins))
	for index, vin := range vins {
		carResponse, err := o.dcar.GetCarWithResponse(ctx, vin)
		if err != nil {
			return nil, err
		}

		if carResponse.JSON200 != nil {
			cars[index] = carResponse.JSON200.ToModelBase()
		} else {
			statusCode := carResponse.StatusCode()
			if statusCode == http.StatusNotFound {
				return nil, fmt.Errorf("%w: car %s from fleet %s not in domain", errors.ErrDomainAssertion, vin, fleetID)
			} else {
				return nil, fmt.Errorf("%w: unknown error (domain code %d)", errors.ErrDomainAssertion, statusCode)
			}
		}
	}

	return cars, nil
}

func (o Operations) RemoveCar(_ context.Context, fleetID model.FleetID, vin model.Vin) error {
	return o.database.RemoveCarFromFleet(fleetID, vin)
}

func (o Operations) GetCar(ctx context.Context, fleetID model.FleetID, vin model.Vin) (*model.Car, error) {
	// TODO maybe add database operation for find single
	vins, err := o.database.GetCarsForFleet(fleetID)
	if err != nil {
		return nil, err
	}

	for _, foundVin := range vins {
		if foundVin == vin {
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
	}

	return nil, errors.ErrCarNotInFleet
}

func (o Operations) AddCarToFleet(ctx context.Context, fleetID model.FleetID, vin model.Vin) (*model.CarBase, error) {
	carResponse, err := o.dcar.GetCarWithResponse(ctx, vin)
	if err != nil {
		return nil, err
	}

	if carResponse.JSON200 == nil {
		statusCode := carResponse.StatusCode()
		if statusCode == http.StatusNotFound {
			return nil, errors.ErrCarNotFound
		} else {
			return nil, fmt.Errorf("%w: unknown error (domain code %d)", errors.ErrDomainAssertion, statusCode)
		}
	}

	err = o.database.AddCarToFleet(fleetID, vin)
	if err != nil {
		return nil, err
	}

	baseData := carResponse.JSON200.ToModelBase()
	return &baseData, nil
}
