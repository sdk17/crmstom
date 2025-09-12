# Многоэтапная сборка для оптимизации размера образа
FROM golang:1.21-alpine AS builder

# Установка необходимых пакетов
RUN apk add --no-cache git ca-certificates tzdata

# Установка рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Финальный образ
FROM alpine:latest

# Установка CA сертификатов и временной зоны
RUN apk --no-cache add ca-certificates tzdata

# Создание пользователя для безопасности
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Установка рабочей директории
WORKDIR /root/

# Копирование исполняемого файла из builder
COPY --from=builder /app/main .

# Копирование статических файлов
COPY --from=builder /app/static ./static

# Изменение владельца файлов
RUN chown -R appuser:appgroup /root/

# Переключение на непривилегированного пользователя
USER appuser

# Открытие порта
EXPOSE 8080

# Команда запуска
CMD ["./main"]
