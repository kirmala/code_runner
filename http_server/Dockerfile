FROM golang:1.23.3-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY http_server .

RUN apk add --no-cache make

RUN go build cmd/app/main.go

FROM alpine AS runner

WORKDIR app

RUN apk add --no-cache curl

COPY --from=build /build/main ./main

COPY --from=build /build/cmd/app/config/config.yaml ./config.yaml 

CMD ["/app/main", "--config=/app/config.yaml"]
