package app

import (
	"context"

	"github.com/marcusyip/golang-wire-mongo/config"
	"github.com/marcusyip/golang-wire-mongo/web"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type App struct {
	conf      *config.Config
	webServer *web.Server
	dbClient  *mongo.Client
}

func (app *App) Start() {
	app.dbClient.Connect(context.Background())
	app.webServer.Start()
}

func NewApp(
	conf *config.Config,
	webServer *web.Server,
	dbClient *mongo.Client,
) *App {
	return &App{
		conf:      conf,
		webServer: webServer,
		dbClient:  dbClient,
	}
}
