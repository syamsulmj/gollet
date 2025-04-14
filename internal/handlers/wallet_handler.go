package handlers

import (
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
		"user_id":  wallet.UserId,
		"balance":  utils.CentsToMoney(int64(wallet.Balance), wallet.Currency),
		"currency": wallet.Currency,
	}})
}
