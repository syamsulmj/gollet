package repositories

import (
	"gollet/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) (*models.Transaction, error) {
	result := r.db.Create(transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func (r *TransactionRepository) FindByUserId(userId uint) ([]models.Transaction, error) {
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
