package models

import (
	"gorm.io/gorm"
	"time"
)

type OrderStatusHistory struct {
	gorm.Model
	OrderID       uint
	RiderID       uint   // This ensures we can track which rider was involved in each status update.
	CurrentStatus string // new status
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
