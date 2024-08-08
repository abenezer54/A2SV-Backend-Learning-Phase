package controllers

import (
	"net/http"
	"time"

	"task-manager-api/domains"
	"task-manager-api/infrastructure"
	"task-manager-api/usecases"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskController struct {
	TaskService *usecases.TaskService
}

func NewTaskController(taskService *usecases.TaskService) *TaskController {
	return &TaskController{
		TaskService: taskService,
	}
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var req struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description" binding:"required"`
		Status      string    `json:"status" binding:"required"`
		DueDate     time.Time `json:"due_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userClaims := c.MustGet("userClaims").(*infrastructure.Claims)
	creatorID, _ := primitive.ObjectIDFromHex(userClaims.UserID)

	task, err := tc.TaskService.CreateTask(c.Request.Context(), req.Title, req.Description, req.DueDate, creatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) GetTaskByCreatorID(c *gin.Context) {
	taskID := c.Param("id")
	taskObjID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userClaims := c.MustGet("userClaims").(*infrastructure.Claims)
	creatorID, err := primitive.ObjectIDFromHex(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	task, err := tc.TaskService.GetTaskByIDAndCreator(c.Request.Context(), taskObjID, creatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Admins EndPoint
func (tc *TaskController) GetAllTasks(c *gin.Context) {
	userClaims := c.MustGet("userClaims").(*infrastructure.Claims)
	creatorID, err := primitive.ObjectIDFromHex(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	tasks, err := tc.TaskService.GetTasksByCreator(c.Request.Context(), creatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) UpdateTaskByCreatorID(c *gin.Context) {
	taskID := c.Param("id")
	taskObjID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description" binding:"required"`
		Completed   bool      `json:"completed" binding:"required"`
		DueDate     time.Time `json:"due_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userClaims := c.MustGet("userClaims").(*infrastructure.Claims)
	creatorID, err := primitive.ObjectIDFromHex(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	task, err := tc.TaskService.UpdateTaskByCreatorID(c.Request.Context(), taskObjID, creatorID, req.Title, req.Description, req.Completed, req.DueDate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Admin
func (rh *TaskController) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	var task domains.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task.ID = objectID
	if err := rh.TaskService.UpdateTask(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) DeleteTaskByCreatorID(c *gin.Context) {
	taskID := c.Param("id")
	taskObjID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userClaims := c.MustGet("userClaims").(*infrastructure.Claims)
	creatorID, err := primitive.ObjectIDFromHex(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	err = tc.TaskService.DeleteTaskByCreatorID(c.Request.Context(), taskObjID, creatorID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (rh *TaskController) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if err := rh.TaskService.DeleteTask(taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
