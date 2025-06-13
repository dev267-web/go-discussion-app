package tag

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	authmw "go-discussion-app/internal/auth"
	"go-discussion-app/models"
	"go-discussion-app/pkg/jwtutil"
)

// MockTagRepository is a mock implementation of tag.TagRepository
type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) GetAll(ctx context.Context) ([]models.Tag, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Tag), args.Error(1)
}

func (m *MockTagRepository) GetByName(ctx context.Context, name string) (*models.Tag, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *MockTagRepository) Create(ctx context.Context, name string) (int, error) {
	args := m.Called(ctx, name)
	return args.Int(0), args.Error(1)
}

// Helper to generate a JWT token for testing
func generateTestTokenTag(userID int) string {
	token, err := jwtutil.GenerateToken(userID)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate test token: %v", err))
	}
	return token
}

// Helper to set up the Gin router with TagController
// It creates a real TagService but with a mocked TagRepository.
func setupTagTestRouter(mockRepo TagRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create real service with mock repository
	tagService := NewService(mockRepo) // NewService is from tag package (service.go)
	tagController := NewController(tagService) // NewController is from tag package (controller.go)

	// As per internal/tag/routes.go, the /tags route is registered on a group.
	// The comment suggests this group should be protected.
	protectedGroup := router.Group("/")
	protectedGroup.Use(authmw.JWTAuthMiddleware())
	{
		protectedGroup.GET("/tags", tagController.ListHandler)
	}
	return router
}

func performTagRequest(r http.Handler, method, path, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil) // GET requests typically don't have a body
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// --- ListAllTags Tests (GET /tags) ---

func TestListTags_Success(t *testing.T) {
	mockRepo := new(MockTagRepository)
	router := setupTagTestRouter(mockRepo)
	token := generateTestTokenTag(1) // User ID 1

	expectedTags := []models.Tag{
		{ID: 1, Name: "Go", CreatedAt: time.Now()},
		{ID: 2, Name: "Testing", CreatedAt: time.Now()},
	}
	mockRepo.On("GetAll", mock.Anything).Return(expectedTags, nil)

	w := performTagRequest(router, "GET", "/tags", token)

	assert.Equal(t, http.StatusOK, w.Code)
	var tags []models.Tag
	err := json.Unmarshal(w.Body.Bytes(), &tags)
	assert.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, expectedTags[0].Name, tags[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestListTags_Success_NoTags(t *testing.T) {
	mockRepo := new(MockTagRepository)
	router := setupTagTestRouter(mockRepo)
	token := generateTestTokenTag(1)

	expectedTags := []models.Tag{} // Empty slice
	mockRepo.On("GetAll", mock.Anything).Return(expectedTags, nil)

	w := performTagRequest(router, "GET", "/tags", token)

	assert.Equal(t, http.StatusOK, w.Code)
	var tags []models.Tag
	err := json.Unmarshal(w.Body.Bytes(), &tags)
	assert.NoError(t, err)
	assert.Len(t, tags, 0)
	mockRepo.AssertExpectations(t)
}

func TestListTags_ServiceError(t *testing.T) {
	mockRepo := new(MockTagRepository)
	router := setupTagTestRouter(mockRepo)
	token := generateTestTokenTag(1)

	mockRepo.On("GetAll", mock.Anything).Return(nil, assert.AnError)

	w := performTagRequest(router, "GET", "/tags", token)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "server error", resp["error"])
	mockRepo.AssertExpectations(t)
}

func TestListTags_Unauthorized_NoToken(t *testing.T) {
	mockRepo := new(MockTagRepository)
	router := setupTagTestRouter(mockRepo)
	// No token provided

	w := performTagRequest(router, "GET", "/tags", "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	// Repository's GetAll should not be called
	mockRepo.AssertNotCalled(t, "GetAll", mock.Anything)
}

// Note: Tests for Create, GetByID/Name, Delete are not included as these functionalities
// are not present in the current TagController or TagService.
// Listing discussions by tag is handled by DiscussionController.
// Authorization tests beyond requiring a valid token for ListTags are not applicable
// as there are no modification endpoints.
