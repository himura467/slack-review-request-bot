variable "google_service_account_email" {
  description = "The email of the Google service account to use for authentication"
  type        = string
}

variable "google_project_id" {
  description = "The ID of the Google Cloud project"
  type        = string
}

variable "google_region" {
  description = "The Google Cloud region to deploy resources"
  type        = string
  default     = "asia-northeast1"
}
