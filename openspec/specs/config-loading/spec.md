## Purpose

Environment variable based configuration for vordr runtime settings.
## Requirements
### Requirement: Environment variable based configuration
The config SHALL add three new environment variables.

#### Scenario: Default values for new fields
- **WHEN** no new environment variables are set
- **THEN** Load SHALL return ExcludeImages as empty slice, AlertSeverity as "HIGH", RetryMaxBackoff as 15m

#### Scenario: VORDR_EXCLUDE_IMAGES
- **WHEN** VORDR_EXCLUDE_IMAGES is `aquasec/trivy:*,*devcontainer*`
- **THEN** ExcludeImages SHALL be ["aquasec/trivy:*", "*devcontainer*"]

#### Scenario: VORDR_ALERT_SEVERITY validation
- **WHEN** VORDR_ALERT_SEVERITY is set to "INVALID"
- **THEN** Load SHALL return an error

#### Scenario: VORDR_RETRY_MAX_BACKOFF
- **WHEN** VORDR_RETRY_MAX_BACKOFF is `20m`
- **THEN** RetryMaxBackoff SHALL be 20 minutes

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

