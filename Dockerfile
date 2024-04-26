FROM golang:alpine AS builder

ENV ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]