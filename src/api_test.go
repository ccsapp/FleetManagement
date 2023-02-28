package main

import (
	"PFleetManagement/environment"
	"PFleetManagement/infrastructure/database"
	"PFleetManagement/testdata"
	"PFleetManagement/testhelpers"
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
	fleetDB            database.FleetDB
	app                *echo.Echo
	recordingFormatter *testhelpers.RecordingFormatter
}

func (suite *ApiTestSuite) SetupSuite() {
	environment.SetupTestingEnvironment(
		"https://carservice.kit.edu",
		"https://rentalmanagement.kit.edu",
	)

	// generate a collection name so that concurrent executions do not interfere
	collectionPrefix := fmt.Sprintf("test-%d-", time.Now().Unix())
	environment.GetEnvironment().SetAppCollectionPrefix(collectionPrefix)

	var err error
	suite.fleetDB, err = database.OpenDatabase(environment.GetEnvironment())
	if err != nil {
		suite.handleDbConnectionError(err)
	}

	suite.app, err = newApp(suite.fleetDB)
	if err != nil {
		suite.T().Fatal(err.Error())
	}

	// we need to initially clear the database since by default, an empty fleet is inserted into the database
	suite.clearCollection()
}

func (suite *ApiTestSuite) handleDbConnectionError(err error) {
	// if local setup mode is disabled, we fail without any additional checks
	if !environment.GetEnvironment().IsLocalSetupMode() {
		suite.T().Fatal(err.Error())
	}

	running, dockerErr := testhelpers.IsMongoDbContainerRunning()
	if dockerErr != nil {
		suite.T().Fatal(dockerErr.Error())
	}
	if !running {
		suite.T().Fatal("MongoDB container is not running. " +
			"Please start it with 'docker compose up -d' and try again.")
	}
	fmt.Println("MongoDB container seems to be running, but the connection could not be established. " +
		"Please check the logs for more information.")
	suite.T().Fatal(err.Error())
}

func (suite *ApiTestSuite) clearCollection() {
	if err := suite.fleetDB.DropCollection(context.Background()); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *ApiTestSuite) SetupTest() {
	suite.recordingFormatter = testhelpers.NewRecordingFormatter()
}

func (suite *ApiTestSuite) TearDownTest() {
	// generate the sequence diagram for the test
	suite.recordingFormatter.SetOutFileName(suite.T().Name())
	suite.recordingFormatter.SetTitle(suite.T().Name())

	diagramFormatter := apitest.SequenceDiagram()
	diagramFormatter.Format(suite.recordingFormatter.GetRecorder())

	// clear the collection after each test
	suite.clearCollection()
}

