package domains

import (
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
