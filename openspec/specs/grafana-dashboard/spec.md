## ADDED Requirements

### Requirement: Grafana dashboard JSON definition
The project SHALL include a Grafana dashboard JSON file at grafana/vordr-dashboard.json that visualizes vordr metrics.

#### Scenario: Dashboard import
- **WHEN** the JSON file is imported into Grafana
- **THEN** it SHALL create a working dashboard with panels for CVE overview, severity distribution, scan status, and DB age

#### Scenario: Dashboard panels
- **WHEN** viewing the dashboard with vordr metrics available
- **THEN** it SHALL show: stat panels (total/scanned/pending/errors), severity bar chart, per-image table, DB age gauge, and scan duration history
