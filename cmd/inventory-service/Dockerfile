FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-service ./cmd/inventory-service

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /build/inventory-service ./inventory-service

CMD ["./inventory-service"]