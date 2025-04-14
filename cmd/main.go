package main

import (
	"fmt"
	"log"

	"gollet/internal/database"
	"gollet/internal/handlers"
	"gollet/internal/repositories"
	"gollet/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gollet/cmd/routes"
)

func main() {
	fmt.Println("Hello, World!")
	fmt.Println("This is a simple Go program.")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	database.ConnectDB()
	db := database.DB // Get the GORM database instance

	// Auto migrate the database schema
	database.MigrateDB()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, walletRepo)
	walletService := services.NewWalletService(walletRepo, transactionRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	walletHandler := handlers.NewWalletHandler(walletService)

	// Initialize the router
	router := gin.Default()
	routes.SetupRoutes(router, walletHandler, userHandler)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	fmt.Println("Server started on port 8080")
	fmt.Println("You can access the API at http://localhost:8080/api")
}
