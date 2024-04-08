package handlers

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

// AnyTime 是一个自定义的 sqlmock 参数匹配器，用来匹配任何 time.Time 类型的参数
type AnyTime struct{}

// Match 满足 sqlmock.Argument 接口，用于匹配任意的 time.Time 类型的参数
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// handlers/handlers_test.go

type MockRedisClient struct {
	// 这个 map 可以用来存储键值对，模拟 Redis 的行为
	Data map[string]string
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// 将值转换为字符串并存储在 Data map 中，模拟 Set 操作
	if strVal, ok := value.(string); ok {
		m.Data[key] = strVal
	}
	return nil
}

func TestCreateAd(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// 设置数据库操作的预期
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO Advertisements").
		WithArgs("Test Ad", AnyTime{}, AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// 初始化 gin 路由器
	r := gin.Default()

	// 初始化 MockRedisClient
	mockRedis := &MockRedisClient{Data: make(map[string]string)}

	// 使用 mock 数据库和 mock Redis 客户端
	// 注意：这里要确保 CreateAd 函数已经修改为接受 RedisClientInterface
	r.POST("/ad", CreateAd(db, mockRedis))

	// 创建测试请求
	ad := map[string]interface{}{
		"title":   "Test Ad",
		"startAt": "2022-01-01T00:00:00Z",
		"endAt":   "2022-01-10T00:00:00Z",
	}
	adJSON, _ := json.Marshal(ad)
	req, _ := http.NewRequest("POST", "/ad", bytes.NewBuffer(adJSON))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求并记录响应
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 校验响应状态码
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// 校验 Redis 操作
	if val, ok := mockRedis.Data["ad_title:Test Ad"]; !ok || val != "Test Ad" {
		t.Errorf("Expected 'ad_title:Test Ad' to be set in Redis with value 'Test Ad'")
	}

	// 校验数据库操作的预期是否都已满足
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
