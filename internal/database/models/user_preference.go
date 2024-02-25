package models

import (
	"gorm.io/gorm"
)

type UserPreference struct {
	gorm.Model
	UserID         uint
	FoodCategoryID uint
}
