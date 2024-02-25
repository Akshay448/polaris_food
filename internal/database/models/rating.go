package models

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	OrderID   uint
	RatedByID uint    // UserID of -  who gives the rating
	RatedToID uint    // UserID to whom rating is given
	Stars     float64 // 1 to 5 stars scale
	Comment   string  // Optional
}
