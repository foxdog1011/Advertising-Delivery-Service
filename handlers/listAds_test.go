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
	// 如果键不存在，模拟 Redis 的行为，返回一个错误
	return "", errors.New("redis: nil") // Redis 在键不存在时返回的错误
}

func mockDatabase(t *testing.T) (*sql.DB, func()) {
	// 创建一个 SQLite 内存数据库
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}

	// 调用 setupDatabase 来执行数据库设置，比如创建表、插入测试数据等
	setupDatabase(db)

	// 返回数据库连接和一个清理函数
	return db, func() { db.Close() }
}

// setupDatabase 用于设置数据库环境，比如创建所需的表和插入初始数据
func setupDatabase(db *sql.DB) {
	// 创建表
	const createTableSQL = `CREATE TABLE IF NOT EXISTS Advertisements (
        adid INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        startat DATETIME NOT NULL,
        endat DATETIME NOT NULL
    );`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// 插入测试数据（根据需要调整）
	// ...
}
func TestListAds(t *testing.T) {
	// 创建一个模拟的数据库连接
	db, teardown := mockDatabase(t) // 需要实现该函数，建立一个模拟的数据库环境
	defer teardown()
	mockRedis := &MockRedisClient{Data: make(map[string]string)}
	// 创建一个新的 Gin 路由器
	router := gin.Default()

	// 注册你的路由和处理函数
	router.GET("/ads", ListAds(db, mockredis))

	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest("GET", "/ads", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// 记录响应
	w := httptest.NewRecorder()

	// 调用相应的处理函数
	router.ServeHTTP(w, req)

	// 检查响应码是否如预期
	if w.Code != http.StatusOK {
		t.Errorf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	// 可以进一步检查响应体是否符合预期
}
