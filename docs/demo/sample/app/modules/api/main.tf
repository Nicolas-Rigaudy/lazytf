variable "env" {
  type = string
}

variable "replicas" {
  type = number
}

resource "terraform_data" "service" {
  input = { name = "api-${var.env}", replicas = var.replicas }
}

resource "terraform_data" "replica" {
  count = var.replicas
  input = "api-${var.env}-${count.index}"
}
