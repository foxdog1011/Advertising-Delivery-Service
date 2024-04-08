# 階段 1: 構建環境
# 使用官方的 Golang 鏡像作為構建環境
FROM golang:1.18 as builder
# 設置工作目錄，這是後續命令的執行上下文
WORKDIR /app
# 將當前目錄下的所有文件複製到容器的 /app 目錄下
COPY . .
# 下載依賴項
RUN go mod tidy
# 靜態編譯 Go 應用程序，禁用 CGO 並設置目標操作系統和架構
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 階段 2: 運行環境
# 使用輕量級的 Alpine Linux 鏡像作為運行環境
FROM alpine:latest
# 為了確保應用程序的兼容性，更新 Alpine 的包索引並安裝 libc6-compat
# libc6-compat 提供了對 glibc 库的兼容，一些 Go 應用可能需要這個
RUN apk update && apk add --no-cache libc6-compat
# 設置工作目錄
WORKDIR /root/
# 從構建階段複製編譯好的 Go 可執行文件到當前工作目錄
COPY --from=builder /app/main .
# 為可執行文件設置執行權限
RUN chmod +x ./main
# 宣告容器運行時應監聽的端口
EXPOSE 8080
# 配置容器啟動後執行的命令
CMD ["./main"]
