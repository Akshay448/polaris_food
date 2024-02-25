package restaurant

import (
	"errors"
	"gorm.io/gorm"
	"polaris_food/internal/database/models"
)

type RestaurantService interface {
	RegisterRestaurant(restaurant models.Restaurant, menuItems []models.MenuItem) error
	GetMenu(restaurantID uint) ([]models.MenuItem, error)
	SuggestRestaurants(foodCategoryName string, deliveryTime int) ([]models.Restaurant, error)
}

type restaurantService struct {
	db *gorm.DB
}

func NewRestaurantService(db *gorm.DB) RestaurantService {
	return &restaurantService{db}
}

// RegisterRestaurant saves new restaurant details to the database along with the menu items
func (s *restaurantService) RegisterRestaurant(restaurant models.Restaurant, menuItems []models.MenuItem) error {
	// Start a transaction
	tx := s.db.Begin()

	// Create the restaurant record
	if err := tx.Create(&restaurant).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the RestaurantID for each menu item and save them to the database
	for i := range menuItems {
		// additional logic to check for new food category, skipping it for now
		//categoryId := createUpdateCategory(tx)
		//menuItems[i].CategoryID = categoryId

		menuItems[i].RestaurantID = restaurant.ID
		if err := tx.Create(&menuItems[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

// GetMenu returns all the menu items from a given restaurant id
func (s *restaurantService) GetMenu(restaurantID uint) ([]models.MenuItem, error) {
	var menuItems []models.MenuItem
	err := s.db.Where("restaurant_id = ?", restaurantID).Find(&menuItems).Error
	if err != nil {
		return nil, err
	}

	return menuItems, nil
}

// SuggestRestaurants returns all the nearby restaurants,
// given food category within desired delivery time
func (s *restaurantService) SuggestRestaurants(foodCategoryName string, deliveryTime int) ([]models.Restaurant, error) {

	// first extract the category_id using foodCategoryName from table food_categories
	// chinese      1
	// italian      2
	// north indian 3
	var foodCategory models.FoodCategory
	var suggestedRestaurants []models.Restaurant
	err := s.db.Table("food_categories").Where("name = ?", foodCategoryName).Select("id").First(&foodCategory).Error
	// If there's an error or the category is not found, fetch all open restaurants with the desired delivery time
	if err != nil {
		errAll := s.db.
			Where("is_open = TRUE AND delivery_time <= ?", deliveryTime).
			Find(&suggestedRestaurants).Error

		if errAll != nil {
			return nil, errAll // Return an error if there's an issue fetching all restaurants
		}

		return suggestedRestaurants, nil // Return all restaurants meeting the open and delivery time criteria
	}
	foodCategoryID := foodCategory.ID

	// Now that we have the category ID, we can query for restaurants and menu items
	// joining tables restaurants with menu_items
	// then filtering the joined table with food_categoryId, open restaurants and desired delivery time
	err = s.db.
		Joins("JOIN menu_items ON restaurants.id = menu_items.restaurant_id").
		Where("menu_items.category_id = ? AND restaurants.is_open = TRUE AND restaurants.delivery_time <= ?", foodCategoryID, deliveryTime).
		Group("restaurants.id").
		Find(&suggestedRestaurants).Error

	if err != nil {
		return nil, err
	}

	// Location-based filtering would be done here - skipping it for now
	// filterRestaurantsByNearestToUserLocation(suggestedRestaurants)

	return suggestedRestaurants, nil
}

// createUpdateCategory checks if a category exists and returns its ID, or creates a new one if it does not exist.
func createUpdateCategory(tx *gorm.DB, categoryName string) (uint, error) {
	var category models.FoodCategory

	// Check if the category already exists
	if err := tx.Where("name = ?", categoryName).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Category does not exist, create a new one
			newCategory := models.FoodCategory{
				Name: categoryName,
				// Optionally set the description or any other fields here
			}
			if err := tx.Create(&newCategory).Error; err != nil {
				// Return an error if the category cannot be created
				return 0, err
			}
			// Return the new category's ID
			return newCategory.ID, nil
		} else {
			// Return any other error encountered
			return 0, err
		}
	}

	// Return the existing category's ID if found
	return category.ID, nil
}
