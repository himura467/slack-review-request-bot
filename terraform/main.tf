terraform {
  required_version = "1.10.2"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.31.1"
    }
  }
  
  backend "s3" {
    profile      = "himura"
    bucket       = "slack-review-request-bot-terraform-state"
    key          = "terraform.tfstate"
    region       = "ap-northeast-1"
    acl          = "private"
    encrypt      = true
    use_lockfile = true
  }
}
