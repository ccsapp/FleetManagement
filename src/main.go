package main

import (
	"PFleetManagement/infrastructure/database"
	"PFleetManagement/infrastructure/dcar"
	"PFleetManagement/logic/operations"
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

const (
	ExitConfigErr       = 9
	ExitDomainClientErr = 10
	ExitDatabaseOpenErr = 11
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
		// TODO make env run
		allowOrigins = []string{"*.cloud.iai.kit.edu"}
	}

	domainServer := os.Getenv(EnvDomainServer)
	if domainServer == "" {
		domainServer = "https://cm-d-carimplementation-mock-api.cloud.iai.kit.edu"

		// TODO make env run
		// return nil, errors.New("no domain server given")
	}

	timeoutString := os.Getenv(EnvDomainTimeout)
	domainTimeout, err := time.ParseDuration(timeoutString)
	if err != nil {
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
		os.Exit(ExitConfigErr)
	}

	if len(config.allowOrigins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			Skipper:      middleware.DefaultSkipper,
			AllowOrigins: config.allowOrigins,
			AllowMethods: []string{
				http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete,
			},
		}))
	}

	fleetDb := database.OpenDatabase()
	err = fleetDb.AddFleet("xk48jpgz") // TODO manage fleets correctly
	if err != nil {
		e.Logger.Fatal(err)
		os.Exit(ExitDatabaseOpenErr)
	}

	dcarClient, err := dcar.NewClientWithResponses(config.domainServer, func(c *dcar.Client) error {
		c.Client = &http.Client{
			Timeout: config.domainTimeout,
		}

		return nil
	})

	if err != nil {
		e.Logger.Fatal(err)
		os.Exit(ExitDomainClientErr)
	}

	operationsInstance := operations.NewOperations(fleetDb, dcarClient)
	controllerInstance := api.NewController(operationsInstance)

	api.RegisterHandlers(e, controllerInstance)

	e.Logger.Fatal(e.Start(":80"))
}
