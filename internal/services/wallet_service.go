package services

import (
	"errors"
	"fmt"
	"gollet/internal/database"
	"gollet/internal/models"
	"gollet/internal/repositories"

	"gorm.io/gorm"
)

type WalletService struct {
	walletRepository      *repositories.WalletRepository
	transactionRepository *repositories.TransactionRepository
}

func NewWalletService(walletRepository *repositories.WalletRepository, transactionRepository *repositories.TransactionRepository) *WalletService {
	return &WalletService{
		walletRepository:      walletRepository,
		transactionRepository: transactionRepository,
	}
}

func (s *WalletService) Deposit(userId uint, amount uint) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	wallet, err := s.walletRepository.FindByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find wallet: %w", err)
	}

	// Use DB transaction to ensure transaction and wallet update are created not independently
	// and rollback if any of them fails
	// This is important to ensure that the wallet balance and transaction are always in sync
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Update the wallet balance
		wallet.Balance += amount
		err = s.walletRepository.UpdateBalance(wallet)
		if err != nil {
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}

		_, err = s.transactionRepository.Create(&models.Transaction{
			UserId:          userId,
			TransactionType: "deposit",
			Amount:          amount,
			Currency:        wallet.Currency,
			Metadata:        `{"some_tracking": "12345_track"}`,
		})
		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletService) Withdraw(userId uint, amount uint) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("insufficient balance")
	}

	wallet, err := s.walletRepository.FindByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find wallet: %w", err)
	}

	if wallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// update the wallet balance
		wallet.Balance -= amount
		err = s.walletRepository.UpdateBalance(wallet)
		if err != nil {
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}

		_, err = s.transactionRepository.Create(&models.Transaction{
			UserId:          userId,
			TransactionType: "withdraw",
			Amount:          amount,
			Currency:        wallet.Currency,
			Metadata:        `{"some_tracking": "12345_track"}`,
		})
		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletService) GetBalance(userId uint) (*models.Wallet, error) {
	return s.walletRepository.FindByUserId(userId)
}

func (s *WalletService) GetTransactionHistory(userId uint) ([]models.Transaction, error) {
	return s.transactionRepository.FindByUserId(userId)
}
