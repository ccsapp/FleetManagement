package operations

import (
	"PFleetManagement/infrastructure/dcar"
	"PFleetManagement/logic/fleetErrors"
	"PFleetManagement/logic/model"
	"PFleetManagement/mocks"
	"context"
	"errors"
	openapiTypes "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var carBase1 = model.CarBase{
	Brand: "Tesla",
	Model: "Model X",
	ProductionDate: openapiTypes.Date{
		Time: time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
	},
	Vin: "3B7HF13Y81G193584",
}

var car1 = dcar.Car{
	Brand:       "Tesla",
	DynamicData: dcar.DynamicData{},
	Model:       "Model X",
	ProductionDate: openapiTypes.Date{
		Time: time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
	},
	TechnicalSpecification: dcar.TechnicalSpecification{},
	Vin:                    "3B7HF13Y81G193584",
}

var modelCar1 = model.Car{
	Brand: "Tesla",
	Model: "Model X",
	ProductionDate: openapiTypes.Date{
		Time: time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
	},
	Vin: "3B7HF13Y81G193584",
}

var vins = []model.Vin{
	"3B7HF13Y81G193584",
}

var cars = []model.CarBase{
	carBase1,
}

func TestOperations_AddCarToFleet_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(&dcar.GetCarResponse{
		JSON200: &car1,
	}, nil)
	mockDatabase.EXPECT().AddCarToFleet(ctx, fleetID, vin).Return(nil)

	carBase, err := operations.AddCarToFleet(ctx, fleetID, vin)

	assert.Nil(t, err)
	assert.Equal(t, &carBase1, carBase)
}

func TestOperations_AddCarToFleet_domainError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)
	domainError := errors.New("domain error")

	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(nil, domainError)

	carBase, err := operations.AddCarToFleet(ctx, fleetID, vin)

	assert.ErrorIs(t, err, domainError)
	assert.Nil(t, carBase)
}

func TestOperations_AddCarToFleet_carNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193585"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(&dcar.GetCarResponse{
		HTTPResponse: &http.Response{
			StatusCode: http.StatusNotFound,
		},
	}, nil)

	carBase, err := operations.AddCarToFleet(ctx, fleetID, vin)

	assert.ErrorIs(t, err, fleetErrors.ErrCarNotFound)
	assert.Nil(t, carBase)
}

func TestOperations_AddCarToFleet_unexpectedDomainStatusCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193585"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(&dcar.GetCarResponse{
		HTTPResponse: &http.Response{
			StatusCode: http.StatusTeapot,
		},
	}, nil)

	carBase, err := operations.AddCarToFleet(ctx, fleetID, vin)

	assert.ErrorIs(t, err, fleetErrors.ErrDomainAssertion)
	assert.Nil(t, carBase)
}

func TestOperations_AddCarToFleet_databaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	databaseError := errors.New("database error")

	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(&dcar.GetCarResponse{
		JSON200: &car1,
	}, nil)
	mockDatabase.EXPECT().AddCarToFleet(ctx, fleetID, vin).Return(databaseError)

	carBase, err := operations.AddCarToFleet(ctx, fleetID, vin)

	assert.ErrorIs(t, err, databaseError)
	assert.Nil(t, carBase)
}

func TestOperations_GetCar_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	mockDatabase.EXPECT().IsCarInFleet(ctx, fleetID, vin).Return(true, nil)
	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(&dcar.GetCarResponse{
		JSON200: &car1,
	}, nil)

	car, err := operations.GetCar(ctx, fleetID, vin)

	assert.Nil(t, err)
	assert.Equal(t, &modelCar1, car)
}

func TestOperations_GetCar_databaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	databaseError := errors.New("database error")

	mockDatabase.EXPECT().IsCarInFleet(ctx, fleetID, vin).Return(false, databaseError)

	car, err := operations.GetCar(ctx, fleetID, vin)

	assert.ErrorIs(t, err, databaseError)
	assert.Nil(t, car)
}

func TestOperations_GetCar_notInFleet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193585"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	mockDatabase.EXPECT().IsCarInFleet(ctx, fleetID, vin).Return(false, nil)

	car, err := operations.GetCar(ctx, fleetID, vin)

	assert.ErrorIs(t, err, fleetErrors.ErrCarNotInFleet)
	assert.Nil(t, car)
}

func TestOperations_GetCar_domainError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	domainError := errors.New("domain error")

	mockDatabase.EXPECT().IsCarInFleet(ctx, fleetID, vin).Return(true, nil)
	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(nil, domainError)

	car, err := operations.GetCar(ctx, fleetID, vin)

	assert.ErrorIs(t, err, domainError)
	assert.Nil(t, car)
}

