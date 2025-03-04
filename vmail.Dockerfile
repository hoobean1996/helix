FROM golang:1.23-alpine

# 安装sqlite和gcc编译工具链，这些是CGO必须的
RUN apk add --no-cache gcc musl-dev sqlite-dev tzdata ca-certificates

# 设置工作目录
WORKDIR /app

# 复制go模块依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译bin/vmail/main.go，启用CGO
ENV CGO_ENABLED=1
RUN go build -o /app/vmail ./bin/vmail/main.go

# 创建工作目录
RUN mkdir -p /data

# 设置时区
ENV TZ=Asia/Shanghai

# 设置启动命令
ENTRYPOINT ["/app/vmail"]