FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache gcc musl-dev && \
    go mod download && \
    go mod verify && \
    go build -ldflags="-s -w" -o main .

FROM alpine:latest
ENV TZ=Asia/Shanghai
RUN apk add --no-cache alpine-conf ca-certificates  && \
    /usr/sbin/setup-timezone -z Asia/Shanghai && \
    apk del alpine-conf && \
    rm -rf /var/cache/apk/*
COPY --from=builder /app/main /app/main
CMD ["sh", "-c", "/app/main"]
EXPOSE 8199
