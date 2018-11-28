package entities

import "github.com/google/go-cloud/wire"

var WireSet = wire.NewSet(
	NewErrorEntity,
	NewAccountEntity,
	NewAccessEntity,
)
