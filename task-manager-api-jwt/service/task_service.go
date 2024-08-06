package service

import (
	"context"
	"errors"
	"time"

	"task-manager-api/models"
	"task-manager-api/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskService struct {
	repo *repository.TaskRepoMongo
}

func NewTaskService(repo *repository.TaskRepoMongo) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (ts *TaskService) CreateTask(ctx context.Context, title, description string, dueDate time.Time, creatorID primitive.ObjectID) (*models.Task, error) {
	task := &models.Task{
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		CreatorID:   creatorID,
	}
	existingTask, err := ts.repo.GetTaskByID(task.ID.Hex())
	if err == nil && existingTask != nil {
		return nil, errors.New("task already exists")
	}

	t, err := ts.repo.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (ts *TaskService) GetTasksByCreator(ctx context.Context, creatorID primitive.ObjectID) ([]*models.Task, error) {
	return ts.repo.FindTasksByCreator(ctx, creatorID)
}

func (ts *TaskService) GetTaskByIDAndCreator(ctx context.Context, taskID, creatorID primitive.ObjectID) (*models.Task, error) {
	return ts.repo.FindTaskByIDAndCreator(ctx, taskID, creatorID)
}
//Admin
func (ts *TaskService) GetTaskByID(id string) (*models.Task, error) {
	task, err := ts.repo.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	return task, nil
}
//Admin
func (ts *TaskService) GetAllTasks() ([]*models.Task, error) {
	tasks, err := ts.repo.GetAllTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ts *TaskService) UpdateTaskByCreatorID(ctx context.Context, taskID, creatorID primitive.ObjectID, title, description string, completed bool, dueDate time.Time) (*models.Task, error) {
	task, err := ts.repo.FindTaskByIDAndCreator(ctx, taskID, creatorID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, mongo.ErrNoDocuments
	}

	task.Title = title
	task.Description = description
	task.Completed = completed
	task.DueDate = dueDate

	err = ts.repo.UpdateTaskByCreatorID(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}


// Admin
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

func (ts *TaskService) DeleteTaskByCreatorID(ctx context.Context, taskID, creatorID primitive.ObjectID) error {
	task, err := ts.repo.FindTaskByIDAndCreator(ctx, taskID, creatorID)
	if err != nil {
		return err
	}

	if task == nil {
		return mongo.ErrNoDocuments
	}

	err = ts.repo.DeleteTaskByCreatorID(ctx, taskID, creatorID)
	if err != nil {
		return err
	}

	return nil
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
