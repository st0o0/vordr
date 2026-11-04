## Purpose

Docker Compose deployment configuration for vordr.
## Requirements
### Requirement: Docker Compose for local development
The project SHALL include a docker-compose.yml that runs vordr as a single container with the Trivy binary bundled.

#### Scenario: Compose up
- **WHEN** running `docker compose up`
- **THEN** a single vordr container SHALL start with Docker socket and trivy-cache volume mounted

#### Scenario: Trivy cache persistence
- **WHEN** the container is restarted
- **THEN** the Trivy vulnerability database SHALL be preserved via a named volume at VORDR_TRIVY_CACHE_DIR

#### Scenario: Docker socket access
- **WHEN** vordr starts
- **THEN** it SHALL have read-only access to the Docker socket (`/var/run/docker.sock`)

#### Scenario: No Trivy sidecar
- **WHEN** inspecting the docker-compose.yml
- **THEN** there SHALL be only one service (vordr), no separate trivy server

### Requirement: Deployment documentation
The docker-compose files SHALL document server mode deployment as an option.

#### Scenario: Production compose documents server mode
- **WHEN** a user reads `docker-compose.yml`
- **THEN** comments explain how to add vordr-server and set VORDR_TRIVY_SERVER_URL

#### Scenario: Dev compose includes server mode
- **WHEN** `docker-compose.dev.yml` is used for development
- **THEN** a vordr-server service is defined for testing server mode
- **AND** vordr can be configured to use it via VORDR_TRIVY_SERVER_URL environment variable

