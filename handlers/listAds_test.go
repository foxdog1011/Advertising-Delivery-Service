package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	// 检查键是否存在
	if val, ok := m.Data[key]; ok {
		// 如果键存在，返回值
		return val, nil
	}
	// 如果键不存在，模擬 Redis 的行为，返回一个错误
	return "", errors.New("redis: nil") // Redis 在键不存在时返回的错误
}

func mockDatabase(t *testing.T) (*sql.DB, func()) {
	// 創建一个 SQLite 内存資料库
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}

	// 调用 setupDatabase 来執行資料库設置，比如創建表、插入测试資料等
	setupDatabase(db)

	// 返回資料库連接和一個清理函数
	return db, func() { db.Close() }
}

// setupDatabase 用于設置資料库環境，比如創建所需的表和插入初始資料
func setupDatabase(db *sql.DB) {
	// 創建表
	const createTableSQL = `CREATE TABLE IF NOT EXISTS Advertisements (
        adid INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        startat DATETIME NOT NULL,
        endat DATETIME NOT NULL
    );`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

}
func TestListAds(t *testing.T) {
	// 創建一个模擬的資料库連接
	db, teardown := mockDatabase(t) // 需要實現该函数，建立一个模擬的資料库環境
	defer teardown()
	mockRedis := &MockRedisClient{Data: make(map[string]string)}
	// 創建一个新的 Gin 路由器
	router := gin.Default()

	// 註冊路由和處理函数
	router.GET("/ads", ListAds(db, mockRedis))

	// 創建一个新的 HTTP 请求
	req, err := http.NewRequest("GET", "/ads", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// 紀錄
	w := httptest.NewRecorder()

	// 调用相應的處理函数
	router.ServeHTTP(w, req)

	// 检查響應码是否如预期
	if w.Code != http.StatusOK {
		t.Errorf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

}