func (suite *ApiTestSuite) TearDownSuite() {
	// close the database connection when the program exits
	if err := suite.fleetDB.CleanUpDatabase(); err != nil {
		suite.T().Fatal(err)
	}
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (suite *ApiTestSuite) newApiTest() *apitest.APITest {
	return apitest.New().
		Debug().
		Handler(suite.app).
		Report(suite.recordingFormatter)
}

func (suite *ApiTestSuite) newApiTestWithMocks(mocks []*apitest.Mock) *apitest.APITest {
	return apitest.New().
		Mocks(mocks...).
		Debug().
		Handler(suite.app).
		Report(suite.recordingFormatter)
}

func (suite *ApiTestSuite) newApiTestWithCarMock() *apitest.APITest {
	return suite.newApiTestWithMocks(suite.newCarMock())
}

func (suite *ApiTestSuite) newApiTestWithCarAndRentalMocks() *apitest.APITest {
	return suite.newApiTestWithMocks(append(suite.newCarMock(), suite.newRentalMock()...))
}

func (suite *ApiTestSuite) newCarMock() []*apitest.Mock {
	return []*apitest.Mock{
		apitest.NewMock().
			Get(environment.GetEnvironment().GetCarServerUrl() + "/cars/" + testdata.VinCar).
			RespondWith().Status(http.StatusOK).Body(testdata.ExampleCar).End(),
		apitest.NewMock().
			Get(environment.GetEnvironment().GetCarServerUrl() + "/cars/" + testdata.VinCar2).
			RespondWith().Status(http.StatusOK).Body(testdata.ExampleCar2).End(),
		apitest.NewMock().
			Get(environment.GetEnvironment().GetCarServerUrl() + "/cars/" + testdata.UnknownVin).
			RespondWith().Status(http.StatusNotFound).End(),
	}
}

func (suite *ApiTestSuite) newRentalMock() []*apitest.Mock {
	return []*apitest.Mock{
		apitest.NewMock().
			Get(environment.GetEnvironment().GetRentalServerUrl() + "/cars/" + testdata.VinCar + "/rentalStatus").
			RespondWith().Status(http.StatusNoContent).End(),
		apitest.NewMock().
			Get(environment.GetEnvironment().GetRentalServerUrl() + "/cars/" + testdata.VinCar2 + "/rentalStatus").
			RespondWith().Status(http.StatusOK).Body(testdata.ExampleRental).End(),
	}
}

func (suite *ApiTestSuite) TestGetCars_invalidFleetId() {
	suite.newApiTest().
		Get("/fleets/abc/cars").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestGetCars_unknownFleet() {
	suite.newApiTest().
		Get("/fleets/xk48jpgz/cars").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestGetCars_successEmpty() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	suite.newApiTest().
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
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar2).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCar2Response).
		End()
	suite.newApiTestWithCarMock().
		Get("/fleets/" + testdata.FleetId + "/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleFleetOverview).
		End()
}

func (suite *ApiTestSuite) TestGetCar_invalidFleetId() {
	suite.newApiTest().
		Get("/fleets/abc/cars/G1YZ23J9P58034278").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestGetCar_invalidVin() {
	suite.newApiTest().
		Get("/fleets/xk48jpgz/cars/G1YZ23J9P5803427").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestGetCar_unknownFleet() {
	suite.newApiTest().
		Get("/fleets/xk48jpgz/cars/G1YZ23J9P58034278").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestGetCar_unknownCar() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	suite.newApiTest().
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
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId2 + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	suite.newApiTest().
		Get("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestGetCar_success_noRental() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	suite.newApiTestWithCarAndRentalMocks().
		Get("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCar).
		End()
}

func (suite *ApiTestSuite) TestGetCar_success_rental() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar2).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCar2Response).
		End()
	suite.newApiTestWithCarAndRentalMocks().
		Get("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar2).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCar2WithRental).
		End()
}

func (suite *ApiTestSuite) TestAddCar_invalidFleetId() {
	suite.newApiTest().
		Put("/fleets/abc/cars/G1YZ23J9P58034278").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestAddCar_invalidVin() {
	suite.newApiTest().
		Put("/fleets/" + testdata.FleetId + "/cars/G1YZ23J9P5803427").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestAddCar_unknownFleet() {
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestAddCar_unknownCar() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.UnknownVin).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestAddCar_duplicate() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}

func (suite *ApiTestSuite) TestAddCar_success() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	suite.newApiTestWithCarMock().
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
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId2 + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_invalidFleetId() {
	suite.newApiTest().
		Delete("/fleets/abc/cars/G1YZ23J9P58034278").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_invalidVin() {
	suite.newApiTest().
		Delete("/fleets/" + testdata.FleetId + "/cars/G1YZ23J9P5803427").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_unknownFleet() {
	suite.newApiTest().
		Delete("/fleets/" + testdata.FleetId + "/cars/G1YZ23J9P58034278").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_unknownCar() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	suite.newApiTest().
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
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId2 + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	suite.newApiTest().
		Delete("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_success() {
	if err := suite.fleetDB.AddFleet(context.Background(), testdata.FleetId); err != nil {
		suite.T().Fatal(err)
	}
	suite.newApiTestWithCarMock().
		Put("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarResponse).
		End()
	suite.newApiTest().
		Delete("/fleets/" + testdata.FleetId + "/cars/" + testdata.VinCar).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}
