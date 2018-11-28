package repositories

import (
	"context"

	"github.com/marcusyip/golang-wire-mongo/models"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
)

type AccessRepository struct {
	db *mongo.Database
}

func (r *AccessRepository) CollectionName() string {
	return "accesses"
}

func (r *AccessRepository) Collection() *mongo.Collection {
	return r.db.Collection(r.CollectionName())
}

func (r *AccessRepository) FindOne(ctx context.Context, filter interface{}) (*models.Access, error) {
	access := &models.Access{}
	err := r.Collection().FindOne(ctx, filter).Decode(access)
	if err != nil {
		return nil, err
	}
	return access, nil
}

func (r *AccessRepository) Create(ctx context.Context, access *models.Access) error {
	_, err := r.Collection().InsertOne(ctx, access)
	if err != nil {
		return err
	}
	return nil
}

func (r *AccessRepository) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	options ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	return r.Collection().UpdateOne(ctx, filter, update, options...)
}

func NewAccessRepository(db *mongo.Database) *AccessRepository {
	return &AccessRepository{
		db: db,
	}
}
