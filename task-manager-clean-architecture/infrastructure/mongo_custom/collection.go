package mongo_custom

import (
	"context"
	"task-manager-api/domains"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollection struct {
	*mongo.Collection
}
func NewMongoCollection(collection *mongo.Collection) *MongoCollection {
	return &MongoCollection{
		Collection: collection,
	}
}

func (c *MongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.Collection.InsertOne(ctx, document, opts...)
}
func (c *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur domains.CursorInterface, err error) {
	cursor, err := c.Collection.Find(ctx, filter, opts...)
	return NewMongoCursor(cursor), err
}

func (c *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) domains.SingleResultInterface {
	res := c.Collection.FindOne(ctx, filter, opts...)
	return NewSingleResult(res)
}

func (c *MongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.Collection.UpdateOne(ctx, filter, update, opts...)
}

func (c *MongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.Collection.DeleteOne(ctx, filter, opts...)
}

func (c *MongoCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return c.Collection.CountDocuments(ctx, filter, opts...)
}
