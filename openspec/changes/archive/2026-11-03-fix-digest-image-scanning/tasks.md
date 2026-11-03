## 1. ImageRef Changes

- [x] 1.1 Add `ImageID string` field to `ImageRef` in `internal/scanner/types.go`
- [x] 1.2 Add `ScanTarget() string` method — returns ImageID if set, otherwise String()

## 2. ParseImageRef Digest Handling

- [x] 2.1 Update `ParseImageRef` in `internal/watcher/watcher.go` to strip `@sha256:...` before tag parsing
- [x] 2.2 Add test cases for digest references (`image:tag@sha256:...`, `image@sha256:...`)

## 3. Watcher ImageID Extraction

- [x] 3.1 Update initial container listing to set `ref.ImageID` from `container.ImageID`
- [x] 3.2 Update event handler to inspect container for ImageID on start events
- [x] 3.3 Add `ImageID` parameter to `ParseImageRef` or set it separately after parsing

## 4. Worker Pool Scan Target

- [x] 4.1 Change `pool.scan()` to use `ref.ScanTarget()` instead of `ref.String()` for Trivy CLI argument
- [x] 4.2 Keep `ref.String()` for store key and log messages

## 5. Verification

- [x] 5.1 Run `go build ./...` and `go vet ./...`
- [x] 5.2 Run `go test ./...` — all tests pass
- [x] 5.3 Run `docker compose up`, verify digest-pulled images scan successfully
