package auth

import (
	"context"
	"errors"

	huacct "github.com/marcusyip/golang-wire-mongo/domains/account"
	"github.com/marcusyip/golang-wire-mongo/models"
	repos "github.com/marcusyip/golang-wire-mongo/repositories"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredential = errors.New("invalid credential")

type Manager struct {
	acctMgr  *huacct.Manager
	userRepo *repos.UserRepository
}

func (mgr *Manager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

type AuthorizeParams struct {
	Username string
	Password string
}

func (mgr *Manager) Authorize(ctx context.Context, params *AuthorizeParams) (*models.User, error) {
	user, err := mgr.userRepo.FindOne(ctx, bson.D{{"username", params.Username}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrInvalidCredential
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(params.Password))
	if err != nil {
		return nil, ErrInvalidCredential
	}
	return user, nil
}

func NewManager(
	acctMgr *huacct.Manager,
	userRepo *repos.UserRepository,
) *Manager {
	return &Manager{
		acctMgr:  acctMgr,
		userRepo: userRepo,
	}
}
