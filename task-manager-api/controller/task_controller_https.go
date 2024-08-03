package controller

import (
	"net/http"

	"task-manager-api/models"
	"task-manager-api/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RouteHandler struct {
	TaskService *service.TaskService
}

func NewRouteHandler(taskService *service.TaskService) *RouteHandler {
	return &RouteHandler{
		TaskService: taskService,
	}
}

func (rh *RouteHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := rh.TaskService.CreateTask(&task); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (rh *RouteHandler) GetTaskByID(c *gin.Context) {
	taskID := c.Param("id")
	task, err := rh.TaskService.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (rh *RouteHandler) GetAllTasks(c *gin.Context) {
	tasks, err := rh.TaskService.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (rh *RouteHandler) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	var task models.Task
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

func (rh *RouteHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if err := rh.TaskService.DeleteTask(taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
