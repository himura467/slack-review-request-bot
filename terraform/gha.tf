# https://github.com/google-github-actions/auth#preferred-direct-workload-identity-federation

resource "google_project_service" "iam" {
  project = var.google_project_id
  service = "iam.googleapis.com"
}

resource "google_iam_workload_identity_pool" "gha" {
  provider                  = google-beta
  workload_identity_pool_id = "gha-pool"
  display_name              = "GHA Workload Identity Pool"
  depends_on                = [google_project_service.iam]
}

resource "google_iam_workload_identity_pool_provider" "gha" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.gha.workload_identity_pool_id
  workload_identity_pool_provider_id = "gha-provider"
  display_name                       = "GHA Provider"
  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.repository" = "assertion.repository"
  }
  attribute_condition = <<EOT
    assertion.repository == "${var.github_org}/${var.github_repo}" &&
    assertion.event_name == "pull_request" &&
    assertion.base_ref == "refs/heads/main"
  EOT
  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

resource "google_project_iam_member" "gha_gar_writer" {
  project = var.google_project_id
  role    = "roles/artifactregistry.writer"
  member  = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.gha.name}/attribute.repository/${var.github_org}/${var.github_repo}"
}

resource "google_project_iam_member" "gha_cloud_run_developer" {
  project = var.google_project_id
  role    = "roles/run.developer"
  member  = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.gha.name}/attribute.repository/${var.github_org}/${var.github_repo}"
}

resource "google_service_account_iam_member" "gha_act_as_compute_sa" {
  service_account_id = "projects/${var.google_project_id}/serviceAccounts/${data.google_project.this.number}-compute@developer.gserviceaccount.com"
  role               = "roles/iam.serviceAccountUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.gha.name}/attribute.repository/${var.github_org}/${var.github_repo}"
}
