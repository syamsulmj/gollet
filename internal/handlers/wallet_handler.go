package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"gollet/internal/services"
	"gollet/internal/utils"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

type DepositRequest struct {
	Amount   float64 `json:"amount" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
}

type WithdrawRequest struct {
	Amount   float64 `json:"amount" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
}

type TransferRequest struct {
	RecipientUserId uint    `json:"recipient_user_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
}

type TransactionHistoryResponse struct {
	ID              uint   `json:"id"`
	UserID          uint   `json:"user_id"`
	Amount          string `json:"amount"`
	TransactionType string `json:"transaction_type"`
	Metadata        string `json:"metadata"`
	Currency        string `json:"currency"`
	CreatedAt       string `json:"created_at"`
}

func (h *WalletHandler) Deposit(c *gin.Context) {
	// In Go: Gets URL parameter "userId" and assigns it to userIdStr
	// Equivalent Elixir: user_id_str = conn.params["user_id"]
	//
	// Go uses := for declaration + assignment in one line
	// Elixir uses = for assignment and doesn't need explicit declaration
	// Go uses camelCase while Elixir uses snake_case by convention
	userIdStr := c.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	amount := utils.MoneyToCents(req.Amount)
	wallet, err := h.walletService.Deposit(uint(userId), uint(amount))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deposit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id":       wallet.ID,
		"user_id":  wallet.UserID,
		"balance":  utils.CentsToMoney(int64(wallet.Balance), wallet.Currency),
		"currency": wallet.Currency,
	}})
}

func (h *WalletHandler) Withdraw(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	amount := utils.MoneyToCents(req.Amount)
	wallet, err := h.walletService.Withdraw(uint(userId), uint(amount))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to withdraw: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id":       wallet.ID,
		"user_id":  wallet.UserID,
		"balance":  utils.CentsToMoney(int64(wallet.Balance), wallet.Currency),
		"currency": wallet.Currency,
	}})
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	amount := utils.MoneyToCents(req.Amount)
	wallet, err := h.walletService.Transfer(uint(userId), req.RecipientUserId, uint(amount))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to transfer: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id":       wallet.ID,
		"user_id":  wallet.UserID,
		"balance":  utils.CentsToMoney(int64(wallet.Balance), wallet.Currency),
		"currency": wallet.Currency,
	}})
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	wallet, err := h.walletService.GetBalance(uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get balance: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id":       wallet.ID,
		"user_id":  wallet.UserID,
		"balance":  utils.CentsToMoney(int64(wallet.Balance), wallet.Currency),
		"currency": wallet.Currency,
	}})
}

func (h *WalletHandler) GetTransactionHistory(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	transactions, err := h.walletService.GetTransactionHistory(uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get transaction history: %v", err)})
		return
	}

	var response []TransactionHistoryResponse
	for _, transaction := range transactions {
		response = append(response, TransactionHistoryResponse{
			ID:              transaction.ID,
			UserID:          transaction.UserID,
			Amount:          utils.CentsToMoney(int64(transaction.Amount), transaction.Currency),
			TransactionType: transaction.TransactionType,
			Metadata:        transaction.Metadata,
			Currency:        transaction.Currency,
			CreatedAt:       transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
