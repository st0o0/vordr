## ADDED Requirements

### Requirement: Docker Compose for local development
The project SHALL include a docker-compose.yml that runs vordr alongside a Trivy server instance.

#### Scenario: Compose up
- **WHEN** running `docker compose up`
- **THEN** both vordr and trivy-server containers SHALL start
- **AND** vordr SHALL be able to reach trivy at the configured address

#### Scenario: Trivy cache persistence
- **WHEN** containers are restarted
- **THEN** the Trivy vulnerability database SHALL be preserved via a named volume

#### Scenario: Docker socket access
- **WHEN** vordr starts
- **THEN** it SHALL have read-only access to the Docker socket (`/var/run/docker.sock`)
