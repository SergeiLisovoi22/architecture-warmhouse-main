FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o temperature-api .

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/temperature-api .
COPY --from=builder /app/docs ./docs

EXPOSE 8081
CMD ["./temperature-api"]