## ADDED Requirements

### Requirement: Details JSON endpoint
The HTTP server SHALL expose GET /details returning a JSON array of all tracked images with their vulnerability details.

#### Scenario: All images
- **WHEN** GET /details is called
- **THEN** the server SHALL return a JSON array with one entry per tracked image, including image, tag, status, scanned_at, and vulnerabilities list

#### Scenario: Filter by image
- **WHEN** GET /details?image=nginx:latest is called
- **THEN** the server SHALL return only the entry for that image, or 404 if not found

#### Scenario: Pending image
- **WHEN** an image is in Pending state
- **THEN** the vulnerabilities list SHALL be null/empty and status SHALL be "pending"

#### Scenario: Scanned image with CVEs
- **WHEN** an image has been scanned with vulnerabilities
- **THEN** each vulnerability SHALL include id, package, installed version, fixed version, and severity
