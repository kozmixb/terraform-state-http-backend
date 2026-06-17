# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project follows semantic versioning where practical.

## [Unreleased]

## [1.2.0] - 2026-06-17

### Changed

- Replaced Echo with the Go standard library `net/http`.
- Replaced CGO-based `github.com/mattn/go-sqlite3` with pure-Go `modernc.org/sqlite`.
- Switched the Docker runtime image to `scratch`.
- Made the Docker build CGO-free.
- Added BuildKit target platform arguments for faster multi-architecture builds.
- Kept the runtime container non-root with user `65532:65532`.
- Reduced Docker image size from about `25.5 MB` to about `14.5 MB`.

### Fixed

- SQLite no longer requires native CGO dependencies in the Docker image.
- SQLite tests now close database handles cleanly before temporary directory cleanup.

### Verified

- `go test ./...`
- Docker image build
- Terraform end-to-end test with file driver
- Terraform end-to-end test with SQLite driver
- `govulncheck`
- Trivy filesystem scan
- Trivy image scan

## [1.1.0] - 2026-06-17

### Added

- Terraform-compatible state locking and unlocking.
- Optional HTTP Basic Auth using `BASIC_AUTH_USERNAME` and `BASIC_AUTH_PASSWORD`.
- GitHub Actions workflow for Go tests on push and pull request.
- GitHub Actions workflow for Docker Hub publishing on GitHub release.
- GitHub Actions security workflow with `govulncheck`, Trivy, and dependency review.
- Expanded automated test coverage for controllers, drivers, routes, service logic, and auth.
- Professional README with setup, configuration, storage, security, and release documentation.

### Changed

- Terraform state is stored and returned as raw JSON bytes without reformatting.
- SQLite storage uses parameterized queries.
- Storage drivers return errors instead of silently ignoring failures.
- Docker image runs as a non-root user.
- Docker image uses pinned Alpine images.
- Project tooling updated to Go `1.26.4`.
- Go dependencies updated to current secure versions.

### Fixed

- Docker image build copies the compiled binary correctly.
- SQLite support works in the Docker image.
- File storage creates the `storage/` directory automatically.
- Lock conflicts return Terraform-compatible responses.
- Vulnerable indirect Go modules were updated.

### Verified

- `go test ./...`
- Docker image build
- Terraform end-to-end test with file driver
- Terraform end-to-end test with SQLite driver
- `govulncheck`
- Trivy filesystem scan
- Trivy image scan
