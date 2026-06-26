FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o filekeeper ./cmd/api-server

# переносим готовый бинарник в img

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/filekeeper .

ENTRYPOINT ["./filekeeper"]

