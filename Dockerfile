FROM golang:alpine AS builder
WORKDIR /app

# 先复制依赖文件，利用Docker缓存机制
COPY go.mod go.sum ./
RUN go mod download

# 设置环境变量并复制源代码
ENV CGO_ENABLED=0
COPY . .
RUN go build -ldflags="-s -w" -o main .

FROM alpine
ENV TZ=Asia/Shanghai
RUN apk add --no-cache alpine-conf ca-certificates  && \
    /usr/sbin/setup-timezone -z Asia/Shanghai && \
    apk del alpine-conf && \
    rm -rf /var/cache/apk/*
COPY --from=builder /app/main /app/main
CMD /app/main
EXPOSE 8299