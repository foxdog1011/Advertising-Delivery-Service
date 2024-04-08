package handlers

import (
	"ad-service/cache"
	"ad-service/models"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func ListAds(db *sql.DB, rdb cache.RedisClientInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 构建缓存键
		cacheKey := "ads_list:" + c.Request.URL.RawQuery // 使用完整的查询字符串作为缓存键的一部分
		// 尝试从 Redis 中获取缓存的广告列表
		cachedAds, err := rdb.Get(context.Background(), cacheKey)
		if err == nil {
			// 如果成功获取到缓存，则直接返回缓存的结果
			c.Header("Content-Type", "application/json")
			c.String(http.StatusOK, cachedAds)
			return
		}
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		limit = max(min(limit, 100), 1)

		now := time.Now()
		params := []interface{}{now, now}

		query := "SELECT adid, title, startat, endat FROM Advertisements WHERE startat <= $1 AND endat >= $2"

		// 处理多选条件：国家、平台（性别和年龄在这个示例中被假定为单选，如需要支持多选，可按照国家的处理方式进行调整）
		if countries, ok := c.GetQueryArray("country"); ok && len(countries) > 0 {
			query += " AND country = ANY($" + strconv.Itoa(len(params)+1) + ")"
			params = append(params, pq.Array(countries))
		}
		if platforms, ok := c.GetQueryArray("platform"); ok && len(platforms) > 0 {
			query += " AND platform = ANY($" + strconv.Itoa(len(params)+1) + ")"
			params = append(params, pq.Array(platforms))
		}

		// 添加年龄和性别条件，假设为单选
		if age, ok := c.GetQuery("age"); ok {
			query += " AND age = $" + strconv.Itoa(len(params)+1)
			params = append(params, age)
		}
		if gender, ok := c.GetQuery("gender"); ok {
			query += " AND gender = $" + strconv.Itoa(len(params)+1)
			params = append(params, gender)
		}

		// 添加排序和分页
		query += " ORDER BY endat ASC LIMIT $" + strconv.Itoa(len(params)+1) + " OFFSET $" + strconv.Itoa(len(params)+2)
		params = append(params, limit, offset)

		rows, err := db.Query(query, params...)
		if err != nil {
			log.Printf("Failed to list ads: %v", err) // 增加错误日志
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list ads"})
			return
		}

		defer rows.Close()

		var ads []models.Ad
		for rows.Next() {
			var ad models.Ad
			if err := rows.Scan(&ad.ID, &ad.Title, &ad.StartAt, &ad.EndAt); err != nil {
				log.Printf("Failed to scan ad: %v", err) // 增加错误日志
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan ad"})
				return
			}
			ads = append(ads, ad)
		}
		// 将广告列表序列化为 JSON
		adsJSON, err := json.Marshal(gin.H{"items": ads})
		if err != nil {
			log.Printf("Failed to serialize ads: %v", err)
			// 处理错误...
		}

		// 将序列化后的 JSON 字符串缓存到 Redis 中，设置适当的过期时间，例如 30 分钟
		err = rdb.Set(context.Background(), cacheKey, string(adsJSON), 30*time.Minute)
		if err != nil {
			log.Printf("Failed to cache ads list in Redis: %v", err)
			// 可以选择记录错误，但不一定要中断整个流程
		}

		// 返回查询结果
		c.JSON(http.StatusOK, gin.H{"items": ads}) // 统一响应结构
	}
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
