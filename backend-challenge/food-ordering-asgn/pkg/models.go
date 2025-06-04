package pkg

// Package pkg provides core files for food ordering system

// The models.go file defines the core data models for the food ordering system.
// It includes structures for products, orders, and other relevant entities.

import (
	"encoding/json"
	"time"
)

type Product struct {
	ID        uint `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Price     float64 `gorm:"not null"`
	Category  string `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Order struct {
	ID        uint `gorm:"primaryKey"`
	Items     []map[string]any `gorm:"type:jsonb"`
	Products  []Product `gorm:"many2many:order_products;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type OrderReq struct {
	Items []map[string]any `json:"items"`
	Products []Product `json:"products"`
}

type ApiResonse struct {
	StatusCode int
	Type string
	Message string
}

func (res *ApiResonse) Serialize() []byte {
	bytes, _ := json.Marshal(res)
	return bytes
}