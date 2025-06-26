FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o order-service ./cmd/order-service

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /build/order-service ./order-service

CMD ["./order-service"]