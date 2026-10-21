# Sample project for the README demo. terraform_data is built into Terraform,
# so init downloads nothing and plan/apply touch no cloud account.
terraform {
  backend "local" {}
}

variable "env" {
  type = string
}

variable "replicas" {
  type = number
}

module "api" {
  source   = "./modules/api"
  env      = var.env
  replicas = var.replicas
}

module "db" {
  source = "./modules/db"
  env    = var.env
}

resource "terraform_data" "dns" {
  input = "app.${var.env}.example.com"
}
