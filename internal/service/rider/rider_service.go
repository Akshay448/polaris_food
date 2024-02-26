package rider

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math"
	"polaris_food/internal/database/models"
)

type RiderService interface {
	RegisterRider(user models.User, riderProfile models.RiderProfile) error
	UpdateRiderLocation(riderID uint, latitude, longitude float64) error
	GetRiderOrders(riderID uint) ([]models.Order, error)
	GetNearestAvailableRider(restaurantID uint) (models.RiderProfile, error)
}

type riderService struct {
	db *gorm.DB
}

func NewRiderService(db *gorm.DB) RiderService {
	return &riderService{db}
}

// RegisterRider registers a new rider in rider profile and the users table,
// assuming a rider is also a user - just an assumption for this design
func (s *riderService) RegisterRider(user models.User, riderProfile models.RiderProfile) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)

	// Start a transaction
	tx := s.db.Begin()

	// Create the user record
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the UserID in the RiderProfile
	riderProfile.UserID = user.ID

	// Create the rider profile record
	if err := tx.Create(&riderProfile).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// UpdateRiderLocation updates the current location of the rider
func (s *riderService) UpdateRiderLocation(riderID uint, latitude, longitude float64) error {
	return s.db.Model(&models.RiderProfile{}).Where("id = ?",
		riderID).Updates(map[string]interface{}{"latitude": latitude, "longitude": longitude}).Error
}

// GetRiderOrders retrieves the history of orders completed by the rider
func (s *riderService) GetRiderOrders(riderID uint) ([]models.Order, error) {
	var orders []models.Order
	// Directly fetch orders where the RiderID matches the given riderID
	err := s.db.Where("rider_id = ? AND status = ?", riderID, "delivered").
		Find(&orders).
		Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// GetNearestAvailableRider finds the nearest available rider to a restaurant.
func (s *riderService) GetNearestAvailableRider(restaurantID uint) (models.RiderProfile, error) {
	var nearestRider models.RiderProfile
	var availableRiders []models.RiderProfile
	var restaurant models.Restaurant

	// Retrieve the specified restaurant's details
	if err := s.db.First(&restaurant, restaurantID).Error; err != nil {
		return models.RiderProfile{}, err
	}

	// Retrieve all available riders
	// there should also be a filter where city = "city of the restaurant"
	// it will reduce the number of available riders to loop over
	if err := s.db.Where("availability_status = ? and is_delivering = ?", true,
		false).Find(&availableRiders).Error; err != nil {
		return models.RiderProfile{}, err
	}

	minDistance := math.MaxFloat64

	// Loop through available riders and find the nearest one
	for _, rider := range availableRiders {
		distance := Haversine(restaurant.Latitude, restaurant.Longitude, rider.Latitude, rider.Longitude)
		if distance < minDistance {
			minDistance = distance
			nearestRider = rider
		}
	}

	if nearestRider.ID == 0 {
		// No available riders found
		return models.RiderProfile{}, gorm.ErrRecordNotFound
	}

	return nearestRider, nil
}

// Haversine function calculates the distance between two points on Earth.
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	var R = float64(6371) // Earth's radius in kilometers
	var φ1 = lat1 * math.Pi / 180
	var φ2 = lat2 * math.Pi / 180
	var Δφ = (lat2 - lat1) * math.Pi / 180
	var Δλ = (lon2 - lon1) * math.Pi / 180

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	d := R * c // Distance in kilometers

	return d
}
