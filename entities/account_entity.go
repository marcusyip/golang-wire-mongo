package entities

import (
	"github.com/marcusyip/golang-wire-mongo/models"
)

type Account struct {
	EntityObject
	ID    string `json:"id"`
	Email string `json:"email"`

	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Provider    string `json:"provider"`
	FacebookID  string `json:"facebook_id"`

	EmailVerified bool `json:"email_verified"`
}

type AccountEntity struct {
}

func (e *AccountEntity) New(m *models.User) *Account {
	return &Account{
		EntityObject:  EntityObject{"account"},
		ID:            m.ID.Hex(),
		Email:         m.Email,
		Username:      m.Username,
		DisplayName:   m.DisplayName,
		AvatarURL:     m.AvatarURL,
		Provider:      m.Provider,
		FacebookID:    m.FacebookID,
		EmailVerified: m.EmailVerified,
	}
}

func NewAccountEntity() *AccountEntity {
	return &AccountEntity{}
}
