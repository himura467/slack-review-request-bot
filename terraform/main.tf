locals {
  app_name = "slack-review-request-bot"
}

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

provider "google" {
  project = var.google_project_id
  region  = var.google_region
}

resource "google_artifact_registry_repository" "slack_review_request_bot_repo" {
  location               = var.google_region
  repository_id          = local.app_name
  format                 = "DOCKER"
  cleanup_policy_dry_run = false
  cleanup_policies {
    id     = "keep-minimum-versions"
    action = "KEEP"
    most_recent_versions {
      keep_count = 3
    }
  }
  cleanup_policies {
    id     = "delete-old-versions"
    action = "DELETE"
    condition {
      tag_state  = "ANY"
      older_than = "30d"
    }
  }
}

resource "google_project_service" "artifact_registry_api" {
  project = var.google_project_id
  service = "artifactregistry.googleapis.com"
}

resource "terraform_data" "docker_push" {
  triggers_replace = [ timestamp() ]

  provisioner "local-exec" {
    command = <<EOF
      echo "Logging in to Artifact Registry..."
      gcloud auth print-access-token --impersonate-service-account ${var.google_service_account_email} | docker login -u oauth2accesstoken --password-stdin ${var.google_region}-docker.pkg.dev

      echo "Tagging ${local.app_name} image..."
      docker tag ${local.app_name}:latest ${var.google_region}-docker.pkg.dev/${var.google_project_id}/${local.app_name}/${local.app_name}:latest

      echo "Pushing ${local.app_name} image to Artifact Registry..."
      docker push ${var.google_region}-docker.pkg.dev/${var.google_project_id}/${local.app_name}/${local.app_name}:latest
    EOF
  }

  depends_on = [
    google_artifact_registry_repository.slack_review_request_bot_repo,
    google_project_service.artifact_registry_api,
  ]
}
