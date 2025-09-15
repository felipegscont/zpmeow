FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o zpmeow-server cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/zpmeow-server .
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

EXPOSE 8080

CMD ["./zpmeow-server"]
