# Service folder name inside ./cmd
ARG SERVICE_PATH
##
## Stage: Base
##
FROM golang:1.22-alpine AS base
WORKDIR /app
ADD . /app
RUN go mod vendor \
    && go mod download
##
## Stage: Development
##
FROM base AS development
ARG SERVICE_PATH
ARG DEBUG_PORT
RUN go install github.com/air-verse/air@latest
RUN apk add build-base \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && sed "s/{SERVICE_PATH}/${SERVICE_PATH}/g; s/{DEBUG_PORT}/${DEBUG_PORT}/g" .air.toml > /root/.air.toml
ENTRYPOINT ["air", "-c", "/root/.air.toml"]

##
## Stage: Build
##
FROM base AS builder
ARG SERVICE_PATH
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o go-service "./cmd/${SERVICE_PATH}/main.go"

##
## Stage: Release
##
FROM alpine:latest as release
WORKDIR /release
COPY ./.env ./.env
COPY --from=builder /app/config ./config
COPY --from=builder /app/go-service .
ENTRYPOINT ["/release/go-service"]
