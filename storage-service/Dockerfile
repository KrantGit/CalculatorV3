FROM golang:1.25.3

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o app ./cmd/app

CMD ["./app"]