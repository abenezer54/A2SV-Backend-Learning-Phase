package domains

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Completed   bool               `bson:"completed" json:"completed"`
	DueDate     time.Time          `bson:"due_date" json:"due_date"`
	CreatorID   primitive.ObjectID `bson:"creator_id" json:"creator_id"`
}

func NewTask(title, description string, completed bool, dueDate time.Time, creatorId primitive.ObjectID) *Task {
	return &Task{
		ID:          primitive.NewObjectID(),
		Title:       title,
		Description: description,
		Completed:   completed,
		DueDate:     dueDate,
		CreatorID:   creatorId,
	}
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task *Task) (*Task, error)
	FindTasksByCreator(ctx context.Context, creatorID primitive.ObjectID) ([]*Task, error)
	FindTaskByIDAndCreator(ctx context.Context, taskID, creatorID primitive.ObjectID) (*Task, error)
	GetTaskByID(id string) (*Task, error)
	UpdateTaskByCreatorID(ctx context.Context, task *Task) error
	DeleteTaskByCreatorID(ctx context.Context, taskID, creatorID primitive.ObjectID) error
}

type TaskUsecase interface {
	CreateTask(ctx context.Context, title, description string, dueDate time.Time, creatorID primitive.ObjectID) (*Task, error)
	GetTasksByCreator(ctx context.Context, creatorID primitive.ObjectID) ([]*Task, error)
	GetTaskByIDAndCreator(ctx context.Context, taskID, creatorID primitive.ObjectID) (*Task, error)
	GetTaskByID(ctx context.Context, id string) (*Task, error)
	UpdateTaskByCreatorID(ctx context.Context, taskID, creatorID primitive.ObjectID, title, description string, completed bool, dueDate time.Time) (*Task, error)
	DeleteTaskByCreatorID(ctx context.Context, taskID, creatorID primitive.ObjectID) error
}
