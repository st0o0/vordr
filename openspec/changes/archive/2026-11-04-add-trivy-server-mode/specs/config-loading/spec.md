## ADDED Requirements

### Requirement: Environment variable loading
Config.Load SHALL read VORDR_TRIVY_SERVER_URL from the environment. When set, it indicates server mode. When empty or unset, standalone mode is used.

#### Scenario: Server URL configured
- **WHEN** VORDR_TRIVY_SERVER_URL is set to `http://trivy:4954`
- **THEN** Config.TrivyServerURL is `http://trivy:4954`

#### Scenario: Server URL not configured
- **WHEN** VORDR_TRIVY_SERVER_URL is not set
- **THEN** Config.TrivyServerURL is empty string
- **AND** standalone mode is used (existing behavior unchanged)

#### Scenario: Server mode logged at startup
- **WHEN** vordr starts with VORDR_TRIVY_SERVER_URL set
- **THEN** the startup log includes `trivy_server_url` field with the configured URL
