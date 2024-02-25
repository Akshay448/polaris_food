package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"polaris_food/internal/database/models"
	"polaris_food/internal/service/restaurant"
	"strconv"
)

// RegisterRestaurant handles the registration of a new restaurant.
// It expects restaurant details and the menu items of the restaurant in the request body,
// validates them, and saves the restaurant to the database.
func RegisterRestaurant(restaurantService restaurant.RestaurantService) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Bind JSON payload to a request struct
		request := struct {
			Restaurant models.Restaurant `json:"restaurant"`
			MenuItems  []models.MenuItem `json:"menuItems"`
		}{}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Validate the Restaurant data
		if err := validateRestaurant(request.Restaurant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save Restaurant to the database
		if err := restaurantService.RegisterRestaurant(request.Restaurant, request.MenuItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
			return
		}

		// Return response: success or error message
		c.JSON(http.StatusOK, gin.H{
			"message": "Restaurant registered successfully",
		})
	}
}

// SuggestRestaurants suggests restaurants to a user based on the kind of food they want and the time within which they want the food.
// It receives 2 parameters - foodCategoryName and desired delivery time.
func SuggestRestaurants(restaurantService restaurant.RestaurantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract query parameters
		foodCategoryName := c.Query("food_category")
		deliveryTimeStr := c.Query("delivery_time")
		deliveryTime, err := strconv.Atoi(deliveryTimeStr) // Convert delivery time to int
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery time format"})
			return
		}

		// Query the database for matching restaurants based on category name and delivery time
		suggestedRestaurants, err := restaurantService.SuggestRestaurants(foodCategoryName, deliveryTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch suggested restaurants"})
			return
		}

		// Return a list of suggested restaurants
		c.JSON(http.StatusOK, suggestedRestaurants)
	}
}

// ProvideMenu returns the menu for a specific restaurant.
// It expects a restaurant ID as a path parameter.
func ProvideMenu(restaurantService restaurant.RestaurantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract restaurant ID from path
		restaurantID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			return
		}

		// Query the database for the restaurant's menu
		menu, err := restaurantService.GetMenu(uint(restaurantID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurant menu"})
			return
		}

		// Return the menu
		c.JSON(http.StatusOK, menu)
	}
}

// --------------- restaurant utils ------------------
// Add a function to validate Restaurant data
func validateRestaurant(restaurant models.Restaurant) error {
	// Your validation logic goes here
	// For example, check if the required fields are present
	if restaurant.Name == "" || (restaurant.Latitude == 0.0 && restaurant.Longitude == 0.0) {
		return errors.New("name and Location are required fields for a restaurant")
	}
	return nil
}