func TestOperations_GetCar_unexpectedDomainStatusCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	mockDatabase.EXPECT().IsCarInFleet(ctx, fleetID, vin).Return(true, nil)
	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(&dcar.GetCarResponse{
		HTTPResponse: &http.Response{
			StatusCode: http.StatusTeapot,
		},
	}, nil)

	car, err := operations.GetCar(ctx, fleetID, vin)

	assert.ErrorIs(t, err, fleetErrors.ErrDomainAssertion)
	assert.Nil(t, car)
}

func TestOperations_RemoveCar_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	mockDatabase.EXPECT().RemoveCarFromFleet(ctx, fleetID, vin).Return(nil)

	err := operations.RemoveCar(ctx, fleetID, vin)

	assert.Nil(t, err)
}

func TestOperations_RemoveCar_databaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	databaseError := errors.New("database error")

	operations := NewOperations(mockDatabase, mockDCar)

	mockDatabase.EXPECT().RemoveCarFromFleet(ctx, fleetID, vin).Return(databaseError)

	err := operations.RemoveCar(ctx, fleetID, vin)

	assert.ErrorIs(t, err, databaseError)
}

func TestOperations_GetCarsInFleet_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	mockDatabase.EXPECT().GetCarsForFleet(ctx, fleetID).Return(vins, nil)
	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(&dcar.GetCarResponse{
		JSON200: &car1,
	}, nil)

	retCars, err := operations.GetCarsInFleet(ctx, fleetID)

	assert.Nil(t, err)
	assert.Equal(t, cars, retCars)
}

func TestOperations_GetCarsInFleet_databaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	databaseError := errors.New("database error")

	operations := NewOperations(mockDatabase, mockDCar)

	mockDatabase.EXPECT().GetCarsForFleet(ctx, fleetID).Return(nil, databaseError)

	retCars, err := operations.GetCarsInFleet(ctx, fleetID)

	assert.ErrorIs(t, err, databaseError)
	assert.Nil(t, retCars)
}

func TestOperations_GetCarsInFleet_domainError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	domainError := errors.New("domain error")

	operations := NewOperations(mockDatabase, mockDCar)

	mockDatabase.EXPECT().GetCarsForFleet(ctx, fleetID).Return(vins, nil)
	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(nil, domainError)

	retCars, err := operations.GetCarsInFleet(ctx, fleetID)

	assert.ErrorIs(t, err, domainError)
	assert.Nil(t, retCars)
}

func TestOperations_GetCarsInFleet_carNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	mockDatabase.EXPECT().GetCarsForFleet(ctx, fleetID).Return(vins, nil)
	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(&dcar.GetCarResponse{
		HTTPResponse: &http.Response{
			StatusCode: http.StatusNotFound,
		},
	}, nil)

	retCars, err := operations.GetCarsInFleet(ctx, fleetID)

	assert.ErrorIs(t, err, fleetErrors.ErrDomainAssertion)
	assert.Nil(t, retCars)
}

func TestOperations_GetCarsInFleet_unexpectedDomainStatusCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin := "3B7HF13Y81G193584"

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	mockDatabase.EXPECT().GetCarsForFleet(ctx, fleetID).Return(vins, nil)
	mockDCar.EXPECT().GetCarWithResponse(ctx, vin).Return(&dcar.GetCarResponse{
		HTTPResponse: &http.Response{
			StatusCode: http.StatusTeapot,
		},
	}, nil)

	retCars, err := operations.GetCarsInFleet(ctx, fleetID)

	assert.ErrorIs(t, err, fleetErrors.ErrDomainAssertion)
	assert.Nil(t, retCars)
}

func TestOperations_GetCarsInFleet_secondDomainCallFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fleetID := "jJd9jb8I"
	vin1 := "3B7HF13Y81G193584"
	vin2 := "3B7HF13Y81G193585"
	multipleVins := []string{vin1, vin2}

	ctx := context.Background()

	mockDatabase := mocks.NewMockFleetDB(ctrl)
	mockDCar := mocks.NewMockClientWithResponsesInterface(ctrl)

	operations := NewOperations(mockDatabase, mockDCar)

	mockDatabase.EXPECT().GetCarsForFleet(ctx, fleetID).Return(multipleVins, nil)
	firstCall := mockDCar.EXPECT().GetCarWithResponse(ctx, vin1).
		Return(&dcar.GetCarResponse{
			JSON200: &car1,
		}, nil)
	mockDCar.EXPECT().GetCarWithResponse(ctx, vin2).
		After(firstCall).
		Return(&dcar.GetCarResponse{
			HTTPResponse: &http.Response{
				StatusCode: http.StatusNotFound,
			},
		}, nil)

	retCars, err := operations.GetCarsInFleet(ctx, fleetID)

	assert.ErrorIs(t, err, fleetErrors.ErrDomainAssertion)
	assert.Nil(t, retCars)
}
