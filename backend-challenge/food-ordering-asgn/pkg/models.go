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

type OrderItem struct {
	ProductID uint `gorm:"for"`
	Quantity  int  `json:"quantity"`
}

type OrderReq struct {
	Items    []OrderItem `json:"items"`
	Products []Product   `json:"products"`
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

// Each coupon has a unique code and can be associated with multiple source files
type Coupon struct {
	ID         uint          `gorm:"primaryKey"`                    
	Code       string        `gorm:"unique;not null"`              
	SourceFile []CouponSource `gorm:"many2many:coupon_sources;constraint:OnDelete:CASCADE"`    
}

// A source can contain multiple coupons, and coupons can come from multiple sources
type CouponSource struct {
	ID     uint     `gorm:"primaryKey;autoIncrement"`                    
	Source string   `gorm:"not null;unique"`                      
	Coupon []Coupon `gorm:"many2many:coupon_sources;constraint:OnDelete:CASCADE"`     
}

func (CouponSource) TableName() string {
	return "coupon_source"
}