# Slack Review Request Bot

A Slack bot for managing review requests with automatic reviewer assignment, built with Go and deployed on Google Cloud Run.

## Features

- Random reviewer assignment
- Manual reviewer selection
- Urgent mode (online reviewers only)
- Reviewer reassignment

## Prerequisites

- Go 1.24.2
- Docker
- 1Password CLI (`op`)
- Terraform 1.10.2
- Google Cloud SDK
- AWS CLI (for Terraform state backend)

## Local Development

### 1. Setup Environment

Generate `reviewer_map.json` from 1Password:

```sh
./scripts/setup.sh
```

This script fetches reviewer configuration from 1Password and creates the required `reviewer_map.json` file.

### 2. Environment Variables

The application uses 1Password for secret management. Ensure you have access to the required vault and item:

- `OP_VAULT_NAME`: 1Password vault name (default: "Slack Review Request Bot")
- `OP_ITEM_NAME`: 1Password item name (default: "Secrets")

Required Slack credentials in 1Password:

- `SLACK_OAUTH_TOKEN`: Slack Bot User OAuth Token
- `SLACK_SIGNING_SECRET`: Slack Signing Secret

### 3. Run Locally

```sh
go run cmd/slack-events-api/main.go
```

### 4. Build Docker Image

```sh
OP_VAULT_NAME="Slack Review Request Bot" OP_ITEM_NAME="Secrets" op run --env-file app.env -- ./scripts/build.sh
```

## Usage

1. Invite the bot to your Slack channel
2. Mention the bot: `@bot-name Please review this`
3. Select reviewer option (Random/Urgent/Manual)
4. Add ✅ reaction when review is complete

## Deployment

### Infrastructure Setup

The project uses Terraform for infrastructure management with:

- Google Cloud Run for application hosting
- Google Artifact Registry for container images

### Deploy to Google Cloud Run

```sh
./scripts/deploy.sh <aws-profile> <gcloud-config-name>
```

Example:

```sh
./scripts/deploy.sh himura my-gcloud-config
```

This script:

1. Sets up AWS and Google Cloud credentials
2. Initializes Terraform
3. Applies infrastructure changes

## Project Structure

```
slack-review-request-bot/
├── cmd/slack-events-api/  # Application entry point
├── internal/              # Internal packages
│   ├── config/            # Configuration management
│   ├── domain/            # Domain models and interfaces
│   ├── infrastructure/    # External service implementations
│   ├── interface/rest/    # HTTP handlers and routing
│   └── usecase/           # Business logic
├── scripts/               # Build and deployment scripts
├── terraform/             # Infrastructure as Code
└── .github/workflows/     # CI/CD pipelines
```

## Architecture

The application follows Clean Architecture principles with dependency injection using Google Wire:

- **Domain Layer**: Business entities and repository interfaces
- **Usecase Layer**: Application business logic
- **Interface Layer**: HTTP handlers and external interfaces
- **Infrastructure Layer**: External service implementations

## Configuration

### Reviewer Configuration

The bot uses `reviewer_map.json` for reviewer assignment, which is automatically generated from 1Password during setup.

### Environment Variables

All sensitive configuration is managed through 1Password:

- `SLACK_OAUTH_TOKEN`: Slack Bot User OAuth Token
- `SLACK_SIGNING_SECRET`: Slack App Signing Secret

## Tech Stack

- **Language**: Go 1.24.2
- **Framework**: chi/v5 (HTTP router)
- **Dependency Injection**: Google Wire
- **Slack Integration**: slack-go/slack
- **Infrastructure**: Terraform
- **Cloud Platform**: Google Cloud Run
- **Container Registry**: Google Artifact Registry
- **Secret Management**: 1Password
- **State Backend**: AWS S3

## CI/CD

The project includes GitHub Actions workflows for continuous deployment, automatically building and deploying changes to Google Cloud Run.
