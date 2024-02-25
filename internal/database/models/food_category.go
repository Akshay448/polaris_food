package models

import (
	"gorm.io/gorm"
)

type FoodCategory struct {
	gorm.Model
	Name        string
	Description string
}
