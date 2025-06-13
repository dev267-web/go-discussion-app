package comment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	authmw "go-discussion-app/internal/auth" // Renamed to avoid conflict
	"go-discussion-app/models"
	"go-discussion-app/pkg/jwtutil"
)

// MockCommentService is a mock implementation of comment.Service
type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) AddComment(ctx context.Context, discussionID, userID int, content string) (int, error) {
	args := m.Called(ctx, discussionID, userID, content)
	return args.Int(0), args.Error(1)
}

func (m *MockCommentService) GetComments(ctx context.Context, discussionID int) ([]models.Comment, error) {
	args := m.Called(ctx, discussionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Comment), args.Error(1)
}

// Helper to generate a JWT token for testing
func generateTestTokenComment(userID int) string {
	token, err := jwtutil.GenerateToken(userID)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate test token: %v", err))
	}
	return token
}

// Helper to set up the Gin router with CommentController and middleware
func setupCommentTestRouter(mockService Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	commentController := NewController(mockService)

	// Apply JWT middleware to the group where comment routes are defined
	// This assumes that these routes are intended to be protected.
	authedRoutes := router.Group("/") // Or specific group like /api/v1
	authedRoutes.Use(authmw.JWTAuthMiddleware())
	{
		// Path: /discussions/:id/comments
		// The :id here is discussionID
		authedRoutes.POST("/discussions/:id/comments", commentController.Create)
		authedRoutes.GET("/discussions/:id/comments", commentController.List)
	}
	return router
}

func performCommentRequest(r http.Handler, method, path, token string, body interface{}) *httptest.ResponseRecorder {
	var reqBodyReader *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBodyReader = bytes.NewBuffer(jsonData)
	} else {
		reqBodyReader = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(method, path, reqBodyReader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// --- Create Comment Tests (POST /discussions/:discussionID/comments) ---

func TestCreateComment_Success(t *testing.T) {
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService)

	actingUserID := 1
	discussionID := 10
	token := generateTestTokenComment(actingUserID)
	dto := CreateCommentDTO{Content: "This is a test comment."}
	expectedCommentID := 123

	mockService.On("AddComment", mock.Anything, discussionID, actingUserID, dto.Content).Return(expectedCommentID, nil)

	w := performCommentRequest(router, "POST", fmt.Sprintf("/discussions/%d/comments", discussionID), token, dto)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedCommentID, resp["id"])
	mockService.AssertExpectations(t)
}

func TestCreateComment_Unauthorized_NoToken(t *testing.T) {
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService)
	discussionID := 10
	dto := CreateCommentDTO{Content: "Test comment"}

	w := performCommentRequest(router, "POST", fmt.Sprintf("/discussions/%d/comments", discussionID), "", dto) // No token

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	// Service should not be called
	mockService.AssertNotCalled(t, "AddComment", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateComment_InvalidDiscussionID_Format(t *testing.T) {
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService)
	actingUserID := 1
	token := generateTestTokenComment(actingUserID)
	dto := CreateCommentDTO{Content: "Test comment"}

	w := performCommentRequest(router, "POST", "/discussions/invalidID/comments", token, dto)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid discussion ID", resp["error"])
}

func TestCreateComment_InvalidPayload_BindingError(t *testing.T) {
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService)
	actingUserID := 1
	discussionID := 10
	token := generateTestTokenComment(actingUserID)

	w := performCommentRequest(router, "POST", fmt.Sprintf("/discussions/%d/comments", discussionID), token, "not-json")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid payload", resp["error"])
}

func TestCreateComment_InvalidPayload_ValidationError(t *testing.T) {
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService)
	actingUserID := 1
	discussionID := 10
	token := generateTestTokenComment(actingUserID)
	dto := CreateCommentDTO{Content: ""} // Empty content, fails validation

	w := performCommentRequest(router, "POST", fmt.Sprintf("/discussions/%d/comments", discussionID), token, dto)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "content is required", resp["error"]) // Error from dto.Validate()
}

func TestCreateComment_ServiceError(t *testing.T) {
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService)
	actingUserID := 1
	discussionID := 10
	token := generateTestTokenComment(actingUserID)
	dto := CreateCommentDTO{Content: "Valid comment"}

	mockService.On("AddComment", mock.Anything, discussionID, actingUserID, dto.Content).Return(0, assert.AnError)

	w := performCommentRequest(router, "POST", fmt.Sprintf("/discussions/%d/comments", discussionID), token, dto)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "could not add comment", resp["error"])
	mockService.AssertExpectations(t)
}

// --- List Comments for Discussion Tests (GET /discussions/:discussionID/comments) ---

func TestListComments_Success(t *testing.T) {
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService)
	discussionID := 10
	token := generateTestTokenComment(1) // Token needed if route is protected by JWTAuthMiddleware

	expectedComments := []models.Comment{
		{ID: 1, DiscussionID: discussionID, UserID: 1, Content: "Comment 1"},
		{ID: 2, DiscussionID: discussionID, UserID: 2, Content: "Comment 2"},
	}
	mockService.On("GetComments", mock.Anything, discussionID).Return(expectedComments, nil)

	w := performCommentRequest(router, "GET", fmt.Sprintf("/discussions/%d/comments", discussionID), token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var comments []models.Comment
	err := json.Unmarshal(w.Body.Bytes(), &comments)
	assert.NoError(t, err)
	assert.Len(t, comments, 2)
	assert.Equal(t, expectedComments[0].Content, comments[0].Content)
	mockService.AssertExpectations(t)
}

func TestListComments_Success_NoComments(t *testing.T) {
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService)
	discussionID := 10
	token := generateTestTokenComment(1)
	expectedComments := []models.Comment{} // Empty slice

	mockService.On("GetComments", mock.Anything, discussionID).Return(expectedComments, nil)

	w := performCommentRequest(router, "GET", fmt.Sprintf("/discussions/%d/comments", discussionID), token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var comments []models.Comment
	err := json.Unmarshal(w.Body.Bytes(), &comments)
	assert.NoError(t, err)
	assert.Len(t, comments, 0)
	mockService.AssertExpectations(t)
}

func TestListComments_InvalidDiscussionID_Format(t *testing.T) {
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService)
	token := generateTestTokenComment(1)

	w := performCommentRequest(router, "GET", "/discussions/invalidID/comments", token, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid discussion ID", resp["error"])
}

func TestListComments_ServiceError(t *testing.T) {
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService)
	discussionID := 10
	token := generateTestTokenComment(1)

	mockService.On("GetComments", mock.Anything, discussionID).Return(nil, assert.AnError)

	w := performCommentRequest(router, "GET", fmt.Sprintf("/discussions/%d/comments", discussionID), token, nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "could not fetch comments", resp["error"])
	mockService.AssertExpectations(t)
}

func TestListComments_Unauthorized_NoToken(t *testing.T) {
	// This test assumes the GET /discussions/:id/comments route is protected by JWTAuthMiddleware.
	// If it's a public route, this test would fail (expect 200) or need adjustment.
	// Based on setupCommentTestRouter, it IS protected.
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService)
	discussionID := 10

	w := performCommentRequest(router, "GET", fmt.Sprintf("/discussions/%d/comments", discussionID), "", nil) // No token

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Note: Tests for Update and Delete are not included as these functionalities
// are not present in the provided CommentController or CommentService.
// If they were, tests similar to those in user/controller_test.go or discussion/controller_test.go
// for Update/Delete (including AuthZ checks for author) would be added here.
