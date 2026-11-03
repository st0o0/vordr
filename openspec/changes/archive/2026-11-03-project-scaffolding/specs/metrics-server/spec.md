## ADDED Requirements

### Requirement: Custom Prometheus registry
The metrics package SHALL create a custom `prometheus.NewRegistry()` and MUST NOT use the global DefaultRegisterer.

#### Scenario: Registry isolation
- **WHEN** the metrics server starts
- **THEN** only explicitly registered vordr metrics SHALL appear on /metrics (no go_* or process_* unless explicitly added)

### Requirement: Metrics HTTP server
The metrics package SHALL provide an HTTP server on a configurable address (default `:9550`) with `/metrics` and `/healthz` endpoints.

#### Scenario: Metrics endpoint
- **WHEN** HTTP GET /metrics is called
- **THEN** the server SHALL return Prometheus text format with registered metrics and status 200

#### Scenario: Health endpoint
- **WHEN** HTTP GET /healthz is called
- **THEN** the server SHALL return a JSON response with status "ok" and HTTP 200

#### Scenario: Server timeout
- **WHEN** the server is created
- **THEN** ReadHeaderTimeout SHALL be set to 5 seconds

### Requirement: Graceful server shutdown
The metrics server SHALL shut down gracefully when its context is cancelled.

#### Scenario: Context cancellation
- **WHEN** the parent context is cancelled
- **THEN** the server SHALL stop accepting new connections and drain existing ones
