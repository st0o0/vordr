## ADDED Requirements

### Requirement: Embedded dashboard served from binary
The Vordr binary SHALL serve a self-contained HTML dashboard via `go:embed`. The `internal/ui` package SHALL expose a `Handler()` function returning an `http.Handler` that serves the embedded HTML file.

#### Scenario: Dashboard handler serves HTML
- **WHEN** `ui.Handler()` receives a GET request
- **THEN** the response SHALL have status 200, Content-Type `text/html`, and contain the dashboard HTML

#### Scenario: Embedded file compiled into binary
- **WHEN** the binary is built
- **THEN** `dashboard.html` SHALL be embedded via `//go:embed` and served without external file dependencies

### Requirement: Dashboard displays vulnerability summary
The dashboard SHALL fetch `/details` and display aggregate severity counts as summary tiles.

#### Scenario: Summary tiles shown
- **WHEN** the dashboard loads and `/details` returns image data
- **THEN** the dashboard SHALL display total vulnerability count and per-severity counts (Critical, High, Medium, Low)

#### Scenario: No images tracked
- **WHEN** `/details` returns an empty array
- **THEN** the summary tiles SHALL show zero counts

### Requirement: Dashboard lists images with status
The dashboard SHALL display all tracked images in a list showing image name, scan status, vulnerability count, and last scan time.

#### Scenario: Image list populated
- **WHEN** `/details` returns image entries
- **THEN** each image SHALL be displayed with its name, tag, status, vulnerability count, and relative scan time

#### Scenario: Pending image display
- **WHEN** an image has status "pending"
- **THEN** the vulnerability count SHALL show "-" and last scan time SHALL show "-"

### Requirement: Per-image vulnerability detail view
The dashboard SHALL support expanding an image row to show its individual vulnerabilities.

#### Scenario: Expand image details
- **WHEN** a user clicks on an image row
- **THEN** the row SHALL expand to show a table of vulnerabilities with CVE ID, package name, installed version, fixed version, and severity

#### Scenario: Collapse image details
- **WHEN** a user clicks on an already-expanded image row
- **THEN** the vulnerability detail table SHALL collapse

### Requirement: Dashboard auto-refresh
The dashboard SHALL periodically refresh data from the server without requiring manual page reload.

#### Scenario: Periodic data refresh
- **WHEN** 30 seconds have elapsed since the last fetch
- **THEN** the dashboard SHALL re-fetch `/details` and `/healthz` and update the displayed data

#### Scenario: Last refresh timestamp
- **WHEN** data is refreshed
- **THEN** the dashboard SHALL display the time of the last successful refresh

### Requirement: Dashboard displays server health info
The dashboard SHALL fetch `/healthz` and display version and uptime in the header.

#### Scenario: Health info in header
- **WHEN** `/healthz` returns successfully
- **THEN** the dashboard header SHALL show the Vordr version and uptime

#### Scenario: Health status indicator
- **WHEN** `/healthz` returns status "ok"
- **THEN** the dashboard SHALL display a green health indicator
