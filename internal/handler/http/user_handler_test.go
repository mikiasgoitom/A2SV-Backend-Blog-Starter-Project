package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/dto"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/mocks"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/validator"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)

	// Register custom validators
	validator.RegisterCustomValidators()

	// Run the tests
	os.Exit(m.Run())
}

func setupRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", handler.CreateUser)
	r.POST("/login", handler.Login)
	r.GET("/users/:id", handler.GetUser)
	return r
}

func TestCreateUser(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase()
	handler := NewUserHandler(mockUsecase)
	r := setupRouter(handler)

	payload := dto.CreateUserRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "Password123!",
		FirstName: "Test",
		LastName:  "User",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")
}

func TestCreateUser_Fail(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase()
	mockUsecase.ShouldFailCreateUser = true
	handler := NewUserHandler(mockUsecase)
	r := setupRouter(handler)

	// Missing required fields to trigger validation error
	payload := dto.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
		// FirstName and LastName omitted intentionally
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Field validation for 'FirstName' failed on the 'required' tag")
	assert.Contains(t, w.Body.String(), "Field validation for 'LastName' failed on the 'required' tag")
}

func TestLogin(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase()
	handler := NewUserHandler(mockUsecase)
	r := setupRouter(handler)

	payload := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "Password123!",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mock_access_token")
	assert.Contains(t, w.Body.String(), "mock_refresh_token")
}

func TestLogin_Fail(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase()
	mockUsecase.ShouldFailLogin = true
	handler := NewUserHandler(mockUsecase)
	r := setupRouter(handler)

	payload := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "Password123!",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")
}

func TestGetUser(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase()
	handler := NewUserHandler(mockUsecase)
	r := setupRouter(handler)

	id := uuid.New().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+id, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
}

func TestGetUser_Fail(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase()
	mockUsecase.ShouldFailGetByID = true
	handler := NewUserHandler(mockUsecase)
	r := setupRouter(handler)

	id := uuid.New().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+id, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}
