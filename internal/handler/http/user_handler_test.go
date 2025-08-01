package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/dto"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/mocks"
	"github.com/stretchr/testify/assert"
)

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
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
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

	payload := dto.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "user creation failed")
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
	assert.Contains(t, w.Body.String(), "login failed")
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
	assert.Contains(t, w.Body.String(), "user not found")
}
