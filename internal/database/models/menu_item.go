package models

import (
	"gorm.io/gorm"
)

type MenuItem struct {
	gorm.Model
	RestaurantID uint
	Name         string
	Description  string
	Price        float64
	CategoryID   uint
}
