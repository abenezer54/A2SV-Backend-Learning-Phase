package repository

import (
	"task-manager-api/models"
)

type TaskRepository interface {
	CreateTask(task *models.Task) error
	GetTaskByID(id string) (*models.Task, error)
	GetAllTasks() ([]*models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id string) error
}
