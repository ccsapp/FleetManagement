package errors

import "errors"

var (
	ErrFleetNotFound      = errors.New("no such fleet")
	ErrCarNotFound        = errors.New("no such fleet")
	ErrCarNotInFleet      = errors.New("car not in fleet")
	ErrCarAlreadyInFleet  = errors.New("car already in fleet")
	ErrFleetAlreadyExists = errors.New("fleet already exists")
	ErrDomainAssertion    = errors.New("unexpected response from domain service")
	ErrInvalidVin         = errors.New("invalid vin")
	ErrInvalidFleetId     = errors.New("invalid fleet id")
)
