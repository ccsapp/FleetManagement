// Package errors defines the semantic errors which can occur while performing a task process
package errors

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

	// ErrDomainAssertion can occur for different reasons as a result of inconsistency with the Car service
	ErrDomainAssertion = errors.New("unexpected response from domain service")

	// ErrInvalidVin shows that the format of a VIN is invalid
	ErrInvalidVin = errors.New("invalid vin")

	// ErrInvalidFleetId shows that the format of a fleet ID is invalid
	ErrInvalidFleetId = errors.New("invalid fleet id")
)
