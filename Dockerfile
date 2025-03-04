# 构建阶段
FROM golang:1.24.0-alpine AS builder

# 设置工作目录
WORKDIR /app

# 首先复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download && \
    go mod verify

# 复制源代码
COPY . .

# 构建应用
RUN apk add --no-cache gcc musl-dev && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o main .

# 运行阶段
FROM alpine:latest

# 设置时区
ENV TZ=Asia/Shanghai

# 安装必要的包并设置时区
RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    rm -rf /var/cache/apk/*

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8199

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8199/health || exit 1

# 运行应用
CMD ["./main"]
