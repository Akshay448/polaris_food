package main

import (
	"github.com/gin-gonic/gin"
	"polaris_food/api/v1/handlers"
	"polaris_food/internal/service/order"
	"polaris_food/internal/service/ratings"
	"polaris_food/internal/service/restaurant"
	"polaris_food/internal/service/rider"
	"polaris_food/internal/service/user"
)

const CONFIGPATH = "./config"

func main() {
	router := gin.Default()

	// init viper for configuration
	initViperConfig()

	// init db - using sqlite for now
	db := initDb()

	// init services
	userService := user.NewUserService(db)
	orderService := order.NewOrderService(db)
	riderService := rider.NewRiderService(db)
	restaurantService := restaurant.NewRestaurantService(db)
	ratingsService := ratings.NewRatingService(db)

	v1 := router.Group("/api/v1")
	{
		// user handlers
		v1.POST("/register/user", handlers.RegisterUser(userService))
		v1.GET("/users/:id/orders", handlers.ShowUserOrderHistory(userService))
		v1.GET("/users/:id/coupons", handlers.GetCoupons(userService))

		// rider handlers
		v1.POST("/register/rider", handlers.RegisterRider(riderService))
		v1.GET("/riders/nearest/:restaurantId", handlers.FindNearestRider(riderService))
		v1.PUT("/riders/:id/location", handlers.UpdateRiderLocation(riderService))
		v1.GET("/riders/:id/orders", handlers.ShowRiderOrderHistory(riderService))

		// restaurant handlers
		v1.POST("/register/restaurant", handlers.RegisterRestaurant(restaurantService))
		v1.GET("/restaurants/suggest", handlers.SuggestRestaurants(restaurantService))
		v1.GET("/restaurants/:id/menu", handlers.ProvideMenu(restaurantService))

		// orders handlers
		v1.POST("/orders/create", handlers.CreateOrder(orderService))
		v1.POST("/orders/update", handlers.UpdateOrder(orderService))

		// ratings handlers
		v1.POST("/ratings/submit", handlers.SubmitRating(ratingsService))
	}

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
