package utils

import (
	"fmt"
	"math"
)

// CentsToMoney converts an integer amount of cents to a monetary value string
// based on the provided currency. Currently supports USD, EUR, GBP, and JPY.
func CentsToMoney(cents int64, currency string) string {
	switch currency {
	case "USD":
		return fmt.Sprintf("$%.2f", float64(cents)/100)
	case "MYR":
		return fmt.Sprintf("€%.2f", float64(cents)/100)
	default:
		// For unknown currencies, just show the decimal amount
		return fmt.Sprintf("%.2f", float64(cents)/100)
	}
}

// MoneyToCents converts a float monetary value to cents
func MoneyToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}
