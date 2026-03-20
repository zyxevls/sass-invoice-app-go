FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.22

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=builder /sass-invoice-app-go /usr/local/bin/sass-invoice-app-go

EXPOSE 8080
USER app