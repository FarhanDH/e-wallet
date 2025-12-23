package handler

import (
	"ewallet/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo repository.UserRepository
}

// Constructor
func NewUserHandler(repo repository.UserRepository) UserHandler {
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

	err = h.repo.Transfer(req.FromID, req.ToID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer success"})
}
