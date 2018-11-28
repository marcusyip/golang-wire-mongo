package e2e

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/marcusyip/golang-wire-mongo/config"
	"github.com/marcusyip/golang-wire-mongo/models"
	"github.com/marcusyip/golang-wire-mongo/test"
	"github.com/mongodb/mongo-go-driver/bson"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var ctn *test.Container

func TestWeb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Web Suite")
}

var _ = BeforeSuite(func() {
	viper.AddConfigPath("../../config")
	conf := config.ProvideConfig()
	var err error
	ctn, err = test.BuildContainer(conf)
	if err != nil {
		panic(err)
	}
	err = ctn.DbClient.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	_ = ctn.Db.Drop(context.Background())
	err = ctn.DbMigrate.Run()
	if err != nil {
		panic(err)
	}
	ctn.ApiRouter.With(ctn.Engine)
})

var _ = AfterSuite(func() {
	err := ctn.Db.Drop(context.Background())
	if err != nil {
		panic(err)
	}
	err = ctn.DbClient.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}
})

func AddAccessToken(req *http.Request, user *models.User) {
	access, err := ctn.AccessRepo.FindOne(context.Background(), bson.D{{"user_id", user.ID}})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", access)
	req.Header.Add("X-Access-Token", access.AccessToken)
}
