FROM golang:1.23.3-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY consumer .

RUN go build cmd/app/main.go

FROM alpine AS runner

WORKDIR app

COPY --from=build /build/main ./main

COPY --from=build /build/cmd/app/config/config.yaml ./config.yaml

COPY --from=build /build/usecases/services/docker_code_processor/Dockerfile ./Dockerfile


CMD ["/app/main", "--config=/app/config.yaml"]
