package pkg

// Package pkg provides core files for food ordering system

// The models.go file defines the core data models for the food ordering system.
// It includes structures for products, orders, and other relevant entities.

import (
	"encoding/json"
	"time"
)

type Product struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	Price     float64 `gorm:"not null" json:"price"`
	Category  string `gorm:"not null" json:"category"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}

type OrderItem struct {
	ID        uint   `gorm:"primaryKey" json:"-"`
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func(o *OrderItem) TableName() string {
	return "order_items"
}

type Order struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Items     []OrderItem `gorm:"foreignKey:id" json:"items"`
	Products  []Product   `gorm:"many2many:product_list;" json:"products"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"-"`
}

type OrderReq struct {
	CouponCode string      `json:"couponCode"`
	Items      []OrderItem `json:"items"`
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
	// Coupon []Coupon `gorm:"many2many:coupon_sources;constraint:OnDelete:CASCADE"`     
}

func (CouponSource) TableName() string {
	return "coupon_source"
}