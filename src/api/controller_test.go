package api

import (
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

var car1 = model.Car{
	Brand: "Tesla",
	Model: "Model X",
	ProductionDate: openapiTypes.Date{
		Time: time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
	},
	Vin: "3B7HF13Y81G193584",
}

var carBase2 = model.CarBase{
	Brand: "Renault",
	Model: "Megane",
	ProductionDate: openapiTypes.Date{
		Time: time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
	},
	Vin: "3N1CN7AP4FL872456",
}

var carBaseArray = []model.CarBase{
	carBase1, carBase2,
}

func TestController_GetCarsInFleet_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validFleetID := "jJd9jb8I"

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/getCars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockOperations := mocks.NewMockIOperations(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockOperations.EXPECT().GetCarsInFleet(ctx, validFleetID).Return(carBaseArray, nil)
	mockEchoContext.EXPECT().JSON(http.StatusOK, carBaseArray)

	controller := NewController(mockOperations)

	err := controller.GetCarsInFleet(mockEchoContext, validFleetID)

	assert.Nil(t, err)
}

func TestController_GetCarsInFleet_operationsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validFleetID := "jJd9jb8I"

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/getCars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockOperations := mocks.NewMockIOperations(ctrl)

	operationsError := errors.New("operations error")

	mockEchoContext.EXPECT().Request().Return(request)
	mockOperations.EXPECT().GetCarsInFleet(ctx, validFleetID).Return(nil, operationsError)

	controller := NewController(mockOperations)

	err := controller.GetCarsInFleet(mockEchoContext, validFleetID)

	assert.ErrorIs(t, err, operationsError)
}

func TestController_GetCar_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validFleetID := "jJd9jb8I"
	validVin := "3B7HF13Y81G193584"

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/getCar", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockOperations := mocks.NewMockIOperations(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockOperations.EXPECT().GetCar(ctx, validFleetID, validVin).Return(&car1, nil)
	mockEchoContext.EXPECT().JSON(http.StatusOK, &car1)

	controller := NewController(mockOperations)

	err := controller.GetCar(mockEchoContext, validFleetID, validVin)

	assert.Nil(t, err)
}

// TestController_GetCar_invalidVin_success tests that the VIN validity is not checked
// by the controller, and therefore no error is returned.
func TestController_GetCar_invalidVin_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validFleetID := "jJd9jb8I"
	invalidVin := "3B7HF13Y81G193#we"

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/getCar", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockOperations := mocks.NewMockIOperations(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockOperations.EXPECT().GetCar(ctx, validFleetID, invalidVin).Return(&car1, nil)
	mockEchoContext.EXPECT().JSON(http.StatusOK, &car1)

	controller := NewController(mockOperations)

	err := controller.GetCar(mockEchoContext, validFleetID, invalidVin)

	assert.Nil(t, err)
}

// TestController_GetCar_invalidVin_success tests that the fleet id validity is not checked
// by the controller, and therefore no error is returned.
func TestController_GetCar_invalidFleetId_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidFleetID := "jJd9I#"
	validVin := "3B7HF13Y81G193584"

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/getCar", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockOperations := mocks.NewMockIOperations(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockOperations.EXPECT().GetCar(ctx, invalidFleetID, validVin).Return(&car1, nil)
	mockEchoContext.EXPECT().JSON(http.StatusOK, &car1)

	controller := NewController(mockOperations)

	err := controller.GetCar(mockEchoContext, invalidFleetID, validVin)

	assert.Nil(t, err)
}

func TestController_GetCar_operationsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validFleetID := "jJd9jb8I"
	validVin := "3B7HF13Y81G193584"

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/getCar", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockOperations := mocks.NewMockIOperations(ctrl)

	operationsError := errors.New("operations error")

	mockEchoContext.EXPECT().Request().Return(request)
	mockOperations.EXPECT().GetCar(ctx, validFleetID, validVin).Return(nil, operationsError)

	controller := NewController(mockOperations)

	err := controller.GetCar(mockEchoContext, validFleetID, validVin)

	assert.ErrorIs(t, err, operationsError)
}

func TestController_AddCarToFleet_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validFleetID := "jJd9jb8I"
	validVin := "3B7HF13Y81G193584"

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/getCar", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockOperations := mocks.NewMockIOperations(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockOperations.EXPECT().AddCarToFleet(ctx, validFleetID, validVin).Return(&carBase1, nil)
	mockEchoContext.EXPECT().JSON(http.StatusOK, &carBase1)

	controller := NewController(mockOperations)

	err := controller.AddCarToFleet(mockEchoContext, validFleetID, validVin)

	assert.Nil(t, err)
}

func TestController_AddCarToFleet_operationsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validFleetID := "jJd9jb8I"
	validVin := "3B7HF13Y81G193584"

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/getCar", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockOperations := mocks.NewMockIOperations(ctrl)

	operationsError := errors.New("operations error")

	mockEchoContext.EXPECT().Request().Return(request)
	mockOperations.EXPECT().AddCarToFleet(ctx, validFleetID, validVin).Return(nil, operationsError)

	controller := NewController(mockOperations)

	err := controller.AddCarToFleet(mockEchoContext, validFleetID, validVin)

	assert.ErrorIs(t, err, operationsError)
}

func TestController_RemoveCar_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validFleetID := "jJd9jb8I"
	validVin := "3B7HF13Y81G193584"

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/getCar", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockOperations := mocks.NewMockIOperations(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockOperations.EXPECT().RemoveCar(ctx, validFleetID, validVin).Return(nil)
	mockEchoContext.EXPECT().NoContent(http.StatusNoContent)

	controller := NewController(mockOperations)

	err := controller.RemoveCar(mockEchoContext, validFleetID, validVin)

	assert.Nil(t, err)
}

func TestController_RemoveCar_operationsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validFleetID := "jJd9jb8I"
	validVin := "3B7HF13Y81G193584"

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/getCar", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockOperations := mocks.NewMockIOperations(ctrl)

	operationsError := errors.New("operations error")

	mockEchoContext.EXPECT().Request().Return(request)
	mockOperations.EXPECT().RemoveCar(ctx, validFleetID, validVin).Return(operationsError)

	controller := NewController(mockOperations)

	err := controller.RemoveCar(mockEchoContext, validFleetID, validVin)

	assert.ErrorIs(t, err, operationsError)
}
