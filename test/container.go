package test

import (
	"github.com/gin-gonic/gin"
	"github.com/marcusyip/golang-wire-mongo/app"
	"github.com/marcusyip/golang-wire-mongo/db/migrates"
	repos "github.com/marcusyip/golang-wire-mongo/repositories"
	"github.com/marcusyip/golang-wire-mongo/test/factories"
	"github.com/marcusyip/golang-wire-mongo/web/api"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Container struct {
	App       *app.App
	ApiRouter *api.Router
	Engine    *gin.Engine
	DbClient  *mongo.Client
	Db        *mongo.Database
	DbMigrate *migrates.MigrationJob

	UserRepo   *repos.UserRepository
	AccessRepo *repos.AccessRepository

	UserFactory *factories.UserFactory
}

func NewContainer(
	app *app.App,
	apiRouter *api.Router,
	engine *gin.Engine,
	dbClient *mongo.Client,
	db *mongo.Database,
	dbMigrate *migrates.MigrationJob,

	userRepo *repos.UserRepository,
	accessRepo *repos.AccessRepository,

	userFactory *factories.UserFactory,
) *Container {
	return &Container{
		App:       app,
		ApiRouter: apiRouter,
		Engine:    engine,
		DbClient:  dbClient,
		Db:        db,
		DbMigrate: dbMigrate,

		UserRepo:   userRepo,
		AccessRepo: accessRepo,

		UserFactory: userFactory,
	}
}
