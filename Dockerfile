FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN swag init --dir cmd,internal --generalInfo main/main.go
RUN go build -o main ./cmd/main


FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY .env /app/
ENV DB_HOST=db
ENV DB_NAME=subscriptions
ENV DB_PASSWORD=postgres
ENV SERVER_ADDR=0.0.0.0      
ENV SERVER_PORT=8080 

ENV GOOSE_DRIVER=postgres
ENV GOOSE_DBSTRING=postgres://postgres:postgres@db:5432/subscriptions?sslmode=disable
ENV GOOSE_MIGRATION_DIR=./migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY entrypoint.sh /app/
RUN chmod +x entrypoint.sh
EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]