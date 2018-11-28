//+build wireinject

package migrate

import (
	"github.com/google/go-cloud/wire"
	"github.com/marcusyip/golang-wire-mongo/config"
	"github.com/marcusyip/golang-wire-mongo/core/logger"
	"github.com/marcusyip/golang-wire-mongo/db/migrates"
	repos "github.com/marcusyip/golang-wire-mongo/repositories"
)

func NewMigrationJob() (*migrates.MigrationJob, error) {
	wire.Build(
		config.ProvideConfig,
		logger.New,
		repos.WireSet,
		migrates.NewMigrationJob,
	)
	return nil, nil
}
