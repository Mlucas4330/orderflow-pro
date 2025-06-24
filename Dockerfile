FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/orderflow-pro ./cmd/order-service

FROM builder AS tester
RUN go test -v ./...

FROM gcr.io/distroless/base-debian11 AS release

WORKDIR /

COPY --from=builder app/orderflow-pro /orderflow-pro

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/orderflow-pro"]