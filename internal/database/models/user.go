package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string
	Email         string `gorm:"unique"`
	PasswordHash  string
	Role          string   // e.g., "customer", "rider"
	AverageRating *float64 // * represents can be null
}
