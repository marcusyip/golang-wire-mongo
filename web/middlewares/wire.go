package middlewares

import "github.com/google/go-cloud/wire"

var WireSet = wire.NewSet(
	NewTokenAuthenticator,
	NewMongoSession,
)
