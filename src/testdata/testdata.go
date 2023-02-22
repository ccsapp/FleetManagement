package testdata

import _ "embed"

const UnknownVin string = "G1YZ23J9P58034278"
const FleetId string = "xk48jpgz"
const FleetId2 string = "xk49jpgz"

//go:embed exampleCarResponse.json
var ExampleCarResponse string

//go:embed exampleCar.json
var ExampleCar string

const VinCar string = "WVWAA71K08W201030"

//go:embed exampleCar2Response.json
var ExampleCar2Response string

//go:embed exampleCar2.json
var ExampleCar2 string

//go:embed exampleCar2WithRental.json
var ExampleCar2WithRental string

const VinCar2 string = "WVWAA71K08W201031"

//go:embed exampleFleetOverview.json
var ExampleFleetOverview string

//go:embed exampleRental.json
var ExampleRental string
