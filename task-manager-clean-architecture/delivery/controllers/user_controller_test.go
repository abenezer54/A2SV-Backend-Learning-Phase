package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"task-manager-api/domains"
	"task-manager-api/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
	suite.Suite
	router         *gin.Engine
	userUsecase    *mocks.UserUsecase
	userController *UserController
}

func (suite *UserControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.userUsecase = new(mocks.UserUsecase)
	suite.userController = NewUserController(suite.userUsecase)
	suite.router = gin.Default()
}

func (suite *UserControllerTestSuite) TestRegister_Success() {
	reqBody := map[string]string{
		"username": "testuser",
		"password": "testpassword",
		"role":     "user",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	suite.userUsecase.On("RegisterUser", context.Background(), reqBody["username"], reqBody["password"], reqBody["role"]).Return(&domains.User{
		Username: reqBody["username"],
		Role:     reqBody["role"],
	}, nil)

	rec := httptest.NewRecorder()
	suite.router.POST("/register", suite.userController.Register)
	suite.router.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	suite.userUsecase.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestRegister_UsernameAlreadyExists() {
	reqBody := map[string]string{
		"username": "testuser",
		"password": "testpassword",
		"role":     "user",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	// Mock the RegisterUser method to return "username already exists" error
	suite.userUsecase.On("RegisterUser", context.Background(), reqBody["username"], reqBody["password"], reqBody["role"]).Return(nil, errors.New("username already exists"))

	rec := httptest.NewRecorder()
	suite.router.POST("/register", suite.userController.Register)
	suite.router.ServeHTTP(rec, req)

	// Assert that the response status code is HTTP 409 Conflict
	assert.Equal(suite.T(), http.StatusConflict, rec.Code)

	// Optionally, assert that the correct error message was returned in the response body
	var responseBody map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	suite.NoError(err)
	assert.Equal(suite.T(), "Username already exists", responseBody["error"])

	// Ensure the mock expectations were met
	suite.userUsecase.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestLogin_Success() {
	reqBody := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	user := &domains.User{
		Username: reqBody["username"],
		Password: reqBody["password"],
	}

	suite.userUsecase.On("AuthenticateUser", context.Background(), reqBody["username"], reqBody["password"]).Return(user, true)

	rec := httptest.NewRecorder()
	suite.router.POST("/login", suite.userController.Login)
	suite.router.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	suite.userUsecase.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestLogin_InvalidCredentials() {
	reqBody := map[string]string{
		"username": "wronguser",
		"password": "wrongpassword",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	suite.userUsecase.On("AuthenticateUser", context.Background(), reqBody["username"], reqBody["password"]).Return(nil, false)

	rec := httptest.NewRecorder()
	suite.router.POST("/login", suite.userController.Login)
	suite.router.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, rec.Code)
	suite.userUsecase.AssertExpectations(suite.T())
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
