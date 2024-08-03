package router

import (
	"task-manager-api/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(rh *controller.RouteHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/tasks", rh.GetAllTasks)
	router.GET("/tasks/:id", rh.GetTaskByID)
	router.POST("/tasks", rh.CreateTask)
	router.PUT("/tasks/:id", rh.UpdateTask)
	router.DELETE("/tasks/:id", rh.DeleteTask)

	return router
}
