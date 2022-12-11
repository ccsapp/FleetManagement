# FleetManagement

FleetManagement provides the functionality for the capability [Management of the Fleet](https://git.scc.kit.edu/cm-tm/cm-team/projectwork/pse/0-doc-ccs-app-v-2/-/blob/main/pages/capabilities.md) via API endpoints dedicated to individual [use cases](https://git.scc.kit.edu/cm-tm/cm-team/projectwork/pse/0-doc-ccs-app-v-2/-/blob/main/pages/use_case_diagram.md). 

For the implementation of the business logic required for the use cases, FleetManagement orchestrates [Car](https://git.scc.kit.edu/cm-tm/cm-team/projectwork/pse/domain/d-carimpl).

The provided API endpoints of FleetManagement are specified in the [API specification](https://git.scc.kit.edu/cm-tm/cm-team/projectwork/pse/application/p-fleetmanagementdesign). 
 
## Local Setup
To test FleetManagement locally, you can use the MongoDB Docker Compose setup provided in the `dev` folder.

To do so, execute the following commands:
```bash
cd dev
docker-compose up -d
```

This will start a MongoDB instance on port 27017 with a default user with admin privileges.

After that, start the Go server with the following environment variables:

| Environment Variable        | Value           |
|-----------------------------|-----------------|
| `MONGODB_DATABASE_HOST`     | localhost       |
| `MONGODB_DATABASE_NAME`     | ccsappvp2fleet  |
| `MONGODB_DATABASE_USER`     | root            |
| `MONGODB_DATABASE_PASSWORD` | example         |

## General Setup
You also need to set the environment variable `PFL_DOMAIN_SERVER` to the URL of the Car server.
`PFL_ALLOW_ORIGINS` may contain a comma-separated list of allowed origins for CORS requests.
Optionally, you can set a timeout for requests to the Car server with `PFL_DOMAIN_TIMEOUT`
([number with suffix](https://pkg.go.dev/time#ParseDuration)
"ms" for milliseconds, "s" for seconds, "m" for minutes, "h" for hours)
