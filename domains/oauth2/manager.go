package oauth2

import (
	"context"

	"github.com/marcusyip/golang-wire-mongo/models"
	repos "github.com/marcusyip/golang-wire-mongo/repositories"
	"github.com/sirupsen/logrus"
)

var (
	providerTypeGoogle   = "google"
	providerTypeFacebook = "facebook"
)

type Provider interface {
	Authorize(ctx context.Context, params *AuthorizeParams) (*models.User, error)
	AuthCodeURL(state State) (string, error)
}

type State struct {
	RedirectURL string `json:"redirect_url"`
}

type Manager struct {
	logger   *logrus.Logger
	userRepo *repos.UserRepository
	fb       *FacebookProvider
}

func (mgr *Manager) AuthCodeURL(providerType string, state State) (string, error) {
	provider := mgr.getProvider(providerType)
	return provider.AuthCodeURL(state)
}

type AuthorizeParams struct {
	AuthCode    string `json:"auth_code"`
	RedirectURI string `json:"redirect_uri"`
}

func (mgr *Manager) Authorize(ctx context.Context, providerType string, params *AuthorizeParams) (*models.User, error) {
	log := mgr.getLogger("Authorize")
	provider := mgr.getProvider(providerType)
	user, err := provider.Authorize(ctx, params)
	if err != nil {
		log.WithError(err).Error("Failed to get provider")
		return nil, err
	}
	return user, nil
}

func (mgr *Manager) getProvider(providerType string) Provider {
	switch providerType {
	//case providerTypeGoogle:
	//	return newGoogleOAuth2()
	case providerTypeFacebook:
		return mgr.fb
	default:
		return nil
	}
}

func (mgr *Manager) getLogger(method string) *logrus.Entry {
	return mgr.logger.WithFields(logrus.Fields{"domain": "oauth2", "method": method})
}

func NewManager(
	logger *logrus.Logger,
	userRepo *repos.UserRepository,
	fb *FacebookProvider,
) *Manager {
	return &Manager{
		logger:   logger,
		userRepo: userRepo,
		fb:       fb,
	}
}
