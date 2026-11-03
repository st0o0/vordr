## MODIFIED Requirements

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
