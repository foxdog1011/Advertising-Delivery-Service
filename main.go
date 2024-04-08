package main

import (
	"ad-service/database"
	"ad-service/handlers"
	"log"

	"ad-service/cache"

	"github.com/gin-gonic/gin"
)

func main() {
	// 使用 database 包中的 InitDB 函数来初始化資料库连接
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err) // 使用 Fatalf 代替 Fatal 来格式化错误信息
	}
	defer db.Close() // 确保在程序结束前停止連接資料庫

	// 测试資料庫连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err) // 使用 Fatalf 代替 Fatal 来格式化错误信息
	}
	rdb := cache.InitRedis()
	router := gin.Default()

	// 創建处理函数时，傳入資料庫
	// 確保 handlers 包中的 CreateAd 和 ListAds 函数接受 *sql.DB 作为参数
	router.POST("/api/v1/ad", handlers.CreateAd(db, rdb))
	router.GET("/api/v1/ad", handlers.ListAds(db, rdb))

	router.Run(":8080") // 在 8080 port
}
