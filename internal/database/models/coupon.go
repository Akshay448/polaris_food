package models

import (
	"gorm.io/gorm"
	"time"
)

type Coupon struct {
	gorm.Model
	Code          string
	Description   string
	DiscountType  string // "Percentage" or "Flat"
	DiscountValue float64
	ValidFrom     time.Time
	ValidUntil    time.Time
	MinOrderValue float64
	Active        bool
}
