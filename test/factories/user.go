package factories

import (
	"context"

	huauth "github.com/marcusyip/golang-wire-mongo/domains/auth"
	"github.com/marcusyip/golang-wire-mongo/models"
	repos "github.com/marcusyip/golang-wire-mongo/repositories"
	"github.com/marcusyip/golang-wire-mongo/test/options"
	"github.com/icrowley/fake"
)

var UserWithAccess = options.Field{"with_access", true}

type UserFactory struct {
	accessFactory *AccessFactory
	authMgr       *huauth.Manager
	userRepo      *repos.UserRepository
}

func (f *UserFactory) Create(fields ...options.Field) *models.User {
	fMap := options.Fields(fields).Map()

	user := &models.User{
		Model:       *models.NewModel(),
		Username:    fMap.String("username", fake.UserName()),
		Email:       fMap.String("email", fake.EmailAddress()),
		DisplayName: fMap.String("display_name", fake.FullName()),
	}
	if fMap.Has("password") {
		hashedPassword, err := f.authMgr.HashPassword(
			fMap.String("password", fake.Password(8, 15, true, true, true)),
		)
		if err != nil {
			panic(err)
		}
		user.HashedPassword = hashedPassword
	}
	err := f.userRepo.Create(context.Background(), user)
	if err != nil {
		panic(err)
	}
	if fMap.Has("with_access") {
		f.accessFactory.Create(options.Field{"user_id", user.ID})
	}

	return user
}

func NewUserFactory(
	accessFactory *AccessFactory,
	authMgr *huauth.Manager,
	userRepo *repos.UserRepository,
) *UserFactory {
	return &UserFactory{
		accessFactory: accessFactory,
		authMgr:       authMgr,
		userRepo:      userRepo,
	}
}
