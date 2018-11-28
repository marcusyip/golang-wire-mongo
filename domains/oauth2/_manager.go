package oauth2

//go:generate mockgen -destination=../../mocks/mock_managers/mock_oauth2/mock_manager.go -package=mock_oauth2 -source=oauth2_manager.go

import (
	"github.com/sirupsen/logrus"
	"project.scmp.tech/technology/newsroom-system/accounts/core/logger"
	"project.scmp.tech/technology/newsroom-system/accounts/models"
	repos "project.scmp.tech/technology/newsroom-system/accounts/repositories"
)

type OAuth2Manager interface {
	Authorize(providerType string, params *AuthorizeParams) (*models.User, error)
	AuthCodeURL(providerType string, state State) (string, error)
}

type OAuth2ManagerImpl struct {
	logger       *logrus.Logger
	userRepo     repos.UserRepository
	providersMgr ProvidersManager
}

type State struct {
	RedirectURL string `json:"redirect_url"`
}

type AuthorizeParams struct {
	AuthCode    string `json:"auth_code"`
	RedirectURI string `json:"redirect_uri"`
}

func (mgr *OAuth2ManagerImpl) Authorize(providerType string, params *AuthorizeParams) (user *models.User, err error) {
	log := mgr.getLogger("Authorize")
	provider := mgr.providersMgr.GetProvider(ProviderType(providerType))
	user, err = provider.Authorize(params)
	if err != nil {
		log.WithError(err).Error("Failed to get provider")
		return nil, err
	}
	return user, nil
}

func (mgr *OAuth2ManagerImpl) AuthCodeURL(providerType string, state State) (string, error) {
	provider := mgr.providersMgr.GetProvider(ProviderType(providerType))
	return provider.AuthCodeURL(state)
}

func (mgr *OAuth2ManagerImpl) getLogger(method string) *logrus.Entry {
	return mgr.logger.WithFields(logrus.Fields{"manager": "oauth2", "method": method})

}

func GetOAuth2Manager() OAuth2Manager {
	return &OAuth2ManagerImpl{
		logger:       logger.Get(),
		userRepo:     repos.GetUserRepository(),
		providersMgr: GetProvidersManager(),
	}
}
