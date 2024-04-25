# Terraform State: http-backend

## Usage

### Run with go
```shell
go run main.go
```

### Run with docker
```shell
docker build -t 'terraform-state-http-backend' .
docker run --rm -it --name 'terraform-state-http-backend' -v ${pwd}/storage:/storage -p 8080:8080  terraform-state-http-backend
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
