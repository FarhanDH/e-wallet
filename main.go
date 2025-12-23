package main

import (
	"ewallet/config"
	"ewallet/internal/handler"
	"ewallet/internal/repository"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// setup router
	r := gin.Default()

	// Define Routes
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/:id", userHandler.GetUser)
	r.POST("/transfer", userHandler.Transfer)
	r.GET("/users/:id/transactions", userHandler.GetHistory)

	fmt.Println("Server running on port 8080")
	r.Run(":8080")
}
