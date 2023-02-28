package environment

import "os"

func setCarServerUrl(carServerUrl string) {
	_ = os.Setenv(envCarServerUrl, carServerUrl)
}

func setRentalServerUrl(rentalServerUrl string) {
	_ = os.Setenv(envRentalServerUrl, rentalServerUrl)
}
