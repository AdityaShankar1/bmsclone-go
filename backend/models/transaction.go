package models

import (
	"time"
	"gorm.io/gorm"
)

type Transaction struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `json:"user_id"`
	Amount      float64        `json:"amount"`
	Type        string         `json:"type"` // "Investment", "Refund", "Payment"
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
