package routes

import (
	"gollet/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, walletHandler *handlers.WalletHandler, userHandler *handlers.UserHandler) {
	api := router.Group("/api")
	{
		api.POST("/users/signup", userHandler.Signup)
		api.POST("/users/login", userHandler.Login)

		wallet := api.Group("users/:userId/wallets")
		{
			wallet.POST("/deposit", walletHandler.Deposit)
			wallet.POST("/withdraw", walletHandler.Withdraw)
			// wallet.POST("/:userId/transfer", walletHandler.Transfer)
			// wallet.GET("/:userId/balance", walletHandler.GetBalance)
		}
	}
}
