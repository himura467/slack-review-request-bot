FROM golang:1.24.2-bookworm AS build

ARG SLACK_OAUTH_TOKEN
ARG SLACK_SIGNING_SECRET

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build \
    -ldflags "-X github.com/himura467/slack-review-request-bot/internal/config.OAuthToken=$SLACK_OAUTH_TOKEN \
              -X github.com/himura467/slack-review-request-bot/internal/config.SigningSecret=$SLACK_SIGNING_SECRET" \
    -o /go/bin/slack-events-api ./cmd/slack-events-api

FROM gcr.io/distroless/static-debian12 AS slack-events-api

COPY --from=build /go/bin/slack-events-api /app
COPY --from=build /go/src/app/reviewer_map.json /

CMD ["/app"]
