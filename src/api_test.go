package main

import (
	"PFleetManagement/infrastructure/database"
	"PFleetManagement/testdata"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type ApiTestSuite struct {
	suite.Suite
	fleetDB database.FleetDB
	app     *echo.Echo
	config  *Config
}

func (suite *ApiTestSuite) SetupSuite() {
	// load the environment variables for the database layer
	dbConfig, err := database.LoadConfigFromFile("testdata/testdb.env")
	if err != nil {
		suite.T().Fatal(err.Error())
	}

	suite.config = &Config{allowOrigins: []string{"*"}, domainServer: "https://carservice.kit.edu", domainTimeout: 1}

	// generate a collection name so that concurrent executions do not interfere
	dbConfig.CollectionPrefix = fmt.Sprintf("test-%d-", time.Now().Unix())

	suite.fleetDB, err = database.OpenDatabase(dbConfig)

	suite.app, err = newApp(suite.config, suite.fleetDB)
	if err != nil {
		suite.T().Fatal(err.Error())
	}

	// we need to initially clear the database since by default, an empty fleet is inserted into the database
	suite.TearDownTest()
}

func (suite *ApiTestSuite) TearDownSuite() {
	// close the database connection when the program exits
	if err := suite.fleetDB.CleanUpDatabase(); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *ApiTestSuite) TearDownTest() {
	// clear the collection after each test
	if err := suite.fleetDB.DropCollection(context.Background()); err != nil {
		suite.T().Fatal(err)
	}
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func newApiTest(handler http.Handler, name string) *apitest.APITest {
	return apitest.New(name).
		Debug().
		Handler(handler).
		Report(apitest.SequenceDiagram())
}

func newApiTestWithMocks(handler http.Handler, name string, mocks []*apitest.Mock) *apitest.APITest {
	return apitest.New(name).
		Mocks(mocks...).
		Debug().
		Handler(handler).
		Report(apitest.SequenceDiagram())
}

func newCarMock(suite *ApiTestSuite) []*apitest.Mock {
	return []*apitest.Mock{
		apitest.NewMock().
			Get(suite.config.domainServer + "/cars/" + testdata.VinCar).
			RespondWith().Status(http.StatusOK).Body(testdata.ExampleCar).End(),
		apitest.NewMock().
			Get(suite.config.domainServer + "/cars/" + testdata.VinCar2).
			RespondWith().Status(http.StatusOK).Body(testdata.ExampleCar2).End(),
		apitest.NewMock().
			Get(suite.config.domainServer + "/cars/" + testdata.UnknownVin).
			RespondWith().Status(http.StatusNotFound).End(),
	}
}

func (suite *ApiTestSuite) TestGetCars_invalidFleetId() {
	newApiTest(suite.app, "Get cars of a fleet with invalid id").
		Get("/fleets/abc/cars").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestGetCars_unknownFleet() {
	newApiTest(suite.app, "Get cars of a fleet unknown to the system").
		Get("/fleets/xk48jpgz/cars").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestGetCars_successEmpty() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	newApiTest(suite.app, "Get cars with no cars in the fleet").
		Get("/fleets/" + testdata.FleetId + "/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body("[]").
		End()
}

func (suite *ApiTestSuite) TestGetCars_success() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	newApiTestWithMocks(suite.app, "Add car success", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	newApiTestWithMocks(suite.app, "Add car success", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar2).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCar2Response).
		End()
	newApiTestWithMocks(suite.app, "Get cars success", newCarMock(suite)).
		Get("/fleets/" + testdata.FleetId + "/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleFleetOverview).
		End()
}

func (suite *ApiTestSuite) TestGetCar_invalidFleetId() {
	newApiTest(suite.app, "Get car of a fleet with invalid id").
		Get("/fleets/abc/cars/G1YZ23J9P58034278").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestGetCar_invalidVin() {
	newApiTest(suite.app, "Get car with invalid vin").
		Get("/fleets/xk48jpgz/cars/G1YZ23J9P5803427").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestGetCar_unknownFleet() {
	newApiTest(suite.app, "Get car of a fleet unknown to the system").
		Get("/fleets/xk48jpgz/cars/G1YZ23J9P58034278").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestGetCar_unknownCar() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	newApiTest(suite.app, "Get car unknown to the system").
		Get("/fleets/" + testdata.FleetId + "/cars/" + testdata.UnknownVin).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestGetCar_CarInOtherFleet() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId2); err != nil {
		suite.T().Fatal(err)
	}
	newApiTestWithMocks(suite.app, "Add car success", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId2 + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	newApiTest(suite.app, "Get car in different fleet").
		Get("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestGetCar_success() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	newApiTestWithMocks(suite.app, "Add car success", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	newApiTestWithMocks(suite.app, "Get car success", newCarMock(suite)).
		Get("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCar).
		End()
}

func (suite *ApiTestSuite) TestAddCar_invalidFleetId() {
	newApiTest(suite.app, "Add car to a fleet with invalid id").
		Put("/fleets/abc/cars/G1YZ23J9P58034278").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestAddCar_invalidVin() {
	newApiTest(suite.app, "Add car with invalid vin").
		Put("/fleets/" + testdata.FleetId + "/cars/G1YZ23J9P5803427").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestAddCar_unknownFleet() {
	newApiTestWithMocks(suite.app, "Add car to a fleet unknown to the system", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestAddCar_unknownCar() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	newApiTestWithMocks(suite.app, "Add car unknown to the system", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.UnknownVin).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestAddCar_duplicate() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	newApiTestWithMocks(suite.app, "Add car success", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	newApiTestWithMocks(suite.app, "Add car already in fleet", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}

func (suite *ApiTestSuite) TestAddCar_success() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	newApiTestWithMocks(suite.app, "Add car success", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
}

func (suite *ApiTestSuite) TestAddCar_successMultipleFleets() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId2); err != nil {
		suite.T().Fatal(err)
	}
	newApiTestWithMocks(suite.app, "Add car success", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	newApiTestWithMocks(suite.app, "Add car success", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId2 + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_invalidFleetId() {
	newApiTest(suite.app, "Remove car from a fleet with invalid id").
		Delete("/fleets/abc/cars/G1YZ23J9P58034278").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_invalidVin() {
	newApiTest(suite.app, "Remove car with invalid vin").
		Delete("/fleets/" + testdata.FleetId + "/cars/G1YZ23J9P5803427").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_unknownFleet() {
	newApiTest(suite.app, "Remove car from a fleet unknown to the system").
		Delete("/fleets/" + testdata.FleetId + "/cars/G1YZ23J9P58034278").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_unknownCar() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	newApiTest(suite.app, "Remove car unknown to the system").
		Delete("/fleets/" + testdata.FleetId + "/cars/" + testdata.UnknownVin).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_CarInOtherFleet() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId2); err != nil {
		suite.T().Fatal(err)
	}
	newApiTestWithMocks(suite.app, "Add car success", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId2 + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	newApiTest(suite.app, "Remove car in other fleet").
		Delete("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_success() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	newApiTestWithMocks(suite.app, "Add car success", newCarMock(suite)).
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	newApiTest(suite.app, "Remove car success").
		Delete("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}
