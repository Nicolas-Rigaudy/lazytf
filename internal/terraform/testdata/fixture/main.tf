# Fixture for the integration test (go test -tags integration).
terraform {
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.2"
    }
  }
  backend "local" {}
}

variable "name" {
  type = string
}

resource "null_resource" "x" {
  triggers = {
    name = var.name
  }
}
