package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null"`
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
	Wallet   Wallet
	gorm.Model
}

type Wallet struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   uint   `gorm:"unique;not null"`
	Balance  uint   `gorm:"default:0"`  // Amount in cents
	Currency string `gorm:"not null"`   // Currency code (e.g., USD, EUR)
	Metadata string `gorm:"type:jsonb"` // Store extra info as JSON
	gorm.Model
}

type Transaction struct {
	ID              uint   `gorm:"primaryKey"`
	UserID          uint   `gorm:"not null"`
	TransactionType string `gorm:"not null"`   // deposit, withdraw, transfer_sent, transfer_received
	Amount          uint   `gorm:"not null"`   // Amount in cents
	Currency        string `gorm:"not null"`   // Currency code (e.g., USD, EUR)
	Metadata        string `gorm:"type:jsonb"` // Store extra info as JSON
	gorm.Model
}
