package api

import (
	fleetErrors "PFleetManagement/logic/errors"
	stdErrors "errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type errorContent struct {
	Message string `json:"message"`
}

func messageResponse(ctx echo.Context, code int, message string) {
	var err error
	if ctx.Request().Method == http.MethodHead {
		err = ctx.NoContent(code)
	} else {
		err = ctx.JSON(code, errorContent{
			message,
		})
	}

	if err != nil {
		ctx.Logger().Error(err)
	}
}

func FleetErrorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	if stdErrors.Is(err, fleetErrors.ErrFleetNotFound) ||
		stdErrors.Is(err, fleetErrors.ErrCarNotFound) ||
		stdErrors.Is(err, fleetErrors.ErrCarNotInFleet) {
		messageResponse(ctx, http.StatusNotFound, err.Error())
	} else if stdErrors.Is(err, fleetErrors.ErrCarAlreadyInFleet) {
		messageResponse(ctx, http.StatusNoContent, err.Error())
	} else if stdErrors.Is(err, fleetErrors.ErrInvalidFleetId) ||
		stdErrors.Is(err, fleetErrors.ErrInvalidVin) {
		messageResponse(ctx, http.StatusBadRequest, err.Error())
	} else if httpErr, ok := err.(*echo.HTTPError); ok {
		message := http.StatusText(httpErr.Code)
		if containedMessage, ok := httpErr.Message.(string); ok {
			message = containedMessage
		}

		messageResponse(ctx, httpErr.Code, message)
	} else {
		ctx.Logger().Error(err)
		messageResponse(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}
