package main

import (
	"log"

	"task-manager-api/config"
	"task-manager-api/controller"
	"task-manager-api/repository"
	"task-manager-api/router"
	"task-manager-api/service"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	config.ConnectDB()
	dbClient := config.GetDB()

	taskCollection := dbClient.Database("task-manager").Collection("tasks")
	userCollection := dbClient.Database("task-manager").Collection("users")

	taskRepo := repository.NewTaskRepository(taskCollection)
	userRepo := repository.NewUserRepository(userCollection)

	taskService := service.NewTaskService(taskRepo)
	userService := service.NewUserService(userRepo)

	taskController := controller.NewTaskController(taskService)
	userController := controller.NewUserController(userService)

	server := router.NewRouter(taskController, userController)
	if err := server.Run(":8080"); err != nil {
		log.Fatal("Server Run Failed", err)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
