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
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "6.40.0"
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

provider "google-beta" {
  project = var.google_project_id
  region  = var.google_region
}

data "google_project" "this" {
  project_id = var.google_project_id
}

resource "google_artifact_registry_repository" "this" {
  location               = var.google_region
  repository_id          = local.app_name
  format                 = "DOCKER"
  cleanup_policy_dry_run = false
  docker_config {
    immutable_tags = false
  }
  cleanup_policies {
    id     = "keep-minimum-versions"
    action = "KEEP"
    most_recent_versions {
      keep_count = 2
    }
  }
  cleanup_policies {
    id     = "delete-old-versions"
    action = "DELETE"
    condition {
      tag_state  = "ANY"
      older_than = "24h"
    }
  }
}

resource "google_project_service" "artifact_registry" {
  project = var.google_project_id
  service = "artifactregistry.googleapis.com"
}

resource "google_project_service" "cloud_run_admin" {
  project = var.google_project_id
  service = "run.googleapis.com"
}

resource "google_cloud_run_v2_service" "slack_review_request_bot" {
  name                = local.app_name
  location            = var.google_region
  deletion_protection = false
  template {
    containers {
      image = "${var.google_region}-docker.pkg.dev/${var.google_project_id}/${local.app_name}/${local.app_name}:latest"
    }
    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }
  }
  depends_on = [google_project_service.cloud_run_admin]
  lifecycle {
    ignore_changes = [
      template[0].containers[0].image,
      client,
      client_version
    ]
  }
}

data "google_iam_policy" "cloud_run_invoker" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_v2_service_iam_policy" "cloud_run_invoker" {
  project     = google_cloud_run_v2_service.slack_review_request_bot.project
  location    = google_cloud_run_v2_service.slack_review_request_bot.location
  name        = google_cloud_run_v2_service.slack_review_request_bot.name
  policy_data = data.google_iam_policy.cloud_run_invoker.policy_data
}
