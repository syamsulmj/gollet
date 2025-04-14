package services

import (
	"errors"
	"fmt"
	"gollet/internal/database"
	"gollet/internal/models"
	"gollet/internal/repositories"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepository   *repositories.UserRepository
	walletRepository *repositories.WalletRepository
}

func NewUserService(userRepository *repositories.UserRepository, walletRepository *repositories.WalletRepository) *UserService {
	return &UserService{
		userRepository:   userRepository,
		walletRepository: walletRepository,
	}
}

func (s *UserService) Signup(email, password string) (*models.User, error) {
	// Validate email and password
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}
	// Check if password is strong enough
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}
	// Check if email is valid
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("invalid email format")
	}
	// Check if user already exists
	existingUser, err := s.userRepository.FindByEmail(email)
	if existingUser != nil && err == nil {
		return nil, errors.New("user already exists")
	}

	var createdUser *models.User
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		newUser := &models.User{Email: email, Password: string(hashedPassword)}
		createdUser, err = s.userRepository.Create(newUser)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		wallet, err := s.walletRepository.Create(createdUser.ID)
		if err != nil {
			return fmt.Errorf("failed to create wallet: %w", err)
		}

		createdUser.Wallet = *wallet
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return the created user
	return createdUser, nil
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	// Validate email and password
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// Find the user by email
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
