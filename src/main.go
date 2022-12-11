package main

import (
	"PFleetManagement/infrastructure/database"
	"PFleetManagement/infrastructure/dcar"
	fleetErrors "PFleetManagement/logic/errors"
	"PFleetManagement/logic/operations"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"strings"
	"time"

	"PFleetManagement/api"
)

const (
	EnvAllowOrigins  = "PFL_ALLOW_ORIGINS"
	EnvDomainServer  = "PFL_DOMAIN_SERVER"
	EnvDomainTimeout = "PFL_DOMAIN_TIMEOUT"
)

type Config struct {
	allowOrigins  []string
	domainServer  string
	domainTimeout time.Duration
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

	timeoutString := os.Getenv(EnvDomainTimeout)
	var domainTimeout time.Duration

	if timeoutString != "" {
		var err error // declaring with := below would create separate domainTimeout var in this scope
		domainTimeout, err = time.ParseDuration(timeoutString)
		if err != nil {
			return nil, errors.New("invalid domain timeout configured")
		}
	} else {
		domainTimeout = 5 * time.Second
	}

	return &Config{
		allowOrigins,
		domainServer,
		domainTimeout,
	}, nil
}

func main() {
	e := echo.New()
	e.HTTPErrorHandler = api.FleetErrorHandler

	config, err := loadConfig()
	if err != nil {
		e.Logger.Fatal(err)
	}

	if len(config.allowOrigins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: config.allowOrigins,
		}))
	}

	fleetDb, err := database.OpenDatabase()
	if err != nil {
		e.Logger.Fatal(err)
	}
	err = fleetDb.AddFleet(context.TODO(), "xk48jpgz") // TODO manage fleets correctly
	if err != nil && !errors.Is(err, fleetErrors.ErrFleetAlreadyExists) {
		e.Logger.Fatal(err)
	}

	dcarClient, err := dcar.NewClientWithResponses(config.domainServer, func(c *dcar.Client) error {
		c.Client = &http.Client{
			Timeout: config.domainTimeout,
		}

		return nil
	})

	if err != nil {
		e.Logger.Fatal(err)
	}

	operationsInstance := operations.NewOperations(fleetDb, dcarClient)
	controllerInstance := api.NewController(operationsInstance)

	api.RegisterHandlers(e, controllerInstance)

	e.Logger.Fatal(e.Start(":80"))
}
