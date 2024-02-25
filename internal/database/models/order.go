package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID          uint
	RestaurantID    uint
	RiderID         uint
	Status          string // e.g., "created", "accepted or declined", "delivered"
	TotalPrice      float64
	DeliveryAddress string
	CouponId        *uint // can be null
}
