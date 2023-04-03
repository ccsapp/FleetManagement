# FleetManagement

FleetManagement provides the functionality for the capability
[Management of the Fleet](https://github.com/ccsapp/docs/blob/main/pages/capabilities.md)
via API endpoints dedicated to individual
[use cases](https://github.com/ccsapp/docs/blob/main/pages/use_case_diagram.md). 

For the implementation of the business logic required for the use cases, FleetManagement orchestrates 
[Car](https://github.com/ccsapp/Car) to access required data.
Therefore, it depends on the Git repository
[cargotypes](https://github.com/ccsapp/cargotypes)
to provide mappings for the JSON responses.

## Design 

[Task Processes](pages/task_processes.md) 

The provided API endpoints of FleetManagement are specified in the
[API specification](https://github.com/ccsapp/FleetManagementDesign/blob/main/openapi.yaml). 

## :warning: CORS Warning

The current HELM deployment of this microservice allows requests from all origins. This is a security risk and should
be changed in production! Currently, this is needed for the frontend development to be able to access the API.

## Local Setup Mode
To run the microservice Car locally, you can use the MongoDB setup provided in the `dev` directory.

To do so, execute the following commands:
```bash
cd dev
docker compose up -d
```

This will start a MongoDB instance on port 27031 (**non-default port** to avoid collisions with other databases) with
the correct authentication setup.

After that, start the Go server with the following environment variable set:

| Environment Variable | Value | Comment                       | 
|----------------------|-------|-------------------------------|
| `FM_LOCAL_SETUP`     | true  | Enables the local setup mode. |

You might want to set `FM_LOCAL_SETUP` in your IDE's default run configuration.
For example, in IntelliJ IDEA, you can do this [as described here](https://stackoverflow.com/a/32761503).

In the local setup mode, the microservice will use the configuration specified in `environment/localSetup.env`.
It contains the correct database connection information matching the docker compose file such that no further
configuration is required. This information will be embedded into the binary at build time.

However, you can still override the configuration by setting environment variables
described in the "Deployment or Custom Setup" section manually.

The default configuration values of local setup mode can also be found in the table of the "Deployment or Custom Setup"
section.

If the local setup mode is enabled, the integration tests (NOT the application itself) will try to detect if the
correct docker compose stack is running and will print a warning if it is not.

> After you have started the microservice in local setup mode, you can access it at
> [http://localhost:8011](http://localhost:8011).

## Deployment or Custom Setup
Do not use the local setup mode in a deployment or a custom setup, i.e. do not set the `FM_LOCAL_SETUP` environment
variable. Instead, use the following environment variables to configure the microservice:

| Environment Variable          | Local Setup Value                                     | Required for Testing? | Comment                                                                                                                                                  |
|-------------------------------|-------------------------------------------------------|-----------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `MONGODB_CONNECTION_STRING`   | mongodb://root:example@localhost:27031/ccsappvp2fleet | yes                   | The local setup uses a non-default port!                                                                                                                 |
| `MONGODB_DATABASE_NAME`       | ccsappvp2fleet                                        | yes                   |                                                                                                                                                          |
| `FM_EXPOSE_PORT`              | 8011                                                  | no                    | Optional, defaults to 80. This is the port this microservice is exposing. The local setup exposes a non-default port!                                    |
| `FM_COLLECTION_PREFIX`        | localSetup-                                           | no                    | Optional. A (unique) prefix that is prepended to every database collection of this service.                                                              |
| `FM_CAR_SERVER`               | `http://localhost:8001`                               | no                    | The URL of the Car server of the domain layer.                                                                                                           |
| `FM_RENTAL_MANAGEMENT_SERVER` | `http://localhost:8012`                               | no                    | The URL of the RentalManagement server.                                                                                                                  |
| `FM_REQUEST_TIMEOUT`          | 5s                                                    | no                    | Optional. The timeout for requests to the Car and RentalManagement server ([number with suffix](https://pkg.go.dev/time#ParseDuration)). Defaults to 5s. |
| `FM_ALLOW_ORIGINS`            | *                                                     | no                    | Optional. A comma-separated list of allowed origins for CORS requests. By default, no additional origins are allowed.                                    |                          

## Testing

### Test Setup
The Unit Tests of FleetManagement depend on automatically generated Go mocks.
You need to install [mockgen](https://github.com/golang/mock#installation) to generate them.
After the installation, execute `go generate ./...` in the `src` directory of this project.

### Running the Tests
To run the tests locally, choose the local setup mode, or use a custom setup as described above
to configure database access for the integration tests.

> **Please note:** The integration tests will ignore the `FM_COLLECTION_PREFIX` environment variable
> and use dynamically generated collection names to avoid collisions with other tests.

After that, you can run the tests using `go test ./...` in the `src` directory.
