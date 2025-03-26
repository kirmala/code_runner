FROM golang:1.23.3-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build consumer/cmd/app/main.go

FROM alpine AS runner

WORKDIR app

COPY --from=build /build/main ./main

COPY --from=build /build/consumer/cmd/app/config/config.yaml ./config.yaml

COPY --from=build /build/consumer/usecases/services/docker_code_processor/Dockerfile ./Dockerfile


CMD ["/app/main", "--config=/app/config.yaml"]
