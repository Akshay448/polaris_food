package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"polaris_food/internal/database/models"
	"polaris_food/internal/service/user"
	"strconv"
)

// RegisterUser handles the registration of a new user.
// It expects user details in the request body, validates them, and saves the user to the database.
func RegisterUser(userService user.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON payload to a User struct
		var newUser models.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		// Validate the User data
		if err := validateUser(newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := userService.RegisterUser(newUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Return response: success or error message
		c.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
		})
	}
}

// ShowUserOrderHistory shows the history of orders placed by a specific user.
// It expects a user ID and returns a list of orders.
func ShowUserOrderHistory(userService user.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		orders, err := userService.GetUserOrders(uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch order history"})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}

// GetCoupons returns all the coupons for a given user.
// It expects a user ID and returns a list of coupons.
func GetCoupons(userService user.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		orders, err := userService.GetUserCoupons(uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch order history"})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}

func validateUser(user models.User) error {
	// check if the required fields are present
	if user.Username == "" || user.Email == "" || user.PasswordHash == "" {
		return errors.New("username, Email, and Password are required fields")
	}
	return nil
}
