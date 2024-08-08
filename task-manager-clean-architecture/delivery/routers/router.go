package routers

import (
	"task-manager-api/delivery/controllers"
	"task-manager-api/infrastructure"

	"github.com/gin-gonic/gin"
)

func NewRouter(tc *controllers.TaskController, uc *controllers.UserController) *gin.Engine {
	router := gin.Default()
	router.POST("/register", uc.Register)
	router.POST("/login", uc.Login)

	authRoutes := router.Group("/auth")
	authRoutes.Use(infrastructure.JWTAuthMiddleware())
	{
		authRoutes.GET("/tasks", tc.GetAllTasks)
		authRoutes.GET("/tasks/:id", tc.GetTaskByCreatorID)
		authRoutes.POST("/tasks", tc.CreateTask)
		authRoutes.PUT("/tasks/:id", tc.UpdateTaskByCreatorID)
		authRoutes.DELETE("/tasks/:id", tc.DeleteTaskByCreatorID)
	}

	return router
}
