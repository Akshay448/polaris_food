package ratings

import (
	"gorm.io/gorm"
	"polaris_food/internal/database/models"
)

type RatingsService interface {
	SubmitRatings(rating models.Rating) error
}

type ratingsService struct {
	db *gorm.DB
}

func NewRatingService(db *gorm.DB) RatingsService {
	return &ratingsService{db}
}

// SubmitRatings saves the rating from a user(customer or rider) to another user (customer or rider)
func (s *ratingsService) SubmitRatings(rating models.Rating) error {
	// Save the ratings data to the ratings table
	if err := s.db.Create(&rating).Error; err != nil {
		return err
	}

	// Fetch all the ratings done to this user
	var allRatings []models.Rating
	if err := s.db.Where("rated_to_id = ?", rating.RatedToID).Find(&allRatings).Error; err != nil {
		return err
	}

	// Recalculate the average rating for this user
	var totalStars float64
	for _, r := range allRatings {
		totalStars += r.Stars
	}

	var averageRating float64
	if len(allRatings) > 0 {
		averageRating = totalStars / float64(len(allRatings))
	}

	// Save the new average rating to the users table in the column average_rating
	if err := s.db.Model(&models.User{}).Where("id = ?", rating.RatedToID).Update("average_rating", averageRating).Error; err != nil {
		return err
	}

	return nil
}
