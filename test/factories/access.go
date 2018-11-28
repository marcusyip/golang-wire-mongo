package factories

import (
	"context"
	"time"

	"github.com/marcusyip/golang-wire-mongo/core/rand"
	"github.com/marcusyip/golang-wire-mongo/models"
	repos "github.com/marcusyip/golang-wire-mongo/repositories"
	"github.com/marcusyip/golang-wire-mongo/test/options"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type AccessFactory struct {
	accessRepo *repos.AccessRepository
}

func (f *AccessFactory) Create(fields ...options.Field) *models.Access {
	fMap := options.Fields(fields).Map()
	access := &models.Access{
		Model:        *models.NewModel(),
		AccessToken:  fMap.String("access_token", rand.RandString(48)),
		RefreshToken: fMap.String("refresh_token", rand.RandString(48)),
		UserID:       fMap.ObjectID("user_id", objectid.New()),
		ExpireAt:     fMap.Time("expire_at", time.Now().AddDate(0, 1, 0)),
	}
	err := f.accessRepo.Create(context.Background(), access)
	if err != nil {
		panic(err)
	}
	return access
}

func NewAccessFactory(
	accessRepo *repos.AccessRepository,
) *AccessFactory {
	return &AccessFactory{
		accessRepo: accessRepo,
	}
}
