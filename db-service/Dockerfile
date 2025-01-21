FROM golang:1.23.4-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o db-app db-service/cmd/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=build /build/db-app ./
COPY --from=build /build/db-service/config ./config/

CMD ["./db-app"]