package web

import (
	"github.com/google/go-cloud/wire"
	"github.com/marcusyip/golang-wire-mongo/web/api"
	mids "github.com/marcusyip/golang-wire-mongo/web/middlewares"
)

var WireSet = wire.NewSet(
	NewServer,
	mids.WireSet,
	api.WireSet,
)
