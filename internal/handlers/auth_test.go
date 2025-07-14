package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignup_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.MockUserRepository)
	mockUserService := new(mocks.MockUserService)
	handler := NewAuthHandler(mockUserService)

	tenantID := uuid.New().String()

	signupReq := SignupRequest{
		Email:    "test@example.com",
		Password: "test_password",
		TenantID: tenantID,
	}

	body, _ := json.Marshal(signupReq)
	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/signup", handler.SignUp)

	// setup expectation on the mock
	mockRepo.
		On("CreateUser", mock.AnythingOfType("*models.User")).
		Return(nil)

	// Perform request
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "user created successfully")

	// Verify that CreateUser was called
	mockRepo.AssertExpectations(t)
}
