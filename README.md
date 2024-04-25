# Terraform State: http-backend

A Fast, Minimal terraform state backend server.
An easy centralised solution for homelabs

Dockerhub :[image](https://hub.docker.com/r/bencejob/terraform-state-http-backend)

## Usage

### Run with go
```shell
go run main.go
```

### Run with docker
```shell
docker run --rm -it --name 'bencejob/terraform-state-http-backend' -v ${pwd}/storage:/storage -p 8080:8080  bencejob/terraform-state-http-backend
```

Terraform http backend documentation: [docs](https://developer.hashicorp.com/terraform/language/settings/backends/http)

Add the following to your terraform `backend.tf` file:
```hcl
terraform {
  backend "http" {
    address         = "http://localhost:8080/{GROUP_NAME}/{KEY_NAME}"
    update_method   = "POST"
    lock_address    = "http://localhost:8080/{GROUP_NAME}/{KEY_NAME}"
    lock_method     = "PUT"
    unlock_address  = "http://localhost:8080/{GROUP_NAME}/{KEY_NAME}"
    unlock_method   = "DELETE"
  }
}
```

## Build

Build on local machine
```shell
docker build -t 'bencejob/terraform-state-http-backend' .
```

Build for multiple platforms
```shell
docker buildx create --use --name mybuild
docker buildx build --platform linux/386,linux/amd64,linux/arm64 --tag 'bencejob/terraform-state-http-backend' .
```
