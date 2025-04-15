package repositories

import (
	"gollet/internal/models"

	"gorm.io/gorm"
)

// UserRepository is a struct that implements the UserRepository interface
type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindById(id uint) (*models.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(user *models.User) (*models.User, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Preload("Wallet").Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindById(id uint) (*models.User, error) {
	var user models.User
	result := r.db.Preload("Wallet").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
