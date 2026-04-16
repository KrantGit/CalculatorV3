FROM golang:1.25.3

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o calculator ./cmd/api

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/calculator .

RUN adduser -D appuser && chown -R appuser /app
USER appuser

EXPOSE 9091

CMD ["./calculator"]