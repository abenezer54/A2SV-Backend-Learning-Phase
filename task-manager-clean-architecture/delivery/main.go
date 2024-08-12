package main

import (
	"log"

	"task-manager-api/delivery/controllers"
	"task-manager-api/delivery/routers"
	"task-manager-api/infrastructure"
	"task-manager-api/repositories"
	"task-manager-api/usecases"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	infrastructure.ConnectDB()
	dbClient := infrastructure.GetDB()

	taskCollection := dbClient.Database("task-manager").Collection("tasks")
	userCollection := dbClient.Database("task-manager").Collection("users")

	taskRepo := repositories.NewTaskRepositoryMongo(taskCollection)
	userRepo := repositories.NewUserRepositoryMongo(userCollection)

	taskService := usecases.NewTaskUsecase(taskRepo)
	userService := usecases.NewUserUsecase(userRepo)

	taskController := controllers.NewTaskController(taskService)
	userController := controllers.NewUserController(userService)

	server := routers.NewRouter(taskController, userController)
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
