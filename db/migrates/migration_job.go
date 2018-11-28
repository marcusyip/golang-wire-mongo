package migrates

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/sirupsen/logrus"
)

type MigrationJob struct {
	logger *logrus.Logger
	db     *mongo.Database
}

func (job *MigrationJob) Run() error {
	indexView := job.db.Collection("users").Indexes()

	indexName, err := indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{{"username", bsonx.Int32(1)}},
			Options: mongo.NewIndexOptionsBuilder().
				Unique(true).
				Build(),
		},
	)
	if err != nil {
		return err
	}
	fmt.Printf("[migration] Created index %s\n", indexName)
	return nil
}

func NewMigrationJob(
	logger *logrus.Logger,
	db *mongo.Database,
) *MigrationJob {
	return &MigrationJob{
		logger: logger,
		db:     db,
	}
}
