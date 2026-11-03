## MODIFIED Requirements

### Requirement: Custom Prometheus registry
The metrics package SHALL create a custom `prometheus.NewRegistry()` and MUST NOT use the global DefaultRegisterer. The Metrics struct SHALL register a custom Collector that generates vulnerability metrics dynamically from the Store snapshot.

#### Scenario: Registry isolation
- **WHEN** the metrics server starts
- **THEN** only explicitly registered vordr metrics SHALL appear on /metrics (no go_* or process_* unless explicitly added)

#### Scenario: Dynamic metrics from store
- **WHEN** the store contains scanned images
- **THEN** the custom collector SHALL call Store.Snapshot() and generate all vordr_* metrics for the current state
