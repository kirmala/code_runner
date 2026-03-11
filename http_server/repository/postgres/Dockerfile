FROM golang:1.25.0-alpine AS build

WORKDIR /app

RUN apk add --no-cache git

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.27.0

FROM alpine:latest AS run

COPY --from=build /go/bin/goose /usr/local/bin/goose

COPY migrations ./migrations

ENTRYPOINT ["goose", "-dir", "/migrations", "postgres"]