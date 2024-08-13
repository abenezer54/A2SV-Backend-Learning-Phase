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
func (ts *TaskUsecase) GetTaskByID(ctx context.Context, id string) (*domains.Task, error) {
	task, err := ts.repo.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	return task, nil
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
