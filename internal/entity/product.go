package entity

import "time"

type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name" validate:"required"`
	Price     float64   `gorm:"not null" json:"price" validate:"required"`
	Category  string    `gorm:"not null" json:"category" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
