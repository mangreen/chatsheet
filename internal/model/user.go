package model

import "time"

// User 模型用於應用程式使用者
type User struct {
	Email     string     `gorm:"primaryKey" json:"email"`
	Password  string     `gorm:"password" json:"-"`
	CreatedAt *time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt *time.Time `gorm:"default:now()" json:"updated_at"`
}
