package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type Model struct {
	ID        objectid.ObjectID `bson:"_id"`
	CreatedAt time.Time         `bson:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at"`
}

func NewModel() *Model {
	return &Model{
		ID: objectid.New(),
	}
}
