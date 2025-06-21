variable "google_project_id" {
  description = "The ID of the Google Cloud project"
  type        = string
  default     = "slack-review-request-bot"
}

variable "google_region" {
  description = "The Google Cloud region to deploy resources"
  type        = string
  default     = "asia-northeast1"
}

variable "google_service_account_email" {
  description = "The email of the Google service account to use for authentication"
  type        = string
  default     = "mitarashidango0927@gmail.com"
}

variable "github_org" {
  description = "GitHub organization name"
  type        = string
}

variable "github_repo" {
  description = "GitHub repository name"
  type        = string
}
