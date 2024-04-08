package main

import (
	"ad-service/database" // 导入自定义的 database 包
	"ad-service/handlers" // 确保 handlers 包的路径也是正确的
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 使用 database 包中的 InitDB 函数来初始化数据库连接
	db, err := database.InitDB() // 无需再在 main.go 中构建 DSN 字符串
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err) // 使用 Fatalf 代替 Fatal 来格式化错误信息
	}
	defer db.Close() // 确保在程序结束前关闭数据库连接

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err) // 使用 Fatalf 代替 Fatal 来格式化错误信息
	}

	router := gin.Default()

	// 创建处理函数时，传入数据库连接
	// 确保 handlers 包中的 CreateAd 和 ListAds 函数接受 *sql.DB 作为参数
	router.POST("/api/v1/ad", handlers.CreateAd(db))
	router.GET("/api/v1/ad", handlers.ListAds(db))

	router.Run(":8080") // 在 8080 端口启动服务
}
