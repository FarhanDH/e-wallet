package handler

import (
	"context"
	"ewallet/internal/entity"
	"ewallet/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo repository.UserRepositoryInterface
}

// Constructor
func NewUserHandler(repo repository.UserRepositoryInterface) UserHandler {
	return UserHandler{repo: repo}
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	var req struct {
		Name    string `json:"name"`    // User's name
		Balance int    `json:"balance"` // Initial balance for the user
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})

		return
	}

	err = h.repo.CreateUser(req.Name, req.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	// Get ID from url parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // Convert string to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := h.repo.GetUser(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Transfer(c *gin.Context) {
	var req struct {
		FromID int `json:"from_id"`
		ToID   int `json:"to_id"`
		Amount int `json:"amount"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Create new context that has 2 seconds timeout.
	// The parent is c.Request.Context() context from Gin
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	err = h.repo.Transfer(ctx, req.FromID, req.ToID, req.Amount)
	if err != nil {
		if err == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Transaction timeout, please try again."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer success"})
}

func (h *UserHandler) GetHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // Convert string to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	history, err := h.repo.GetTransactionHistory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if history == nil {
		history = []entity.Transaction{}
	}

	c.JSON(http.StatusOK, history)
}
