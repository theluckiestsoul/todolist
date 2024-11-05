# Base stage
FROM golang:1 AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

# Test stage
FROM base AS test
COPY . .
RUN go test -v ./...

# Build stage
FROM base AS build

COPY . .
RUN CGO_ENABLED=0 go build -o todolist ./cmd/server

# Artifact stage
FROM scratch AS artifact
COPY --from=build /app/todolist .

# Final "production" image
FROM alpine:3 AS runtime

WORKDIR /app/

COPY --from=build /app/todolist .
CMD [ "/app/todolist" ]