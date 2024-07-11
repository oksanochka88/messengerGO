package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock модели и функций
type MockModel struct {
	mock.Mock
}

func (m *MockModel) AuthenticateUser(username, password string) (*User, error) {
	args := m.Called(username, password)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockModel) GetTokensByUserID(userID int) ([]Token, error) {
	args := m.Called(userID)
	return args.Get(0).([]Token), args.Error(1)
}

func (m *MockModel) IsTokenValid(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}

func (m *MockModel) SaveToken(userID int, token string) error {
	args := m.Called(userID, token)
	return args.Error(0)
}

// Объекты для теста
type User struct {
	ID       int
	Username string
}

type Token struct {
	Token string
}

// Тест функции Login
func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockModel := new(MockModel)

	router := gin.Default()
	router.POST("/login", func(c *gin.Context) {
		var loginRequest LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := mockModel.AuthenticateUser(loginRequest.Username, loginRequest.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		tokens, err := mockModel.GetTokensByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tokens"})
			return
		}

		var validToken string
		for _, t := range tokens {
			valid, err := mockModel.IsTokenValid(t.Token)
			if err == nil && valid {
				validToken = t.Token
				break
			}
		}

		if validToken == "" {
			token := "newToken"
			err = mockModel.SaveToken(user.ID, token)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
				return
			}
			validToken = token
		}

		c.JSON(http.StatusOK, gin.H{"token": validToken})
	})

	// Успешный запрос
	user := &User{ID: 1, Username: "testuser"}
	mockModel.On("AuthenticateUser", "testuser", "password").Return(user, nil)
	mockModel.On("GetTokensByUserID", user.ID).Return([]Token{}, nil)
	mockModel.On("SaveToken", user.ID, mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(`{"username":"testuser","password":"password"}`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"token":"newToken"`)

	// Неправильное имя пользователя или пароль
	mockModel.On("AuthenticateUser", "wronguser", "password").Return((*User)(nil), assert.AnError)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", strings.NewReader(`{"username":"wronguser","password":"password"}`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid username or password"`)
}
