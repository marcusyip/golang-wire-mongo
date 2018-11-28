package factories

import "github.com/google/go-cloud/wire"

var FactoriesSet = wire.NewSet(
	NewUserFactory,
	NewAccessFactory,
)
