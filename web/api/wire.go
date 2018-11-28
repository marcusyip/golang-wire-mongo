package api

import (
	ctrls "github.com/marcusyip/golang-wire-mongo/web/api/controllers"
	"github.com/google/go-cloud/wire"
)

var WireSet = wire.NewSet(
	ctrls.WireSet,
	NewRouter,
)
