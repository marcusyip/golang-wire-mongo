package repositories

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Ops struct {
	db *mongo.Database
}

func (ops *Ops) FindOne(ctx context.Context, r Repository, filter interface{}) (*bson.D, error) {
	c := ops.db.Collection(r.CollectionName())
	result := &bson.D{}
	err := c.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ops *Ops) Create(ctx context.Context, r Repository, doc interface{}) {

}

func NewOps(db *mongo.Database) *Ops {
	return &Ops{
		db: db,
	}
}
