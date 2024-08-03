package main

import (
	"log"

	"task-manager-api/config"
	"task-manager-api/controller"
	"task-manager-api/repository"
	"task-manager-api/router"
	"task-manager-api/service"
)

func main() {
	config.ConnectDB()
	dbClient := config.GetDB()

	taskCollection := dbClient.Database("task-manager").Collection("tasks")

	taskRepo := repository.NewTaskRepository(taskCollection)

	taskManager := service.NewTaskService(taskRepo)
	routHandler := controller.NewRouteHandler(taskManager)
	server := router.NewRouter(routHandler)
	if err := server.Run(":8080"); err != nil {
		log.Fatal("Server Run Failed", err)
	}
}
