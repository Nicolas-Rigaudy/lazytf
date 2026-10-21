variable "env" {
  type = string
}

resource "terraform_data" "instance" {
  input = { name = "db-${var.env}", size = var.env == "prod" ? "large" : "small" }
}

resource "terraform_data" "backup" {
  input = "db-${var.env}-nightly"
}
