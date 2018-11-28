package app

import "github.com/google/go-cloud/wire"

var AppSet = wire.NewSet(
	NewGinEngine,
	NewApp,
)
