package repository

import (
	"context"
	"fmt"

	"task-manager-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepoMongo struct {
	collection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) *TaskRepoMongo {
	return &TaskRepoMongo{
		collection: collection,
	}
}

func (r *TaskRepoMongo) CreateTask(task *models.Task) error {
	_, err := r.collection.InsertOne(context.TODO(), task)
	return err
}

func (r *TaskRepoMongo) GetTaskByID(id string) (*models.Task, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var task models.Task

	err = r.collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&task)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("task not found")
	}
	return &task, err
}

func (r *TaskRepoMongo) GetAllTasks() ([]*models.Task, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tasks []*models.Task
	for cursor.Next(context.TODO()) {
		var task models.Task
		if err = cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepoMongo) UpdateTask(task *models.Task) error {
	_, err := r.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": task.ID},
		bson.M{"$set": task},
	)
	return err
}

func (r *TaskRepoMongo) DeleteTask(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	return err
}
