package mongo_custom

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type SingleResult struct {
	result *mongo.SingleResult
}

func NewSingleResult(result *mongo.SingleResult) *SingleResult {
	return &SingleResult{
		result: result,
	}
}

func (sr *SingleResult) Decode(v interface{}) error {
	return sr.result.Decode(v)
}
