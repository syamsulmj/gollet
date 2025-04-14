package repositories

import (
	"gollet/internal/models"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(userId uint) (*models.Wallet, error) {
	wallet := models.Wallet{UserId: userId, Balance: 0, Currency: "USD", Metadata: "{}"}
	result := r.db.Create(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

func (r *WalletRepository) FindByUserId(userId uint) (*models.Wallet, error) {
	var wallet models.Wallet
	result := r.db.Where("user_id = ?", userId).First(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

func (r *WalletRepository) UpdateBalance(wallet *models.Wallet) error {
	result := r.db.Save(wallet)
	return result.Error
}
