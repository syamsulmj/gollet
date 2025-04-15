package repositories

import (
	"gollet/internal/models"

	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(userId uint) (*models.Wallet, error)
	FindByUserId(userId uint) (*models.Wallet, error)
	UpdateBalance(wallet *models.Wallet) error
}

type WalletRepositoryImpl struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &WalletRepositoryImpl{db: db}
}

func (r *WalletRepositoryImpl) Create(userId uint) (*models.Wallet, error) {
	wallet := models.Wallet{UserID: userId, Balance: 0, Currency: "USD", Metadata: "{}"}
	result := r.db.Create(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

func (r *WalletRepositoryImpl) FindByUserId(userId uint) (*models.Wallet, error) {
	var wallet models.Wallet
	result := r.db.Where("user_id = ?", userId).First(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

func (r *WalletRepositoryImpl) UpdateBalance(wallet *models.Wallet) error {
	result := r.db.Save(wallet)
	return result.Error
}
