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

	// HTTP responses to HEAD requests must not include a body
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
	// if we already have a response, don't do anything
	if ctx.Response().Committed {
		return
	}

	if stdErrors.Is(err, fleetErrors.ErrFleetNotFound) || stdErrors.Is(err, fleetErrors.ErrCarNotFound) ||
		stdErrors.Is(err, fleetErrors.ErrCarNotInFleet) {

		messageResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	if stdErrors.Is(err, fleetErrors.ErrCarAlreadyInFleet) {
		messageResponse(ctx, http.StatusNoContent, err.Error())
		return
	}

	if stdErrors.Is(err, fleetErrors.ErrInvalidFleetId) || stdErrors.Is(err, fleetErrors.ErrInvalidVin) {
		messageResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if httpErr, ok := err.(*echo.HTTPError); ok {
		// use the HTTP standard error message
		message := http.StatusText(httpErr.Code)

		// except there is a custom message
		if containedMessage, ok := httpErr.Message.(string); ok {
			message = containedMessage
		}

		messageResponse(ctx, httpErr.Code, message)
		return
	}

	// unexpected error
	ctx.Logger().Error(err)
	messageResponse(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
