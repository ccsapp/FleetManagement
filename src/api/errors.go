package api

import (
	"PFleetManagement/logic/fleetErrors"
	"errors"
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
		// if the request was anything but head, return the message in a json object
		err = ctx.JSON(code, errorContent{
			message,
		})
	}

	// if the response could not be written, log why
	if err != nil {
		ctx.Logger().Error(err)
	}
}

// The FleetErrorHandler maps errors returned by the model (i.e. errors from logic/errors) to HTTP
// responses describing the error. This includes determining the correct status codes as defined in
// the specification.
func FleetErrorHandler(err error, ctx echo.Context) {
	// if we already have a response, don't do anything
	if ctx.Response().Committed {
		return
	}

	// "... not found" errors result in a 404 response
	if errors.Is(err, fleetErrors.ErrFleetNotFound) || errors.Is(err, fleetErrors.ErrCarNotFound) ||
		errors.Is(err, fleetErrors.ErrCarNotInFleet) {

		messageResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	// [logic/errors.ErrCarAlreadyInFleet] is not considered a failure (b/c of idempotency)
	if errors.Is(err, fleetErrors.ErrCarAlreadyInFleet) {
		messageResponse(ctx, http.StatusNoContent, err.Error())
		return
	}

	// invalid fleet id or vin, being an invalid/bad request, results in 400
	if errors.Is(err, fleetErrors.ErrInvalidFleetId) || errors.Is(err, fleetErrors.ErrInvalidVin) {
		messageResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// if the returned error is explicitly an HTTP error, return the contained status code
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

	// any other error, which was not handled above; this can be any error returned from library calls
	ctx.Logger().Error(err)
	messageResponse(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
