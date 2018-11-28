package oauth2

import (
	"context"
	"encoding/json"

	"github.com/marcusyip/golang-wire-mongo/config"
	huacct "github.com/marcusyip/golang-wire-mongo/domains/account"
	"github.com/marcusyip/golang-wire-mongo/models"
	repos "github.com/marcusyip/golang-wire-mongo/repositories"
	fb "github.com/huandu/facebook"
	"github.com/mongodb/mongo-go-driver/bson"
	"golang.org/x/oauth2"
	oauth2fb "golang.org/x/oauth2/facebook"
)

// https://www.facebook.com/v2.9/dialog/oauth?
//  client_id=222446594917019
//    &redirect_uri=http%3A%2F%2Flocalhost%3A3600%2Foauth%2Ffacebook%2Fredirect
//    &scope=public_profile,email
//    &response_type=code

type FacebookProvider struct {
	conf     *config.Config
	acctMgr  *huacct.Manager
	userRepo *repos.UserRepository
}

// AuthCodeURL returns oauth google URL
func (p *FacebookProvider) AuthCodeURL(state State) (string, error) {
	oauth2Config := p.getOAuth2Config()
	sState, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	opts := []oauth2.AuthCodeOption{
		// oauth2.SetAuthURLParam("hd", g.googleConf.HostDomain),
		// oauth2.SetAuthURLParam("prompt", "select_account"),
	}
	return oauth2Config.AuthCodeURL(string(sState), opts...), nil
}

func (p *FacebookProvider) Authorize(ctx context.Context, params *AuthorizeParams) (*models.User, error) {
	oauth2Config := p.getOAuth2Config()
	token, err := oauth2Config.Exchange(oauth2.NoContext, params.AuthCode)
	if err != nil {
		return nil, err
	}
	res, err := fb.Get("/me", fb.Params{
		"access_token": token.AccessToken,
		"fields":       "email,name,picture.type(large)",
	})
	if err != nil {
		return nil, err
	}
	var profile FacebookProfile
	res.Decode(&profile)
	return p.register(ctx, &profile, token)
}

func (p *FacebookProvider) getOAuth2Config() *oauth2.Config {
	fbConf := p.conf.Facebook
	return &oauth2.Config{
		ClientID:     fbConf.ClientID,
		ClientSecret: fbConf.ClientSecret,
		Scopes:       []string{"public_profile", "email"},
		RedirectURL:  fbConf.RedirectURI,
		Endpoint:     oauth2fb.Endpoint,
	}
}

func (p *FacebookProvider) register(ctx context.Context, profile *FacebookProfile, token *oauth2.Token) (*models.User, error) {
	filter := bson.D{
		{"facebook_id", profile.ID},
	}
	user, err := p.userRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	if user == nil {
		params := p.createAccountParams(profile)
		params.OAuth2Token = token
		user, err = p.acctMgr.Create(ctx, params)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (p *FacebookProvider) createAccountParams(profile *FacebookProfile) *huacct.CreateParams {
	return &huacct.CreateParams{
		Email:         profile.Email,
		DisplayName:   profile.Name,
		EmailVerified: true,
		AvatarURL:     &(profile.Picture.Data.URL),
		Provider:      string(providerTypeFacebook),
		FacebookID:    &(profile.ID),
	}
}

func NewFacebookProvider(
	conf *config.Config,
	acctMgr *huacct.Manager,
	userRepo *repos.UserRepository,
) *FacebookProvider {
	return &FacebookProvider{
		conf:     conf,
		acctMgr:  acctMgr,
		userRepo: userRepo,
	}
}

type FacebookProfile struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture struct {
		Data struct {
			IsSilhouette bool   `json:"is_silhouette"`
			URL          string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}
