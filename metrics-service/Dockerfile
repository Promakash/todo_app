FROM golang:1.23.4 AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o metrics-app metrics-service/cmd/main.go

FROM ubuntu:22.04 AS runner

WORKDIR /app

COPY --from=build /build/metrics-app ./
COPY --from=build /build/metrics-service/config ./config/

CMD ["./metrics-app"]