package handlers

import (
	"ad-service/models"
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"ad-service/cache"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type CreateAdRequest struct {
	Title      string             `json:"title" binding:"required"`
	StartAt    time.Time          `json:"startAt" binding:"required"`
	EndAt      time.Time          `json:"endAt" binding:"required,gtfield=StartAt"`
	Conditions []models.Condition `json:"conditions"`
}

func newErrorResponse(code int, message string) gin.H {
	return gin.H{"error": ErrorResponse{Code: code, Message: message}}
}

func CreateAd(db *sql.DB, rdb cache.RedisClientInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateAdRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, newErrorResponse(http.StatusBadRequest, "Invalid request data"))
			return
		}

		newAd := models.Ad{
			Title:      req.Title,
			StartAt:    req.StartAt,
			EndAt:      req.EndAt,
			CreatedAt:  time.Now(),
			Conditions: req.Conditions,
		}

		// 省略条件验证...

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to start transaction: %v", err)
			c.JSON(http.StatusInternalServerError, newErrorResponse(http.StatusInternalServerError, "Failed to process request"))
			return
		}
		defer tx.Rollback()

		adInsertStmt := `
        INSERT INTO Advertisements (title, startat, endat, createdat)
        VALUES ($1, $2, $3, $4)
        RETURNING adid`
		if err := tx.QueryRow(adInsertStmt, newAd.Title, newAd.StartAt, newAd.EndAt, newAd.CreatedAt).Scan(&newAd.ID); err != nil {
			log.Printf("Failed to insert new ad: %v", err)
			c.JSON(http.StatusInternalServerError, newErrorResponse(http.StatusInternalServerError, "Failed to insert new ad"))
			return
		}

		// 示例：使用 Redis 缓存广告标题
		ctx := context.Background()
		// 调用 rdb.Set 并获取结果
		// 直接检查从 Set 方法返回的 error
		if err := rdb.Set(ctx, "ad_title:"+newAd.Title, newAd.Title, 0); err != nil {
			log.Printf("Failed to cache ad title in Redis: %v", err)
			// 可以选择记录错误，但不一定要中断整个流程
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(http.StatusInternalServerError, newErrorResponse(http.StatusInternalServerError, "Failed to complete request"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": newAd})
	}
}
