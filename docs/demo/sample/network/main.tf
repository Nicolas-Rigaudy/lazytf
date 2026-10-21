# Second sample project: no var files and no backend config, so LazyTF lists
# a single "(no var-file)" environment and inits with the local backend.
resource "terraform_data" "vpc" {
  input = "10.0.0.0/16"
}
