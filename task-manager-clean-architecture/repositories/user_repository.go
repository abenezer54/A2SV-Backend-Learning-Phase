package repositories

import (
	"context"

	"task-manager-api/domains"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryMongo struct {
	collection *mongo.Collection
}

func NewUserRepositoryMongo(collection *mongo.Collection) *UserRepositoryMongo {
	return &UserRepositoryMongo{
		collection: collection,
	}
}

func (r *UserRepositoryMongo) CreateUser(ctx context.Context, user *domains.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepositoryMongo) FindUserByUsername(ctx context.Context, username string) (*domains.User, error) {
	var user domains.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryMongo) UserExists(ctx context.Context, username string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
