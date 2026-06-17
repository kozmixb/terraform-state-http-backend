# Terraform State HTTP Backend

A small HTTP backend server for storing Terraform state.

This project implements the Terraform HTTP backend protocol with support for state reads, state updates, and state locking. It is intended for self-hosted environments, homelabs, and small teams that want a simple centralized state backend without running a larger storage platform.

Docker Hub: [bencejob/terraform-state-http-backend](https://hub.docker.com/r/bencejob/terraform-state-http-backend)

## Features

- Terraform HTTP backend compatible endpoints
- State locking and unlocking
- File-based storage by default
- Optional SQLite storage
- Optional HTTP Basic Auth
- Small Docker image
- Non-root container runtime user
- Multi-architecture Docker release workflow
- Automated Go tests in GitHub Actions

## Quick Start

Run the backend with Docker:

```shell
docker run --rm \
  --name terraform-state-http-backend \
  -p 8080:8080 \
  -v "${PWD}/storage:/storage" \
  bencejob/terraform-state-http-backend:latest
```

Create a Terraform backend configuration:

```hcl
terraform {
  backend "http" {
    address        = "http://localhost:8080/example/default"
    update_method  = "POST"
    lock_address   = "http://localhost:8080/example/default"
    lock_method    = "PUT"
    unlock_address = "http://localhost:8080/example/default"
    unlock_method  = "DELETE"
  }
}
```

Initialize Terraform:

```shell
terraform init
```

Terraform HTTP backend documentation: [HashiCorp HTTP backend](https://developer.hashicorp.com/terraform/language/settings/backends/http)

## Backend Paths

The backend uses this URL pattern:

```text
/:group/:key
```

Example:

```text
http://localhost:8080/platform/network
```

With the file driver, this stores state at:

```text
storage/platform-network.json
```

Lock data is stored separately:

```text
storage/platform-network.lock
```

## Configuration

Configuration is provided with environment variables.

| Name | Default | Description |
| --- | --- | --- |
| `HTTP_PORT` | `8080` | Port the HTTP server listens on. |
| `DRIVER` | `file` | Storage driver. Supported values: `file`, `sqlite`. |
| `BASIC_AUTH_USERNAME` | unset | Enables HTTP Basic Auth when set with `BASIC_AUTH_PASSWORD`. |
| `BASIC_AUTH_PASSWORD` | unset | Enables HTTP Basic Auth when set with `BASIC_AUTH_USERNAME`. |

## Storage Drivers

### File Driver

The file driver is the default.

```shell
docker run --rm \
  -p 8080:8080 \
  -v "${PWD}/storage:/storage" \
  bencejob/terraform-state-http-backend:latest
```

State and lock files are written under `/storage`. Mount this directory to persistent storage when running in Docker.

### SQLite Driver

Set `DRIVER=sqlite` to use SQLite:

```shell
docker run --rm \
  -p 8080:8080 \
  -e DRIVER=sqlite \
  -v "${PWD}/storage:/storage" \
  bencejob/terraform-state-http-backend:latest
```

SQLite data is stored at:

```text
/storage/database.db
```

## Basic Auth

Basic Auth is disabled by default. To enable it, set both username and password:

```shell
docker run --rm \
  -p 8080:8080 \
  -e BASIC_AUTH_USERNAME=terraform \
  -e BASIC_AUTH_PASSWORD=change-me \
  -v "${PWD}/storage:/storage" \
  bencejob/terraform-state-http-backend:latest
```

Terraform backend configuration with Basic Auth:

```hcl
terraform {
  backend "http" {
    address        = "http://localhost:8080/example/default"
    update_method  = "POST"
    lock_address   = "http://localhost:8080/example/default"
    lock_method    = "PUT"
    unlock_address = "http://localhost:8080/example/default"
    unlock_method  = "DELETE"

    username = "terraform"
    password = "change-me"
  }
}
```

For production usage, prefer passing credentials through environment variables, CI secrets, or Terraform partial backend configuration rather than committing them to source control.

## HTTP API

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/:group/:key` | Read Terraform state. |
| `POST` | `/:group/:key` | Write Terraform state. |
| `PUT` | `/:group/:key` | Acquire a Terraform state lock. |
| `DELETE` | `/:group/:key` | Release a Terraform state lock. |

Common responses:

| Status | Meaning |
| --- | --- |
| `200` | Request succeeded. |
| `404` | State was not found. |
| `409` | Unlock attempted with the wrong lock ID. |
| `423` | State is already locked. |
| `500` | Storage or server error. |

## Local Development

Requirements:

- Go `1.26.0`
- Docker, if building or testing the container image

Run locally:

```shell
go run main.go
```

Run tests:

```shell
go test ./...
```

Build the Docker image:

```shell
docker build -t bencejob/terraform-state-http-backend .
```

Run the locally built image:

```shell
docker run --rm \
  -p 8080:8080 \
  -v "${PWD}/storage:/storage" \
  bencejob/terraform-state-http-backend
```

## Docker Build

Build for the local platform:

```shell
docker build -t bencejob/terraform-state-http-backend .
```

Build for multiple platforms:

```shell
docker buildx create --use --name terraform-state-builder
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --tag bencejob/terraform-state-http-backend:latest \
  .
```

## CI and Releases

The repository includes GitHub Actions workflows for:

- Running `go test ./...` on push and pull request
- Building and publishing the Docker image to Docker Hub when a GitHub release is published

Docker publishing requires these repository secrets:

```text
DOCKERHUB_USERNAME
DOCKERHUB_TOKEN
```

## Security Notes

Terraform state can contain secrets. Treat the backend storage directory, SQLite database, Docker volumes, logs, backups, and access credentials as sensitive.

Recommended practices:

- Enable Basic Auth when exposing the service beyond local development.
- Run behind HTTPS when accessed over a network.
- Restrict network access to trusted clients.
- Persist and back up `/storage`.
- Use dedicated Docker Hub tokens for CI publishing.
- Avoid committing backend credentials into Terraform files.

## License

This project is licensed under the MIT License. See [LICENSE-MIT](./LICENSE-MIT).
