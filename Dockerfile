FROM golang:1.24.2-bookworm AS builder

ARG SLACK_OAUTH_TOKEN
ARG SLACK_SIGNING_SECRET
ARG SLACK_REVIEWER_IDS

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build \
    -ldflags "-X github.com/himura467/slack-review-request-bot/internal/config.OAuthToken=$SLACK_OAUTH_TOKEN \
              -X github.com/himura467/slack-review-request-bot/internal/config.SigningSecret=$SLACK_SIGNING_SECRET \
              -X github.com/himura467/slack-review-request-bot/internal/config.ReviewerIDs=$SLACK_REVIEWER_IDS" \
    -o /go/bin/slack-events-api ./cmd/slack-events-api

FROM scratch AS slack-events-api

COPY --from=builder /go/bin/slack-events-api /app

CMD ["/app"]
