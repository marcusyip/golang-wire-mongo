package account

import "github.com/google/go-cloud/wire"

var AccountSet = wire.NewSet(NewManager)
