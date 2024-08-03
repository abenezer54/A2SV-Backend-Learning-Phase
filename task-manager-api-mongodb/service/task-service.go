package service

import (
	"errors"

	"task-manager-api/models"
	"task-manager-api/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (ts *TaskService) CreateTask(task *models.Task) error {
	existingTask, err := ts.repo.GetTaskByID(task.ID.Hex())
	if err == nil && existingTask != nil {
		return errors.New("task already exists")
	}
	return ts.repo.CreateTask(task)
}

func (ts *TaskService) GetTaskByID(id string) (*models.Task, error) {
	task, err := ts.repo.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (ts *TaskService) GetAllTasks() ([]*models.Task, error) {
	tasks, err := ts.repo.GetAllTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ts *TaskService) UpdateTask(task *models.Task) error {
	existingTask, err := ts.repo.GetTaskByID(task.ID.Hex())
	if err != nil {
		return errors.New("couldn't find task")
	}
	if existingTask == nil {
		return errors.New("task not found")
	}
	return ts.repo.UpdateTask(task)
}

func (ts *TaskService) DeleteTask(taskID string) error {
	existingTask, err := ts.repo.GetTaskByID(taskID)
	if err != nil {
		return errors.New("couldn't find task")
	}
	if existingTask == nil {
		return errors.New("task not found")
	}
	return ts.repo.DeleteTask(taskID)
}
