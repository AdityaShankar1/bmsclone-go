package models

import (
	"time"
	"gorm.io/gorm"
)

type Movie struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Genre     string         `json:"genre"`
	Rating    string         `json:"rating"` // A, U, UA
	PosterURL string         `json:"poster_url"`
	Formats   string         `json:"formats"` // Comma separated: "2D,IMAX"
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Ticket struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	ShowID              uint           `json:"show_id"`
	UserID              uint           `json:"user_id"`
	SeatID              string         `json:"seat_id"`
	PricePaid           float64        `json:"price_paid"`
	Status              string         `json:"status"` // "Confirmed", "Cancelled"
	UserRating          int            `json:"user_rating"` // Added mapping for stars
	IsPaymentSuccessful bool           `json:"is_payment_successful"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserAccount struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	BMSCash           float64        `json:"bms_cash"`
	InvestmentBalance float64        `json:"investment_balance"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}
