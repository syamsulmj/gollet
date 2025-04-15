package repositories

import (
	"gollet/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) (*models.Transaction, error)
	FindByUserId(userId uint) ([]models.Transaction, error)
}

// TransactionRepositoryImpl is a struct that implements the TransactionRepository interface
type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) Create(transaction *models.Transaction) (*models.Transaction, error) {
	result := r.db.Create(transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func (r *TransactionRepositoryImpl) FindByUserId(userId uint) ([]models.Transaction, error) {
	// []models.Transaction declares a slice (dynamic array) of Transaction structs
	// The [] indicates it's a slice that can hold multiple Transaction records
	// models.Transaction refers to the Transaction struct defined in the models package
	var transactions []models.Transaction
	result := r.db.Where("user_id = ?", userId).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}
