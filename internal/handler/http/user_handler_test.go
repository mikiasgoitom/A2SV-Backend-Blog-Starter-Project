package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	handler "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http"
	dto "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/dto"
	mocks "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http/mocks"
	validator "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/infrastructure/validator"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// enable Gin test mode and register custom validators for binding
	gin.SetMode(gin.TestMode)
	validator.RegisterCustomValidators()
	os.Exit(m.Run())
}

func setupRouter(h handler.UserHandlerInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", h.CreateUser)
	r.POST("/login", h.Login)
	r.GET("/users/:id", h.GetUser)
	return r
}

func TestCreateUser(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase()
	h := handler.NewUserHandler(mockUsecase)
	r := setupRouter(h)
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
	h := handler.NewUserHandler(mockUsecase)
	r := setupRouter(h)
	// Missing required fields to trigger validation error
	payload := dto.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
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
	h := handler.NewUserHandler(mockUsecase)
	r := setupRouter(h)
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
	h := handler.NewUserHandler(mockUsecase)
	r := setupRouter(h)
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
	h := handler.NewUserHandler(mockUsecase)
	r := setupRouter(h)
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
	h := handler.NewUserHandler(mockUsecase)
	r := setupRouter(h)
	id := uuid.New().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+id, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}
