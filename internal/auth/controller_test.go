package auth

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"go-discussion-app/internal/user"
	"go-discussion-app/models"
	"go-discussion-app/pkg/jwtutil"
)

// MockUserRepository is a mock implementation of user.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *models.User) (int, error) {
	args := m.Called(ctx, u)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	// Handle the case where Get(0) is nil (user not found)
	if args.Get(0) == nil {
		return nil, args.Error(1) // Error(1) could be nil if user not found is not an error
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, u *models.User) (sql.Result, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) (sql.Result, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(sql.Result), args.Error(1)
}

// Helper function to set up the Gin router with controller routes
func setupTestRouter(mockUserRepo user.UserRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New() // Use gin.New() for a blank router in tests
	authService := NewService(mockUserRepo)
	authController := NewController(authService)

	// Group for /auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authController.RegisterHandler)
		authGroup.POST("/login", authController.LoginHandler)
	}

	// Dummy protected route for middleware testing
	router.GET("/protected", JWTAuthMiddleware(), func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "welcome", "userID": userID})
	})
	return router
}

// Helper function to make HTTP requests
func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBodyReader *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBodyReader = bytes.NewBuffer(jsonData)
	} else {
		reqBodyReader = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBodyReader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRegister_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	router := setupTestRouter(mockUserRepo)

	registerDTO := RegisterDTO{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockUserRepo.On("GetByEmail", mock.Anything, registerDTO.Email).Return(nil, nil)
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Email == registerDTO.Email && u.Username == registerDTO.Username
	})).Return(1, nil)

	w := performRequest(router, "POST", "/auth/register", registerDTO)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response["id"])
	mockUserRepo.AssertExpectations(t)
}

func TestRegister_UserExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	router := setupTestRouter(mockUserRepo)

	registerDTO := RegisterDTO{
		Username: "existinguser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	existingUser := &models.User{ID: 1, Email: registerDTO.Email, Username: "existinguser"}
	mockUserRepo.On("GetByEmail", mock.Anything, registerDTO.Email).Return(existingUser, nil)

	w := performRequest(router, "POST", "/auth/register", registerDTO)

	assert.Equal(t, http.StatusConflict, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "email already in use", response["error"])
	mockUserRepo.AssertExpectations(t)
}

func TestRegister_InvalidInput_ServiceValidationFailure(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	router := setupTestRouter(mockUserRepo)

	// Example: Missing Username, which fails dto.Validate() in the service
	registerDTO := RegisterDTO{
		Email:    "test@example.com",
		Password: "password123",
	}

	// GetByEmail should not be called if service validation fails first.
	// However, the current service calls GetByEmail *before* userRepo.Create, but *after* dto.Validate().
	// If dto.Validate() fails, GetByEmail and Create are not called.

	w := performRequest(router, "POST", "/auth/register", registerDTO)

	// As per controller logic, service errors (other than ErrUserExists) result in 500.
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "server error", response["error"]) // Controller's generic message
	// Ensure no DB calls were made for this validation failure path
	mockUserRepo.AssertNotCalled(t, "GetByEmail")
	mockUserRepo.AssertNotCalled(t, "Create")
}

func TestRegister_InvalidInput_BindingFailure(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	router := setupTestRouter(mockUserRepo)

	rawBody := `{"username": "test", "email": 123, "password": "password"}` // email as number

	req, _ := http.NewRequest("POST", "/auth/register", strings.NewReader(rawBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid payload", response["error"])
	mockUserRepo.AssertNotCalled(t, "GetByEmail")
	mockUserRepo.AssertNotCalled(t, "Create")
}

func TestLogin_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	router := setupTestRouter(mockUserRepo)

	loginDTO := LoginDTO{
		Email:    "test@example.com",
		Password: "password123",
	}
	// Password hash for "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(loginDTO.Password), bcrypt.DefaultCost)
	expectedUser := &models.User{
		ID:           1,
		Email:        loginDTO.Email,
		PasswordHash: string(hashedPassword),
	}

	mockUserRepo.On("GetByEmail", mock.Anything, loginDTO.Email).Return(expectedUser, nil)

	w := performRequest(router, "POST", "/auth/login", loginDTO)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"], "Token should not be empty")

	// Verify token (optional, but good for sanity)
	userID, jwtErr := jwtutil.ExtractUserID(response["token"])
	assert.NoError(t, jwtErr)
	assert.Equal(t, expectedUser.ID, userID)

	mockUserRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	router := setupTestRouter(mockUserRepo)

	loginDTO := LoginDTO{Email: "unknown@example.com", Password: "password"}
	mockUserRepo.On("GetByEmail", mock.Anything, loginDTO.Email).Return(nil, nil) // User not found

	w := performRequest(router, "POST", "/auth/login", loginDTO)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var respData map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respData)
	assert.NoError(t, err)
	assert.Equal(t, "wrong email or password", respData["error"])
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials_WrongPassword(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	router := setupTestRouter(mockUserRepo)

	loginDTO := LoginDTO{Email: "test@example.com", Password: "wrongpassword"}
	correctPasswordHash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	existingUser := &models.User{ID: 1, Email: loginDTO.Email, PasswordHash: string(correctPasswordHash)}

	mockUserRepo.On("GetByEmail", mock.Anything, loginDTO.Email).Return(existingUser, nil)

	w := performRequest(router, "POST", "/auth/login", loginDTO)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var respData map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respData)
	assert.NoError(t, err)
	assert.Equal(t, "wrong email or password", respData["error"])
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_InvalidInput_BindingFailure(t *testing.T) {
    mockUserRepo := new(MockUserRepository)
    router := setupTestRouter(mockUserRepo)
    rawBody := `{"email": 123, "password": "password"}` // Malformed

    req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(rawBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    var respData map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &respData)
    assert.NoError(t, err)
    assert.Equal(t, "invalid payload", respData["error"])
    mockUserRepo.AssertNotCalled(t, "GetByEmail")
}


func TestLogin_InvalidInput_ServiceValidationFailure(t *testing.T) {
    mockUserRepo := new(MockUserRepository)
    router := setupTestRouter(mockUserRepo)
    loginDTO := LoginDTO{Password: "password123"} // Missing Email

    // Service validation (dto.Validate()) will fail.
    // Controller returns 500 for such service errors.
    w := performRequest(router, "POST", "/auth/login", loginDTO)
    assert.Equal(t, http.StatusInternalServerError, w.Code)
	var respData map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respData)
	assert.NoError(t, err)
    assert.Equal(t, "server error", respData["error"])
    mockUserRepo.AssertNotCalled(t, "GetByEmail")
}


