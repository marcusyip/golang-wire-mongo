package oauth2

import (
	"golang.org/x/oauth2"
	"project.scmp.tech/technology/newsroom-system/accounts/models"
)

type Provider interface {
	Authorize(authorizeParams *AuthorizeParams) (*models.User, error)
	AuthCodeURL(state State) (string, error)
}

type ProviderType string

type exchangeToken func(authCode string) (*oauth2.Token, error)

var (
	GoogleProvider ProviderType = "google"
)

type ProvidersManager interface {
	GetProvider(p ProviderType) Provider
}

type ProvidersManagerImpl struct{}

func (mgr *ProvidersManagerImpl) GetProvider(p ProviderType) Provider {
	switch p {
	case GoogleProvider:
		return newGoogleOAuth2()
	default:
		return nil
	}
}

func GetProvidersManager() ProvidersManager {
	return &ProvidersManagerImpl{}
}
