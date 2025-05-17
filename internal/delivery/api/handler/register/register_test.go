package register_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Myakun/personal-secretary-user-api/internal/delivery/api/handler/register"
	"github.com/Myakun/personal-secretary-user-api/internal/presentation/user/registration"
	"github.com/Myakun/personal-secretary-user-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for the Logger interface
type MockLogger struct {
	mock.Mock
}

// Ensure MockLogger implements logger.Logger
var _ logger.Logger = (*MockLogger)(nil)

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Fatal(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Warning(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) DebugW(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) ErrorW(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) FatalW(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) InfoW(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) WarningW(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) DebugWithTag(msg string, tag string) {
	m.Called(msg, tag)
}

func (m *MockLogger) ErrorWithTag(msg string, tag string) {
	m.Called(msg, tag)
}

func (m *MockLogger) FatalWithTag(msg string, tag string) {
	m.Called(msg, tag)
}

func (m *MockLogger) InfoWithTag(msg string, tag string) {
	m.Called(msg, tag)
}

func (m *MockLogger) WarningWithTag(msg string, tag string) {
	m.Called(msg, tag)
}

func (m *MockLogger) DebugWithTagW(msg string, tag string, keysAndValues ...interface{}) {
	m.Called(msg, tag, keysAndValues)
}

func (m *MockLogger) ErrorWithTagW(msg string, tag string, keysAndValues ...interface{}) {
	m.Called(msg, tag, keysAndValues)
}

func (m *MockLogger) FatalWithTagW(msg string, tag string, keysAndValues ...interface{}) {
	m.Called(msg, tag, keysAndValues)
}

func (m *MockLogger) InfoWithTagW(msg string, tag string, keysAndValues ...interface{}) {
	m.Called(msg, tag, keysAndValues)
}

func (m *MockLogger) WarningWithTagW(msg string, tag string, keysAndValues ...interface{}) {
	m.Called(msg, tag, keysAndValues)
}

// Mock for the UserRegistration interface
type MockUserRegistration struct {
	mock.Mock
}

func (m *MockUserRegistration) RegisterUser(ctx context.Context, request registration.RegisterUserRequest) (*registration.RegisterUserResult, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*registration.RegisterUserResult), args.Error(1)
}

func TestRegisterHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		requestBody        map[string]interface{}
		setupMocks         func(*MockLogger, *MockUserRegistration)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "Successful registration",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"name":     "Test User",
				"password": "password123",
			},
			setupMocks: func(mockLogger *MockLogger, mockUserReg *MockUserRegistration) {
				successResponse := &registration.SuccessResponse{
					Token:        "jwt-token",
					RefreshToken: "refresh-token",
				}
				result := &registration.RegisterUserResult{
					Success:         true,
					SuccessResponse: successResponse,
				}

				mockUserReg.On("RegisterUser", mock.Anything, registration.RegisterUserRequest{
					Email:    "test@example.com",
					Name:     "Test User",
					Password: "password123",
				}).Return(result, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"token":         "jwt-token",
				"refresh_token": "refresh-token",
			},
		},
		{
			name:        "Invalid JSON input",
			requestBody: map[string]interface{}{
				// Missing required fields
			},
			setupMocks: func(mockLogger *MockLogger, mockUserReg *MockUserRegistration) {
				mockLogger.On("DebugWithTag", mock.Anything, "API_USER_REGISTER").Return()
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
		},
		{
			name: "Registration error",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"name":     "Test User",
				"password": "password123",
			},
			setupMocks: func(mockLogger *MockLogger, mockUserReg *MockUserRegistration) {
				mockUserReg.On("RegisterUser", mock.Anything, registration.RegisterUserRequest{
					Email:    "test@example.com",
					Name:     "Test User",
					Password: "password123",
				}).Return(nil, errors.New("registration error"))

				mockLogger.On("FatalWithTag", mock.Anything, "API_USER_REGISTER").Return()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
		{
			name: "Validation error",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"name":     "Test User",
				"password": "password123",
			},
			setupMocks: func(mockLogger *MockLogger, mockUserReg *MockUserRegistration) {
				errorResponse := &registration.ErrorResponse{
					Err: "validation error",
				}
				result := &registration.RegisterUserResult{
					Success:       false,
					ErrorResponse: errorResponse,
				}

				mockUserReg.On("RegisterUser", mock.Anything, registration.RegisterUserRequest{
					Email:    "test@example.com",
					Name:     "Test User",
					Password: "password123",
				}).Return(result, nil)
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse: map[string]interface{}{
				"error": "validation error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockLogger := new(MockLogger)
			mockUserReg := new(MockUserRegistration)
			tt.setupMocks(mockLogger, mockUserReg)

			// Create handler
			handler := register.NewRegisterHandler(mockLogger, mockUserReg)

			// Setup router
			router := gin.New()
			router.POST("/register", handler.Register)

			// Create request
			requestJSON, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestJSON))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Assert response body if expected
			if tt.expectedResponse != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			}

			// Verify mocks
			mockLogger.AssertExpectations(t)
			mockUserReg.AssertExpectations(t)
		})
	}
}
