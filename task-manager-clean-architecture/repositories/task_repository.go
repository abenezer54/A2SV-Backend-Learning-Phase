package repositories

import (
	"context"
	"fmt"

	"task-manager-api/domains"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepoMongo struct {
	collection *mongo.Collection
}

func NewTaskRepositoryMongo(collection *mongo.Collection) *TaskRepoMongo {
	return &TaskRepoMongo{
		collection: collection,
	}
}

func (r *TaskRepoMongo) CreateTask(ctx context.Context, task *domains.Task) (*domains.Task, error) {
	_, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *TaskRepoMongo) FindTasksByCreator(ctx context.Context, creatorID primitive.ObjectID) ([]*domains.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"creator_id": creatorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*domains.Task
	for cursor.Next(ctx) {
		var task domains.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepoMongo) FindTaskByIDAndCreator(ctx context.Context, taskID, creatorID primitive.ObjectID) (*domains.Task, error) {
	filter := bson.M{"_id": taskID, "creator_id": creatorID}
	var task domains.Task
	err := r.collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepoMongo) GetTaskByID(id string) (*domains.Task, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var task domains.Task

	err = r.collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&task)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("task not found")
	}
	return &task, err
}

func (r *TaskRepoMongo) GetAllTasks() ([]*domains.Task, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tasks []*domains.Task
	for cursor.Next(context.TODO()) {
		var task domains.Task
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

func (r *TaskRepoMongo) UpdateTaskByCreatorID(ctx context.Context, task *domains.Task) error {
	filter := bson.M{"_id": task.ID, "creator_id": task.CreatorID}
	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"completed":   task.Completed,
			"due_date":    task.DueDate,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Admin
func (r *TaskRepoMongo) UpdateTask(task *domains.Task) error {
	_, err := r.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": task.ID},
		bson.M{"$set": task},
	)
	return err
}

func (r *TaskRepoMongo) DeleteTaskByCreatorID(ctx context.Context, taskID, creatorID primitive.ObjectID) error {
	filter := bson.M{"_id": taskID, "creator_id": creatorID}
	_, err := r.collection.DeleteOne(ctx, filter)
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
