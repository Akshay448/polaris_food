package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"polaris_food/internal/database/models"
	"polaris_food/internal/service/rider"
	"strconv"
)

// RegisterRider handles the registration of a new delivery rider.
// It expects rider details in the request body, validates them, and saves the rider to the database.
func RegisterRider(riderService rider.RiderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON payload to a Rider struct
		type RegisterRiderRequest struct {
			User         models.User         `json:"user"`
			RiderProfile models.RiderProfile `json:"riderProfile"`
		}
		var request RegisterRiderRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Validate the Rider data
		if err := validateRider(request.User, request.RiderProfile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save Rider to the database
		if err := riderService.RegisterRider(request.User, request.RiderProfile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		// Return response: success or error message
		c.JSON(http.StatusOK, gin.H{
			"message": "Rider registered successfully",
		})
	}
}

// FindNearestRider finds the nearest available rider to a restaurant for order pickup.
// It expects a restaurant ID and uses longitude and latitude data to determine minimum distance
func FindNearestRider(riderService rider.RiderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract restaurant ID from request
		restaurantID, err := strconv.ParseUint(c.Param("restaurantId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			return
		}

		// Query the database for nearest available rider
		nearestRider, err := riderService.GetNearestAvailableRider(uint(restaurantID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find nearest rider"})
			return
		}

		// Return the chosen rider's details
		c.JSON(http.StatusOK, nearestRider)
	}
}

// UpdateRiderLocation updates the current location of a rider.
// It expects rider ID and new location coordinates in the request body.
func UpdateRiderLocation(riderService rider.RiderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON payload to latitude and longitude
		var location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}

		if err := c.ShouldBindJSON(&location); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Extract rider ID from path
		riderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rider ID"})
			return
		}

		// Update rider's location in the database
		if err := riderService.UpdateRiderLocation(uint(riderID), location.Latitude, location.Longitude); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rider location"})
			return
		}

		// Return response: success or error message
		c.JSON(http.StatusOK, gin.H{"message": "Rider location updated successfully"})
	}
}

// ShowRiderOrderHistory shows the history of orders completed by a specific rider.
// It expects a rider ID and returns a list of orders.
func ShowRiderOrderHistory(riderService rider.RiderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract rider ID from request
		riderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rider ID"})
			return
		}

		// Query the database for rider's completed orders
		orderHistory, err := riderService.GetRiderOrders(uint(riderID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rider's order history"})
			return
		}

		// Return the list of orders
		c.JSON(http.StatusOK, orderHistory)
	}
}

// Add a function to validate Rider data
func validateRider(rider models.User, riderProfile models.RiderProfile) error {
	// Your validation logic goes here
	// For example, check if the required fields are present
	if rider.Username == "" || rider.Email == "" || rider.PasswordHash == "" {
		return errors.New("username, Email, and Password are required fields")
	}
	if riderProfile.Latitude == 0.0 || riderProfile.Longitude == 0.0 {
		return errors.New("latitude and longitude should be present")
	}
	return nil
}
