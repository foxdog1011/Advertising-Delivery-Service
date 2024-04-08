# 阶段 1: 构建环境
# 使用官方的 Golang 镜像作为构建环境
FROM golang:1.18 as builder
# 设置工作目录，这是后续命令的执行上下文
WORKDIR /app
# 将当前目录下的所有文件复制到容器的 /app 目录下
COPY . .
# 下载依赖项
RUN go mod tidy
# 静态编译 Go 应用程序，禁用 CGO 并设置目标操作系统和架构
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 阶段 2: 运行环境
# 使用轻量级的 Alpine Linux 镜像作为运行环境
FROM alpine:latest
# 为了确保应用程序的兼容性，更新 Alpine 的包索引并安装 libc6-compat
# libc6-compat 提供了对 glibc 库的兼容，一些 Go 应用可能需要这个
RUN apk update && apk add --no-cache libc6-compat
# 设置工作目录
WORKDIR /root/
# 从构建阶段复制编译好的 Go 可执行文件到当前工作目录
COPY --from=builder /app/main .
# 为可执行文件设置执行权限
RUN chmod +x ./main
# 声明容器运行时应监听的端口
EXPOSE 8080
# 配置容器启动后执行的命令
CMD ["./main"]
