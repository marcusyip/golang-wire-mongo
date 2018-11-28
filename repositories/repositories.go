package repositories

import (
	"github.com/google/go-cloud/wire"
)

var WireSet = wire.NewSet(
	NewClient,
	NewDatabase,
	NewUserRepository,
	NewAccessRepository,
)

type Repository interface {
	CollectionName() string
}
