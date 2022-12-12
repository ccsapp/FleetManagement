package model

import "regexp"

var (
	vinRegex     = regexp.MustCompile("^[A-HJ-NPR-Z0-9]{13}[0-9]{4}$")
	fleetIdRegex = regexp.MustCompile("^[a-zA-Z0-9]{8}$")
)

// IsVinValid validates a VIN against the defined format
func IsVinValid(vin Vin) bool {
	return vinRegex.MatchString(vin)
}

// IsFleetIDValid validates a fleet ID against the defined format
func IsFleetIDValid(fleetID FleetID) bool {
	return fleetIdRegex.MatchString(fleetID)
}
