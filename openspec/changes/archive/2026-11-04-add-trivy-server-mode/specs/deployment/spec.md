## ADDED Requirements

### Requirement: Deployment documentation
The docker-compose files SHALL document server mode deployment as an option.

#### Scenario: Production compose documents server mode
- **WHEN** a user reads `docker-compose.yml`
- **THEN** comments explain how to add vordr-server and set VORDR_TRIVY_SERVER_URL

#### Scenario: Dev compose includes server mode
- **WHEN** `docker-compose.dev.yml` is used for development
- **THEN** a vordr-server service is defined for testing server mode
- **AND** vordr can be configured to use it via VORDR_TRIVY_SERVER_URL environment variable
