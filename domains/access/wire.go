package access

import "github.com/google/go-cloud/wire"

var AccessSet = wire.NewSet(NewManager)
