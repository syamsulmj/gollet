package handlers

import (
	"fmt"
	"net/http"

	"gollet/internal/services"
	"gollet/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	WalletID uint   `json:"wallet_id"`
	Balance  string `json:"balance"`
}

func (h *UserHandler) Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	user, err := h.userService.Signup(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create user: %v", err)})
		return
	}

	if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	var response = UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		WalletID: user.Wallet.ID,
		Balance:  utils.CentsToMoney(int64(user.Wallet.Balance), user.Wallet.Currency),
	}
	c.JSON(http.StatusCreated, gin.H{
		"data":    response,
		"message": "User created successfully",
		"status":  "success",
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	var response = UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		WalletID: user.Wallet.ID,
		Balance:  utils.CentsToMoney(int64(user.Wallet.Balance), user.Wallet.Currency),
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "User logged in successfully",
	})
}
