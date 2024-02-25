package models

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	OrderID    uint
	MenuItemID uint
	Quantity   int
	Price      float64 // Price at the time of order
}
