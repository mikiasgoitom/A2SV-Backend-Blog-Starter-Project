package http_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
	handler "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/handler/http"
	contract "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase/contract"
)

// mockAvatarUseCase is a mock implementation of IUserAvatarUseCase for testing
type mockAvatarUseCase struct {
	createCalled      bool
	createUserID      string
	createFileHeader  *multipart.FileHeader
	createReturnMedia *entity.Media
	createReturnErr   error
}

func (m *mockAvatarUseCase) CreateUserAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (*entity.Media, error) {
	m.createCalled = true
	m.createUserID = userID
	m.createFileHeader = file
	return m.createReturnMedia, m.createReturnErr
}

func (m *mockAvatarUseCase) ReadUserAvatarMetadata(ctx context.Context, userID string) (*entity.Media, error) {
	return nil, nil
}
func (m *mockAvatarUseCase) UpdateUserAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (*entity.Media, error) {
	return nil, nil
}
func (m *mockAvatarUseCase) DeleteUserAvatar(ctx context.Context, userID string) error {
	return nil
}

// setupAvatarRouter creates a Gin engine with the avatar route
func setupAvatarRouter(usecase contract.IUserAvatarUseCase) *gin.Engine {
	r := gin.New()
	h := handler.NewUserAvatarHandler(usecase)
	r.POST("/users/:userID/avatar", h.CreateUserAvatar)
	return r
}

func TestCreateUserAvatar_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := &mockAvatarUseCase{
		createReturnMedia: &entity.Media{ID: "1", URL: "http://example.com/avatar.jpg"},
	}
	r := setupAvatarRouter(mockUC)

	// Prepare multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
	part.Write([]byte("fake image data"))
	writer.Close()

	req := httptest.NewRequest("POST", "/users/123/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.True(t, mockUC.createCalled)
	assert.Equal(t, "123", mockUC.createUserID)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "\"data\"")
}

func TestCreateUserAvatar_Fail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockErr := errors.New("upload error")
	mockUC := &mockAvatarUseCase{
		createReturnErr: mockErr,
	}
	r := setupAvatarRouter(mockUC)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
	part.Write([]byte("fake"))
	writer.Close()

	req := httptest.NewRequest("POST", "/users/123/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.True(t, mockUC.createCalled)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "upload error")
}

func TestCreateUserAvatar_BadRequest_NoFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := &mockAvatarUseCase{}
	r := setupAvatarRouter(mockUC)

	req := httptest.NewRequest("POST", "/users/123/avatar", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.False(t, mockUC.createCalled)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Avatar file is required")
}
