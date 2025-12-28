package handler

import (
	"context"
	"errors"
	"ewallet/internal/entity"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define Object Mock
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUser(id int) (entity.User, error) {
	// Record called argument
	args := m.Called(id)
	// return fake data that prepared in test case
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(name string, balance int) error {
	args := m.Called(name, balance)
	return args.Error(0)
}

func (m *MockUserRepository) Transfer(ctx context.Context, fromID, toID, amount int) error {
	args := m.Called(ctx, fromID, toID, amount)
	return args.Error(0)
}

func (m *MockUserRepository) GetTransactionHistory(userID int) ([]entity.Transaction, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

// === UNIT TEST ===
func TestGetUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userHandler := NewUserHandler(mockRepo)

	// expected output from DB
	expectedUser := entity.User{
		ID:      1,
		Name:    "Farhan Test",
		Balance: 100000,
	}

	// Expectation scenario
	mockRepo.On("GetUser", 1).Return(expectedUser, nil)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Pretend there is a request to /users/1
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	userHandler.GetUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Farhan Test")

	mockRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userHandler := NewUserHandler(mockRepo)

	// Error Scenario
	mockRepo.On("GetUser", 99).Return(entity.User{}, errors.New("User not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "99"}}

	userHandler.GetUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userHandler := NewUserHandler(mockRepo)

	nameTest := "Joko"
	balanceTest := 100000
	mockRepo.On("CreateUser", nameTest, balanceTest).Return(nil)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBody := `{"name": "Joko", "balance": 100000}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c.Request = req

	userHandler.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")

	mockRepo.AssertExpectations(t)
}
