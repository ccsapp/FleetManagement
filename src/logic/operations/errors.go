package operations

import "errors"

var (
	ErrFleetNotFound     = errors.New("no such fleet")
	ErrCarNotFound       = errors.New("no such fleet")
	ErrCarNotInFleet     = errors.New("car not in fleet")
	ErrCarAlreadyInFleet = errors.New("car already in fleet")
)
