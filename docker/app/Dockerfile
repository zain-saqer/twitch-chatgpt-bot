# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN GOOS=linux go build -o /twitch-chatgpt-bot-server ./cmd/web

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /twitch-chatgpt-bot-server /app

ENTRYPOINT ["/app"]