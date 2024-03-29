package fleetErrors

import "errors"

var (
	// ErrFleetNotFound shows the non-existence of a fleet with a given identifier
	ErrFleetNotFound = errors.New("no such fleet")

	// ErrCarNotFound shows the non-existence of a car with a given identifier
	ErrCarNotFound = errors.New("no such car")

	// ErrCarNotInFleet shows that, though the car and the fleet may exist, they are not connected
	ErrCarNotInFleet = errors.New("car not in fleet")

	// ErrCarAlreadyInFleet shows that a car with a given VIN is already assigned to a given fleet
	ErrCarAlreadyInFleet = errors.New("car already in fleet")

	// ErrFleetAlreadyExists shows that there already is a fleet with a given fleet ID
	ErrFleetAlreadyExists = errors.New("fleet already exists")

	// ErrDomainAssertion occurs when an unexpected response is received from the Car microservice
	ErrDomainAssertion = errors.New("unexpected response from domain service")

	// ErrRentalManagementAssertion occurs when an unexpected response is received from the RentalManagement microservice
	ErrRentalManagementAssertion = errors.New("unexpected response from rental management service")

	// ErrInvalidVin shows that the format of a VIN is invalid
	ErrInvalidVin = errors.New("invalid vin")

	// ErrInvalidFleetId shows that the format of a fleet ID is invalid
	ErrInvalidFleetId = errors.New("invalid fleet id")
)
