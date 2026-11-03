## MODIFIED Requirements

### Requirement: Worker pool scan target
The worker pool SHALL use `ImageRef.ScanTarget()` as the Trivy CLI argument instead of `ImageRef.String()`. ScanTarget() SHALL return ImageID when available, falling back to Image:Tag.

#### Scenario: Scan with ImageID
- **WHEN** a worker scans an image with ImageID `sha256:52d09...`
- **THEN** the Trivy CLI SHALL be called with `sha256:52d09...` as the image argument

#### Scenario: Scan without ImageID (fallback)
- **WHEN** a worker scans an image with no ImageID
- **THEN** the Trivy CLI SHALL be called with `Image:Tag` as the image argument

#### Scenario: Store key uses display name
- **WHEN** a worker updates the store after scanning
- **THEN** the store key SHALL be `Image:Tag` (not the ImageID)
