package controllers

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"task-manager-api/domains"
	"task-manager-api/infrastructure"
	"task-manager-api/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *mocks.TaskUsecase
	controller  *TaskController
	userClaims  *infrastructure.Claims
}

func (suite *TaskControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
	suite.mockUsecase = new(mocks.TaskUsecase)
	suite.controller = NewTaskController(suite.mockUsecase)

	suite.userClaims = &infrastructure.Claims{
		UserID: primitive.NewObjectID().Hex(),
	}
	suite.router.Use(func(c *gin.Context) {
		c.Set("userClaims", suite.userClaims)
		c.Next()
	})
}

func (suite *TaskControllerTestSuite) TestCreateTask() {
	reqBody := map[string]interface{}{
		"title":       "Test Task",
		"description": "This is a test task",
		"status":      "pending",
		"due_date":    time.Now(),
	}
	reqJSON, _ := json.Marshal(reqBody)

	suite.mockUsecase.On("CreateTask", mock.Anything, reqBody["title"], reqBody["description"], mock.Anything, mock.Anything).Return(&domains.Task{
		ID:          primitive.NewObjectID(),
		Title:       reqBody["title"].(string),
		Description: reqBody["description"].(string),
		DueDate:     reqBody["due_date"].(time.Time),
		Completed:   false,
	}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.POST("/tasks", suite.controller.CreateTask)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUsecase.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestGetTaskByCreatorID() {
	taskID := primitive.NewObjectID()
	suite.mockUsecase.On("GetTaskByIDAndCreator", mock.Anything, taskID, mock.Anything).Return(&domains.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "This is a test task",
		DueDate:     time.Now(),
		Completed:   false,
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/tasks/"+taskID.Hex(), nil)
	w := httptest.NewRecorder()
	suite.router.GET("/tasks/:id", suite.controller.GetTaskByCreatorID)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUsecase.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestGetAllTasks() {
	creatorID, _ := primitive.ObjectIDFromHex(suite.userClaims.UserID)

	suite.mockUsecase.On("GetTasksByCreator", mock.Anything, creatorID).Return([]*domains.Task{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Task 1",
			Description: "This is task 1",
			DueDate:     time.Now(),
			Completed:   false,
		},
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	suite.router.GET("/tasks", suite.controller.GetAllTasks)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUsecase.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestUpdateTaskByCreatorID() {
	taskID := primitive.NewObjectID()
	reqBody := map[string]interface{}{
		"title":       "Updated Task",
		"description": "This is an updated task",
		"completed":   true,
		"due_date":    time.Now(),
	}
	reqJSON, _ := json.Marshal(reqBody)

	suite.mockUsecase.On("UpdateTaskByCreatorID", mock.Anything, taskID, mock.Anything, reqBody["title"], reqBody["description"], reqBody["completed"], mock.MatchedBy(func(t time.Time) bool {
		return t.Equal(reqBody["due_date"].(time.Time))
	})).Return(&domains.Task{
		ID:          taskID,
		Title:       reqBody["title"].(string),
		Description: reqBody["description"].(string),
		DueDate:     reqBody["due_date"].(time.Time),
		Completed:   reqBody["completed"].(bool),
	}, nil)

	req, _ := http.NewRequest(http.MethodPut, "/auth/tasks/"+taskID.Hex(), bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.PUT("/auth/tasks/:id", suite.controller.UpdateTaskByCreatorID)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUsecase.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestDeleteTaskByCreatorID() {
	taskID := primitive.NewObjectID()

	suite.mockUsecase.On("DeleteTaskByCreatorID", mock.Anything, taskID, mock.Anything).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/"+taskID.Hex(), nil)
	w := httptest.NewRecorder()
	suite.router.DELETE("/tasks/:id", suite.controller.DeleteTaskByCreatorID)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUsecase.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestGetTaskByCreatorID_NotFound() {
	taskID := primitive.NewObjectID()
	suite.mockUsecase.On("GetTaskByIDAndCreator", mock.Anything, taskID, mock.Anything).Return(nil, mongo.ErrNoDocuments)

	req, _ := http.NewRequest(http.MethodGet, "/auth/tasks/"+taskID.Hex(), nil)
	w := httptest.NewRecorder()
	suite.router.GET("/auth/tasks/:id", suite.controller.GetTaskByCreatorID)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	suite.mockUsecase.AssertExpectations(suite.T())
}
func (suite *TaskControllerTestSuite) TestCreateTask_ValidationError() {
	reqBody := map[string]interface{}{
		"title": "Test Task", // Missing other required fields
	}
	reqJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.POST("/tasks", suite.controller.CreateTask)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func TestTaskControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TaskControllerTestSuite))
}
