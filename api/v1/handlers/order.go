package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"polaris_food/internal/database/models"
	"polaris_food/internal/service/order"
)

// CreateOrder handles the creation of an order by a user.
// It expects order details in the request body.
func CreateOrder(orderService order.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON payload to an Order struct
		request := struct {
			Order      models.Order       `json:"order"`
			OrderItems []models.OrderItem `json:"orderItems"`
		}{}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the Order data
		if err := validateOrder(request.Order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save Order to the database (with status "created")
		request.Order.Status = "created"
		if err := orderService.CreateOrder(request.Order, request.OrderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}

		// Return response: success or error message
		c.JSON(http.StatusOK, gin.H{
			"message": "Order created successfully",
		})
	}
}

// UpdateOrder updates an order's status to "accepted".
func UpdateOrder(orderService order.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateRequest struct {
			OrderID uint   `json:"OrderID"`
			Status  string `json:"Status"` // The Status field will contain either "accepted" or "declined" or "delivered"
		}

		if err := c.ShouldBindJSON(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Validate the status to ensure it is either "accepted" or "declined" or "delivered"
		if updateRequest.Status != "accepted" && updateRequest.Status != "declined" && updateRequest.Status != "delivered" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
			return
		}

		if err := orderService.UpdateOrder(updateRequest.OrderID, updateRequest.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
			return
		}

		responseMessage := "Order status updated successfully"
		if updateRequest.Status == "declined" {
			responseMessage = "Order has been declined"
		} else if updateRequest.Status == "accepted" {
			responseMessage = "Order accepted successfully"
		} else if updateRequest.Status == "delivered" {
			responseMessage = "Order delivered successfully"
		}

		c.JSON(http.StatusOK, gin.H{
			"message": responseMessage,
		})
	}
}

// Add a function to validate Order data
func validateOrder(order models.Order) error {
	// here validation should also validate for the prices of each item
	// compare with restaurant's prices
	// in case a discount is provided calculate the difference from original price in menu items
	// table and validate the price
	if order.UserID == 0 || order.RestaurantID == 0 || order.TotalPrice <= 0 {
		return errors.New("UserID, RestaurantID, and TotalAmount are required fields for an order")
	}
	return nil
}
