package repositories

import (
	"fmt"

	"github.com/marcusyip/golang-wire-mongo/config"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func NewClient(conf *config.Config) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		conf.Mongodb.Username,
		conf.Mongodb.Password,
		conf.Mongodb.Host,
		conf.Mongodb.Port)
	return mongo.NewClient(uri)
}

func NewDatabase(
	conf *config.Config,
	client *mongo.Client,
) *mongo.Database {
	return client.Database(conf.Mongodb.Database)
}
