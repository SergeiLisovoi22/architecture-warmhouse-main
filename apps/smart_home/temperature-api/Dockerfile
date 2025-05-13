FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o temperature-api .

FROM alpine:3.19
WORKDIR /app

# Копируем только бинарник
COPY --from=builder /app/temperature-api .

# Проверяем бинарник
RUN apk add --no-cache file && \
    file temperature-api && \
    chmod +x temperature-api

# Создаем непривилегированного пользователя
RUN adduser -D -g '' appuser && \
    chown appuser:appuser temperature-api

USER appuser

CMD ["./temperature-api"]