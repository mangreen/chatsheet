package model

import (
	"time"

	uuid "github.com/google/uuid"
)

// UnipileAccount 模型用於儲存連結的第三方帳號
type UnipileAccount struct {
	ID        uuid.UUID  `gorm:"primaryKey;default:gen_random_uuid();not null" json:"id"`
	UserEmail string     `gorm:"not null" json:"user_email"`
	Provider  string     `gorm:"not null" json:"provider"`          // 例如 "linkedin"
	AccountID string     `gorm:"unique;not null" json:"account_id"` // Unipile 返回的 account_id
	CreatedAt *time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt *time.Time `gorm:"default:now()" json:"updated_at"`
}
