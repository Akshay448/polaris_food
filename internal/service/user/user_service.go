package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"polaris_food/internal/database/models"
)

type UserService interface {
	RegisterUser(user models.User) error
	GetUserOrders(userID uint) ([]models.Order, error)
	GetUserCoupons(userID uint) ([]models.Coupon, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db}
}

// RegisterUser registers a new user to the table user
func (s *userService) RegisterUser(user models.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)

	// Save to database
	if err := s.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserOrders retrieves the history of orders placed by a user
func (s *userService) GetUserOrders(userID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := s.db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// GetUserCoupons gets all the active coupons for this user
func (s *userService) GetUserCoupons(userID uint) ([]models.Coupon, error) {
	var coupons []models.Coupon

	// Get user coupons from the UserCoupon table
	var userCoupons []models.UserCoupon
	if err := s.db.Where("user_id = ?", userID).Find(&userCoupons).Error; err != nil {
		return nil, err
	}

	// Iterate through userCoupons to get details from the Coupon table
	for _, userCoupon := range userCoupons {
		var coupon models.Coupon
		if err := s.db.Where("id = ? AND active = ?", userCoupon.CouponID, true).First(&coupon).Error; err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}

	return coupons, nil
}
