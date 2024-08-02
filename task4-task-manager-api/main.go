package main

import (
	"task-manager-api/controller"
	"task-manager-api/router"
	"task-manager-api/service"
)

func main() {
	taskManager := service.NewTaskManager()
	routHandler := controller.NewRoutHandler(taskManager)
	server := router.NewRouter(routHandler)
	server.Run(":8080")
}
