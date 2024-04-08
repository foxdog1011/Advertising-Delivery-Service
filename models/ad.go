package models

import (
	"time"
)

// GenderType 定義了性别的枚举类型
type GenderType string

// 定義 GenderType 的可能值
const (
	GenderMale   GenderType = "M"
	GenderFemale GenderType = "F"
)

// Condition 定義了廣告出现的條件
type Condition struct {
	AgeStart int        `json:"ageStart,omitempty"`
	AgeEnd   int        `json:"ageEnd,omitempty"`
	Country  []string   `json:"country,omitempty"`  // 使用 ISO 3166-1 alpha-2 国家代碼
	Platform []string   `json:"platform,omitempty"` // 可能的值包括 "android", "ios", "web"
	Gender   GenderType `json:"gender,omitempty"`   // 使用 GenderType 類型
}

// Ad 定义了廣告本身的结構
type Ad struct {
	ID         int         `json:"id"`
	Title      string      `json:"title"`
	StartAt    time.Time   `json:"startAt"`
	EndAt      time.Time   `json:"endAt"`
	CreatedAt  time.Time   `json:"createdAt"`
	Conditions []Condition `json:"conditions,omitempty"`
}
