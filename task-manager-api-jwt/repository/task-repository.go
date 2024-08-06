package repository

import (
	"context"

	"task-manager-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskRepository interface {
	CreateTask(task *models.Task) error
	FindTasksByCreator(ctx context.Context, creatorID primitive.ObjectID) ([]*models.Task, error)
	GetTaskByID(id string) (*models.Task, error)
	GetAllTasks() ([]*models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id string) error
}
