# Slack Review Request Bot

A Slack bot for managing review requests with automatic reviewer assignment.

## Features

- Random reviewer assignment
- Manual reviewer selection
- Urgent mode (online reviewers only)
- Reviewer reassignment
- Thread support

## Local Setup

1. Setup environment:

```sh
./scripts/setup.sh
```

2. Run locally:

```sh
go run cmd/slack-events-api/main.go
```

3. Build with Docker:

```sh
./scripts/build.sh
```

## Usage

1. Invite the bot to your Slack channel
2. Mention the bot: `@bot-name Please review this`
3. Select reviewer option (Random/Urgent/Manual)
4. Add ✅ reaction when review is complete

## Configuration

Required environment variables:

- `SLACK_OAUTH_TOKEN`
- `SLACK_SIGNING_SECRET`

The bot uses `reviewer_map.json` for reviewer configuration.

## Deployment

Deploy to Google Cloud Run:

```sh
./scripts/deploy.sh
```

## Tech Stack

- Go 1.24.2
- Slack API
- Docker
- Google Cloud Run
- Terraform
