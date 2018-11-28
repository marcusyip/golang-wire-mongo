package oauth2

import (
	"encoding/json"
	"errors"

	"github.com/marcusyip/golang-wire-mongo/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	people "google.golang.org/api/people/v1"
	mgo "gopkg.in/mgo.v2"
)

type GoogleProvider struct {
}

// AuthCodeURL returns oauth google URL
func (p *GoogleProvider) AuthCodeURL(state State) (string, error) {
	oauth2Config := g.getOAuth2Config()
	sState, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("hd", g.googleConf.HostDomain),
		oauth2.SetAuthURLParam("prompt", "select_account"),
	}
	return oauth2Config.AuthCodeURL(string(sState), opts...), nil
}

// Authorize returns new created or updated user
func (p *GoogleProvider) Authorize(params *AuthorizeParams) (*models.User, error) {
	return g.doAuthorize(params, g.exchangeToken, g.getPersonCall)
}

func (p *GoogleProvider) doAuthorize(params *AuthorizeParams, exchangeToken exchangeToken,
	getPersonCall googleGetPersonCall) (*models.User, error) {

	log := g.getLogger("Authorize")
	token, err := exchangeToken(params.AuthCode)
	if err != nil {
		log.WithError(err).Info("Failed to exchange token")
		return nil, err
	}

	person, err := getPersonCall(token)
	if err != nil {
		log.WithError(err).Error("Failed to get me from google people API")
		return nil, err
	}

	if len(person.EmailAddresses) == 0 {
		return nil, errors.New("empty email addresses")
	}

	err = g.acctMgr.ValidateEmail(person.EmailAddresses[0].Value)
	if err != nil {
		log.WithError(err).Info("Invalid email")
		return nil, err
	}

	googleID := g.getGoogleID(person)
	if len(googleID) == 0 {
		log.WithError(err).WithField("person.metadata.sources", person.Metadata.Sources).Warn("Invalid google ID")
		return nil, errors.New("Cannot find google ID from people API")
	}

	m, err := g.userRepo.FindByGoogleID(googleID)
	if err != nil && err != mgo.ErrNotFound {
		log.WithError(err).WithField("google_id", googleID).Error("Failed to find user by google ID")
		return nil, err
	}

	var user *models.User
	if m == nil {
		user, err = g.createUser(token, person)
		if err != nil {
			log.WithError(err).WithField("google_id", googleID).Error("Failed to create user with google Person")
			return nil, err
		}
	} else {
		user = m.(*models.User)
		err = g.updateUser(user, token, person)
		if err != nil {
			log.WithError(err).WithField("user_id", user.ID).Warn("Failed to update user with google Person")
			// Although error, but no need to return error since user can still use it
		}
	}
	return user, nil
}

func (p *GoogleProvider) getGoogleID(person *people.Person) string {
	for _, source := range person.Metadata.Sources {
		if source.Type == "PROFILE" || source.Type == "DOMAIN_PROFILE" {
			return source.Id
		}
	}
	return ""
}

func (p *GoogleProvider) getOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     g.googleConf.ClientID,
		ClientSecret: g.googleConf.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/user.emails.read",
			"https://www.googleapis.com/auth/calendar",
			"https://www.googleapis.com/auth/contacts.readonly",
		},
		RedirectURL: g.googleConf.RedirectURL,
		Endpoint:    google.Endpoint,
	}
}

func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{}
}
