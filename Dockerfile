FROM golang:1.24.2-bookworm AS builder

ENV CGO_ENABLED=0

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go install -v ./...

FROM scratch AS slack-events-api

COPY --from=builder /go/bin/slack-events-api /app

CMD ["/app"]
