FROM golang:1.25.0-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build consumer/cmd/app/main.go

FROM alpine AS runner

WORKDIR app

COPY --from=build /build/main ./main

COPY --from=build /build/consumer/cmd/app/config/config.yaml ./config.yaml

COPY --from=build /build/consumer/service/docker/runner/Dockerfile ./Dockerfile


CMD ["/app/main", "--config=/app/config.yaml"]
