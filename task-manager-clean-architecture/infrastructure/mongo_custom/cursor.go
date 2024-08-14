package mongo_custom

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCursor struct {
	Cursor *mongo.Cursor
}

func NewMongoCursor(cursor *mongo.Cursor) *MongoCursor {
	return &MongoCursor{
		Cursor: cursor,
	}
}

func (c *MongoCursor) Close(ctx context.Context) error {
	return c.Cursor.Close(ctx)
}

func (c *MongoCursor) Next(ctx context.Context) bool {
	return c.Cursor.Next(ctx)
}

func (c *MongoCursor) Decode(v interface{}) error {
	return c.Cursor.Decode(v)
}

func (c *MongoCursor) Err() error {
	return c.Cursor.Err()
}
