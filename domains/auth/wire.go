package auth

import "github.com/google/go-cloud/wire"

var AuthSet = wire.NewSet(NewManager)
