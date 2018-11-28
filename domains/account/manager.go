package account

import (
	"context"
	"errors"

	"github.com/marcusyip/golang-wire-mongo/models"
	repos "github.com/marcusyip/golang-wire-mongo/repositories"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/oauth2"
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrUsernameAlreadyExists = errors.New("username already exists")

type Manager struct {
	userRepo *repos.UserRepository
}

func (mgr *Manager) FindOne(ctx context.Context, filter bson.D) (*models.User, error) {
	user, err := mgr.userRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (mgr *Manager) FindAll() {
	return
}

type CreateParams struct {
	Email          string        `valid:"email,required"`
	Username       *string       `valid:"length(8)"`
	HashedPassword *string       `valid:"-"`
	DisplayName    string        `valid:"length(3)"`
	AvatarURL      *string       `valid:"-"`
	Provider       string        `valid:"-"`
	OAuth2Token    *oauth2.Token `valid:"-"`
	FacebookID     *string       `valid:"-"`

	EmailVerified bool
}

func (mgr *Manager) Create(ctx context.Context, params *CreateParams) (*models.User, error) {
	user := &models.User{
		Model:         *models.NewModel(),
		Email:         params.Email,
		DisplayName:   params.DisplayName,
		Provider:      params.Provider,
		EmailVerified: params.EmailVerified,
	}
	if params.HashedPassword != nil {
		user.HashedPassword = *params.HashedPassword
	}
	if params.Username != nil {
		user.Username = *params.Username
	}
	if params.AvatarURL != nil {
		user.AvatarURL = *params.AvatarURL
	}
	if params.OAuth2Token != nil {
		user.OAuth2Token = params.OAuth2Token
	}
	err := mgr.userRepo.Create(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

type UpdateParams struct {
	HashedPassword *string `valid:"-"`
	DisplayName    *string `valid:"length(3)"`
	AvatarURL      *string `valid:"-"`
}

func (mgr *Manager) Update(ctx context.Context, ctxUser *models.User, params *UpdateParams) error {
	update := bson.D{}
	if params.DisplayName != nil {
		update = append(update, bson.E{"display_name", *params.DisplayName})
		defer func() {
			ctxUser.DisplayName = *params.DisplayName
		}()
	}
	_, err := mgr.userRepo.UpdateByID(ctx, ctxUser.ID, update)
	return err
}

type ValidateParams struct {
	Email    *string
	Username *string
}

func (mgr *Manager) Validate(ctx context.Context, params *ValidateParams) (result bool, err error) {
	var errs FieldErrors
	user, err := mgr.userRepo.FindOne(ctx, bson.D{{"email", *params.Email}})
	if err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}
	if user != nil {
		errs = append(errs, FieldError{"email", ErrEmailAlreadyExists})
	}
	user, err = mgr.userRepo.FindOne(ctx, bson.D{{"username", *params.Username}})
	if err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}
	if user != nil {
		errs = append(errs, FieldError{"username", ErrUsernameAlreadyExists})
	}
	result = len(errs) > 0
	if result {
		return false, errs
	}
	return true, nil
}

func NewManager(
	userRepo *repos.UserRepository,
) *Manager {
	return &Manager{
		userRepo: userRepo,
	}
}
