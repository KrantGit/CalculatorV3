FROM golang:1.25.3

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o worker ./cmd/worker

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/worker .
COPY --from=builder /app/migrations ./migrations

RUN adduser -D appuser && chown -R appuser /app
USER appuser

CMD ["./worker"]