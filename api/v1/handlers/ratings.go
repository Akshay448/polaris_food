package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"polaris_food/internal/database/models"
	"polaris_food/internal/service/ratings"
)

// SubmitRating gets the new rating submitted to a rider or to a user.
// it expects the request body to have ratedById, ratedToId, orderId, stars, optional comments
func SubmitRating(ratingsService ratings.RatingsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON payload to a Ratings struct
		var rating models.Rating
		if err := c.BindJSON(&rating); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
			return
		}

		// Call ratings service function defined with the ratings interface to save it to the database
		err := ratingsService.SubmitRatings(rating)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not submit rating"})
			return
		}

		// Return response: success submitted rating or error message
		c.JSON(http.StatusOK, gin.H{"message": "Rating submitted successfully"})
	}
}
