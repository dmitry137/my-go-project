FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM alpine:3.19
WORKDIR /app

RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main /app/main

# Значения по умолчанию (в открытом виде, в рамках теста)
ENV PG_HOST=localhost
ENV PG_PORT=5432
ENV PG_USER=myuser
ENV PG_PASSWORD=137599
ENV PG_DBNAME=tasks_db
ENV PG_SSLMODE=disable

EXPOSE 3000
CMD ["/app/main"]