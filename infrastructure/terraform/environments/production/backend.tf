terraform {
  backend "s3" {
    bucket      = "medvault-terraform-state-836734448013"
    key         = "production/terraform.tfstate"
    region      = "us-east-1"
    use_lockfile = true
    encrypt     = true
  }
}
