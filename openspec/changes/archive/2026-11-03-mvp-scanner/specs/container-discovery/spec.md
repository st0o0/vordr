## ADDED Requirements

### Requirement: Initial container listing
The watcher SHALL list all running Docker containers on startup and extract their image references. Each unique image SHALL be added to the store and sent to the scan channel.

#### Scenario: Startup with running containers
- **WHEN** vordr starts and Docker has 5 running containers using 3 unique images
- **THEN** the store SHALL contain 3 image entries with status Pending
- **AND** 3 ImageRef values SHALL be sent to the scan channel with staggered timing

#### Scenario: Startup with no containers
- **WHEN** vordr starts and Docker has no running containers
- **THEN** the store SHALL be empty
- **AND** no values SHALL be sent to the scan channel

### Requirement: Docker event stream
The watcher SHALL subscribe to Docker container events (start, die) and react in real-time. On container start, the image SHALL be added to the store if new. On container die, the container SHALL be removed from the image's container list.

#### Scenario: New container starts with unknown image
- **WHEN** a container starts with an image not yet in the store
- **THEN** the image SHALL be added to the store with status Pending
- **AND** an ImageRef SHALL be sent to the scan channel

#### Scenario: New container starts with known image
- **WHEN** a container starts with an image already in the store
- **THEN** the container SHALL be added to the image's container list
- **AND** no scan SHALL be triggered

#### Scenario: Last container for an image stops
- **WHEN** the last container using a specific image dies
- **THEN** the image SHALL be removed from the store

#### Scenario: One of multiple containers stops
- **WHEN** a container dies but other containers still use the same image
- **THEN** the container SHALL be removed from the image's container list
- **AND** the image SHALL remain in the store

### Requirement: Image reference extraction
The watcher SHALL extract image name and tag from Docker container metadata. The image reference SHALL be the key used in the store.

#### Scenario: Standard image reference
- **WHEN** a container runs `grafana/grafana:13.3.0`
- **THEN** ImageRef SHALL have Image `grafana/grafana` and Tag `13.3.0`

#### Scenario: Image without tag
- **WHEN** a container runs `nginx` (no tag specified)
- **THEN** ImageRef SHALL have Image `nginx` and Tag `latest`

### Requirement: Graceful shutdown
The watcher SHALL stop the event stream and return when its context is cancelled.

#### Scenario: Context cancellation during event stream
- **WHEN** the context is cancelled while watching events
- **THEN** the watcher SHALL stop reading events and return nil
