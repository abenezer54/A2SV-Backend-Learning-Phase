package router

import (
	"task-manager-api/controller"
	"task-manager-api/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(tc *controller.TaskController, uc *controller.UserController) *gin.Engine {
	router := gin.Default()
	router.POST("/register", uc.Register)
	router.POST("/login", uc.Login)

	authRoutes := router.Group("/auth")
	authRoutes.Use(middleware.JWTAuthMiddleware())
	{
		authRoutes.GET("/tasks", tc.GetAllTasks)
		authRoutes.GET("/tasks/:id", tc.GetTaskByCreatorID)
		authRoutes.POST("/tasks", tc.CreateTask)
		authRoutes.PUT("/tasks/:id", tc.UpdateTaskByCreatorID)
		authRoutes.DELETE("/tasks/:id", tc.DeleteTaskByCreatorID)
	}

	return router
}
