FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -ldflags="-s -w" -o main .

FROM alpine
ENV TZ=Asia/Shanghai
RUN apk add --no-cache alpine-conf ca-certificates  && \
    /usr/sbin/setup-timezone -z Asia/Shanghai && \
    apk del alpine-conf && \
    rm -rf /var/cache/apk/*
COPY --from=builder /app/main /app/main
CMD /app/main
EXPOSE 8199