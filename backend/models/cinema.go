package models

import (
	"time"
	"gorm.io/gorm"
)

type Chain struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Theatre struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Name             string         `json:"name"`
	ChainID          uint           `json:"chain_id"`
	Address          string         `json:"address"`
	SupportedFormats string         `json:"supported_formats"` // Comma separated: "2D,IMAX,4DX"
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type Show struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	MovieID   uint           `json:"movie_id"`
	TheatreID uint           `json:"theatre_id"`
	Format    string         `json:"format"`
	Showtime  time.Time      `json:"showtime"`
	Price     float64        `json:"price"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
