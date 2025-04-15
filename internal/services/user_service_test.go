package services

import (
	"gollet/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindById(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, nil
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) Create(userID uint) (*models.Wallet, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletRepository) FindByUserId(userID uint) (*models.Wallet, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletRepository) UpdateBalance(wallet *models.Wallet) error {
	args := m.Called(wallet)
	return args.Error(0)
}

func TestUserServiceSignupFailure(t *testing.T) {
	// Initialize the mock repositories
	mockUserRepo := new(MockUserRepository)
	mockWalletRepo := new(MockWalletRepository)
	userService := NewUserService(mockUserRepo, mockWalletRepo)

	// // Test case 1: Invalid Email
	_, err := userService.Signup("uwu-cat", "password123")
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid email format")

	// Test case 2: Password too short
	_, err = userService.Signup("uwu-cat@f.com", "pass")
	assert.Error(t, err)
	assert.EqualError(t, err, "password must be at least 6 characters long")

	// Test case 3: User already exists
	mockUserRepo.On("FindByEmail", "existing-cat@f.com").Return(&models.User{}, nil).Once()
	_, err = userService.Signup("existing-cat@f.com", "password123")
	assert.Error(t, err)
	assert.EqualError(t, err, "user already exists")
}
