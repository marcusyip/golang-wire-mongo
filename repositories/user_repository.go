package repositories

import (
	"context"
	"time"

	"github.com/marcusyip/golang-wire-mongo/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
)

type UserRepository struct {
	db *mongo.Database
}

func (r *UserRepository) CollectionName() string {
	return "users"
}

func (r *UserRepository) Collection() *mongo.Collection {
	return r.db.Collection(r.CollectionName())
}

func (r *UserRepository) FindOne(ctx context.Context, filter interface{}) (*models.User, error) {
	user := &models.User{}
	err := r.Collection().FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id objectid.ObjectID) (*models.User, error) {
	filter := bson.D{{"_id", id}}
	return r.FindOne(ctx, filter)
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	// user.CreatedAt = time.Now()
	// user.UpdatedAt = time.Now()
	_, err := r.Collection().InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateOne(ctx context.Context, filter bson.D, update bson.D,
	options ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	update = append(update, bson.E{"updated_at", time.Now})
	return r.Collection().UpdateOne(ctx, filter, update, options...)
}

func (r *UserRepository) UpdateByID(ctx context.Context, id objectid.ObjectID, update bson.D,
	options ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	filter := bson.D{{"_id", id}}
	return r.UpdateOne(ctx, filter, update, options...)
}

func NewUserRepository(
	db *mongo.Database,
) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
