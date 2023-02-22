package main

import (
	"PFleetManagement/api"
	"PFleetManagement/infrastructure/database"
	"PFleetManagement/infrastructure/dcar"
	rentalManagement "PFleetManagement/infrastructure/rentalmanagement"
	"PFleetManagement/logic/fleetErrors"
	"PFleetManagement/logic/operations"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	EnvAllowOrigins            = "PFL_ALLOW_ORIGINS"
	EnvDomainServer            = "PFL_DOMAIN_SERVER"
	EnvDomainTimeout           = "PFL_DOMAIN_TIMEOUT"
	EnvRentalManagementServer  = "PFL_RENTAL_MANAGEMENT_SERVER"
	EnvRentalManagementTimeout = "PFL_RENTAL_MANAGEMENT_TIMEOUT"
)

type Config struct {
	allowOrigins            []string
	domainServer            string
	domainTimeout           time.Duration
	rentalManagementServer  string
	rentalManagementTimeout time.Duration
}

// newApp allows production as well as testing to create a new Echo instance for the API
func newApp(config *Config, fleetDb database.FleetDB) (*echo.Echo, error) {
	e := echo.New()
	e.HTTPErrorHandler = api.FleetErrorHandler

	if len(config.allowOrigins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: config.allowOrigins,
		}))
	}

	// validate incoming requests against the OpenAPI spec
	err := api.AddOpenApiValidationMiddleware(e)
	if err != nil {
		return nil, err
	}

	err = fleetDb.AddFleet(context.TODO(), "xk48jpgz") // TODO manage fleets correctly
	if err != nil && !errors.Is(err, fleetErrors.ErrFleetAlreadyExists) {
		return nil, err
	}

	dcarClient, err := dcar.NewClientWithResponses(
		config.domainServer,
		dcar.WithHTTPClient(&http.Client{Timeout: config.domainTimeout}),
	)

	rmClient, err := rentalManagement.NewClientWithResponses(
		config.rentalManagementServer,
		rentalManagement.WithHTTPClient(&http.Client{Timeout: config.rentalManagementTimeout}),
	)

	if err != nil {
		return nil, err
	}

	operationsInstance := operations.NewOperations(fleetDb, dcarClient, rmClient)
	controllerInstance := api.NewController(operationsInstance)

	api.RegisterHandlers(e, controllerInstance)

	return e, nil
}

func parseTimeout(timeoutString string) (time.Duration, error) {
	var timeout time.Duration
	if timeoutString != "" {
		var err error // declaring with := below would create separate timeout var in this scope
		timeout, err = time.ParseDuration(timeoutString)
		if err != nil {
			return 0, errors.New("invalid timeout configured")
		}
		return timeout, nil
	}
	return 5 * time.Second, nil
}

func loadConfig() (*Config, error) {
	allowOriginsString := os.Getenv(EnvAllowOrigins)
	var allowOrigins []string
	if allowOriginsString != "" {
		allowOrigins = strings.Split(allowOriginsString, ",")
	} else {
		allowOrigins = []string{}
	}

	domainServer := os.Getenv(EnvDomainServer)
	if domainServer == "" {
		return nil, errors.New("no domain server given")
	}

	rmServer := os.Getenv(EnvRentalManagementServer)
	if rmServer == "" {
		return nil, errors.New("no rental management server given")
	}

	timeoutString := os.Getenv(EnvDomainTimeout)
	domainTimeout, err := parseTimeout(timeoutString)
	if err != nil {
		return nil, err
	}

	timeoutString = os.Getenv(EnvRentalManagementTimeout)
	rmTimeout, err := parseTimeout(timeoutString)
	if err != nil {
		return nil, err
	}

	return &Config{
		allowOrigins,
		domainServer,
		domainTimeout,
		rmServer,
		rmTimeout,
	}, nil
}

func main() {
	dbConfig, err := database.LoadConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	var fleetDb database.FleetDB
	fleetDb, err = database.OpenDatabase(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	var config *Config
	config, err = loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var e *echo.Echo
	e, err = newApp(config, fleetDb)
	if err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":80"))
}