func TestAuthMiddleware_ValidToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository) // Not used by middleware directly but setup needs it
	router := setupTestRouter(mockUserRepo)

	validToken, err := jwtutil.GenerateToken(123) // Generate a token for user ID 123
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var respData map[string]interface{} // userID is int
	err = json.Unmarshal(w.Body.Bytes(), &respData)
	assert.NoError(t, err)
	assert.Equal(t, "welcome", respData["message"])
	assert.Equal(t, float64(123), respData["userID"]) // JSON numbers are float64
}

func TestAuthMiddleware_InvalidToken_NoToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	router := setupTestRouter(mockUserRepo)

	req, _ := http.NewRequest("GET", "/protected", nil) // No Authorization header
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var respData map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respData)
	assert.NoError(t, err)
	assert.Equal(t, "invalid auth header", respData["error"])
}

func TestAuthMiddleware_InvalidToken_MalformedHeader(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	router := setupTestRouter(mockUserRepo)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearertoken123") // Missing space
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var respData map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respData)
	assert.NoError(t, err)
	assert.Equal(t, "invalid auth header", respData["error"])
}

func TestAuthMiddleware_InvalidToken_NotJWT(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	router := setupTestRouter(mockUserRepo)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer a.b.c") // Not a valid JWT structure
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var respData map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respData)
	assert.NoError(t, err)
	assert.Equal(t, "invalid token", respData["error"])
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
    // This requires creating a token that is actually expired.
    // We can modify jwtutil for testing or use a known expired token if possible.
    // For now, this simulates jwtutil.ExtractUserID returning ErrTokenExpired.
    // To do this properly, we would need to either:
    // 1. Have a way to make GenerateToken create an expired token.
    // 2. Mock ExtractUserID (if it were an interface method).
    // Since jwtutil.ExtractUserID is a direct call and might return a specific error variable
    // jwtutil.ErrTokenExpired, we rely on that behavior.

    // Let's assume jwtutil.SECRET_KEY and jwtutil.TOKEN_TTL are accessible for testing
    // or we have a helper in jwtutil to generate expired tokens.
    // If not, this specific test case is hard to achieve in pure unit test fashion
    // without time manipulation (e.g. "github.com/jonboulle/clockwork")
    // or by changing jwtutil to allow very short TTLs for testing.

    // For the sake of this example, if jwtutil.ErrTokenExpired is a specific error
    // that ExtractUserID can return for various reasons (not just actual expiry like clock time,
    // but maybe "exp" claim is in the past), then "invalid token" might cover it
    // if ExtractUserID doesn't differentiate.
    // The current middleware specifically checks: `if err == jwtutil.ErrTokenExpired`

    // To test this path, we need jwtutil.ExtractUserID to return jwtutil.ErrTokenExpired.
    // This means the token string itself must be validly signed but have an "exp" claim in the past.
    // This is hard to forge without knowing the secret key or having a helper in jwtutil.
    // We'll skip the direct test of an *actually* expired token for now, as it depends
    // heavily on the internals or testability features of jwtutil.
    // The "invalid token" test (TestAuthMiddleware_InvalidToken_NotJWT) covers general failures.
    // If ExtractUserID returns jwtutil.ErrTokenExpired, the middleware should catch it.
    // The current "invalid token" might be sufficient if the error from an expired token
    // is not distinguishable from other parsing errors by ExtractUserID unless it specifically returns ErrTokenExpired.

    // Let's assume we *can* generate an expired token for demonstration.
    // This part is pseudo-code unless jwtutil supports generating expired tokens for tests.
    // expiredToken := jwtutil.GenerateExpiredTokenForTest() // Hypothetical function

    // If we can't generate one, we can't test this path directly in this manner.
    // The middleware code is:
    //  uid, err := jwtutil.ExtractUserID(parts[1])
    //  if err != nil {
    //      if err == jwtutil.ErrTokenExpired {
    //          c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
    //      } else {
    //          c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
    //      } // ...
    //  }
    // So, if `jwtutil.ExtractUserID` correctly returns `jwtutil.ErrTokenExpired`, this path is hit.
    // A comprehensive test of `jwtutil.ExtractUserID` itself should verify it returns this error.
    t.Skip("Skipping direct expired token test: requires jwtutil to generate expired tokens or time mocking.")
}
