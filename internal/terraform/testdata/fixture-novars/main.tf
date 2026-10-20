# Fixture for the integration test (go test -tags integration): local
# backend, no variables, no tfvars. A provider is declared because a project
# without one gets no .terraform directory from init.
terraform {
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.2"
    }
  }
}

resource "null_resource" "x" {}
