# Сначала этап сборки
FROM golang:latest AS builder

WORKDIR /app
COPY . .

RUN go build -o main pionsfu.go

# Затем этап запуска приложения
FROM golang:latest

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
