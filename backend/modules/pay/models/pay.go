package models

import (
	"github.com/goravel/framework/database/orm"
)

// Product is a purchasable item in the shop (point packages, memberships…).
// grant_points ties a purchase to the gamification currency.
type Product struct {
	orm.Timestamps
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"size:255" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	// PriceCents keeps money as the smallest currency unit (no floats).
	PriceCents  int64  `gorm:"default:0" json:"price_cents"`
	Currency    string `gorm:"size:8;default:USD" json:"currency"`
	Stock       int    `gorm:"default:0" json:"stock"`
	GrantPoints int    `gorm:"default:0" json:"grant_points"`
	Image       string `gorm:"size:512" json:"image"`
	IsActive    bool   `gorm:"default:true;index" json:"is_active"`
	Sort        int    `gorm:"default:100" json:"sort"`
	// Type: 'general' (plain purchase) or 'invite' (auto-delivers an invite
	// code to the buyer on payment).
	Type string `gorm:"size:16;default:general" json:"type"`
}

func (Product) TableName() string { return "products" }

// Order tracks a purchase attempt through a payment gateway. Status flow:
// pending → paid / expired / failed.
type Order struct {
	orm.Timestamps
	ID        uint64  `gorm:"primaryKey" json:"id"`
	OrderNo   string  `gorm:"uniqueIndex;size:64" json:"order_no"`
	UserID    uint64  `gorm:"index" json:"user_id"`
	ProductID uint64  `gorm:"index" json:"product_id"`
	Title     string  `gorm:"size:255" json:"title"`
	AmountCents int64 `gorm:"default:0" json:"amount_cents"`
	Currency  string  `gorm:"size:8;default:USD" json:"currency"`
	Gateway   string  `gorm:"size:32" json:"gateway"`
	GatewayRef string `gorm:"size:128" json:"gateway_ref"`
	Status    string  `gorm:"size:32;default:pending;index" json:"status"`
	PaidAt    *string `gorm:"" json:"paid_at"`
	// GrantedCode holds the delivered invite code for type='invite' goods.
	GrantedCode string `gorm:"size:64" json:"granted_code,omitempty"`
}

func (Order) TableName() string { return "orders" }
