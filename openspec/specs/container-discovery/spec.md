## MODIFIED Requirements

### Requirement: Image reference extraction
The watcher SHALL extract image name, tag, and ImageID from Docker container metadata. The ImageID SHALL be used as the Trivy scan target. Image name and tag SHALL be used for Prometheus labels.

#### Scenario: Standard image reference
- **WHEN** a container runs `grafana/grafana:13.3.0` with ImageID `sha256:abc123`
- **THEN** ImageRef SHALL have Image `grafana/grafana`, Tag `13.3.0`, ImageID `sha256:abc123`

#### Scenario: Digest-pulled image
- **WHEN** a container runs `ghcr.io/org/app:rolling@sha256:52d09...` with ImageID `sha256:52d09...`
- **THEN** ImageRef SHALL have Image `ghcr.io/org/app`, Tag `rolling`, ImageID `sha256:52d09...`
- **AND** the `@sha256:...` part SHALL NOT appear in Image or Tag

#### Scenario: Image without tag
- **WHEN** a container runs `nginx` with ImageID `sha256:def456`
- **THEN** ImageRef SHALL have Image `nginx`, Tag `latest`, ImageID `sha256:def456`

### Requirement: Initial container listing
The watcher SHALL extract ImageID from `container.ImageID` during the initial container listing.

#### Scenario: Startup with digest-pulled containers
- **WHEN** vordr starts and containers are running with digest-pulled images
- **THEN** each ImageRef SHALL have the ImageID populated from the Docker API

### Requirement: Docker event stream
The watcher SHALL extract ImageID for new containers discovered via events.

#### Scenario: New container starts
- **WHEN** a container start event is received
- **THEN** the watcher SHALL inspect the container to obtain the ImageID
