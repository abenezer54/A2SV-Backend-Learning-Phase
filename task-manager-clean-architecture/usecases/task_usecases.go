package usecases

import (
	"context"
	"errors"
	"time"

	"task-manager-api/domains"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskUsecase struct {
	repo domains.TaskRepository
}

func NewTaskUsecase(repo domains.TaskRepository) *TaskUsecase {
	return &TaskUsecase{
		repo: repo,
	}
}

func (ts *TaskUsecase) CreateTask(ctx context.Context, title, description string, dueDate time.Time, creatorID primitive.ObjectID) (*domains.Task, error) {
	task := &domains.Task{
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

func (ts *TaskUsecase) GetTasksByCreator(ctx context.Context, creatorID primitive.ObjectID) ([]*domains.Task, error) {
	return ts.repo.FindTasksByCreator(ctx, creatorID)
}

func (ts *TaskUsecase) GetTaskByIDAndCreator(ctx context.Context, taskID, creatorID primitive.ObjectID) (*domains.Task, error) {
	return ts.repo.FindTaskByIDAndCreator(ctx, taskID, creatorID)
}

// Admin
func (ts *TaskUsecase) GetTaskByID(id string) (*domains.Task, error) {
	task, err := ts.repo.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Admin
func (ts *TaskUsecase) GetAllTasks() ([]*domains.Task, error) {
	tasks, err := ts.repo.GetAllTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ts *TaskUsecase) UpdateTaskByCreatorID(ctx context.Context, taskID, creatorID primitive.ObjectID, title, description string, completed bool, dueDate time.Time) (*domains.Task, error) {
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
func (ts *TaskUsecase) UpdateTask(task *domains.Task) error {
	existingTask, err := ts.repo.GetTaskByID(task.ID.Hex())
	if err != nil {
		return errors.New("couldn't find task")
	}
	if existingTask == nil {
		return errors.New("task not found")
	}
	return ts.repo.UpdateTask(task)
}

func (ts *TaskUsecase) DeleteTaskByCreatorID(ctx context.Context, taskID, creatorID primitive.ObjectID) error {
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

func (ts *TaskUsecase) DeleteTask(taskID string) error {
	existingTask, err := ts.repo.GetTaskByID(taskID)
	if err != nil {
		return errors.New("couldn't find task")
	}
	if existingTask == nil {
		return errors.New("task not found")
	}
	return ts.repo.DeleteTask(taskID)
}
