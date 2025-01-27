FROM golang:1.23.4 AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api-app api-gateway/cmd/http/main.go

FROM ubuntu:22.04 AS runner

WORKDIR /app

COPY --from=build /build/api-app ./
COPY --from=build /build/api-gateway/config ./config/

CMD ["./api-app"]