//+build wireinject

package cmd

import (
	"github.com/google/go-cloud/wire"
	huapp "github.com/marcusyip/golang-wire-mongo/app"
	"github.com/marcusyip/golang-wire-mongo/config"
	"github.com/marcusyip/golang-wire-mongo/core/logger"
	huaccess "github.com/marcusyip/golang-wire-mongo/domains/access"
	huacct "github.com/marcusyip/golang-wire-mongo/domains/account"
	huauth "github.com/marcusyip/golang-wire-mongo/domains/auth"
	huoauth "github.com/marcusyip/golang-wire-mongo/domains/oauth2"
	ents "github.com/marcusyip/golang-wire-mongo/entities"
	repos "github.com/marcusyip/golang-wire-mongo/repositories"
	"github.com/marcusyip/golang-wire-mongo/web"
)

func BuildApp(conf *config.Config) (*huapp.App, error) {
	wire.Build(
		logger.New,
		repos.WireSet,
		web.WireSet,
		ents.WireSet,
		huauth.AuthSet,
		huacct.AccountSet,
		huaccess.AccessSet,
		huoauth.OAuth2Set,
		huapp.AppSet,
	)
	return nil, nil
}
