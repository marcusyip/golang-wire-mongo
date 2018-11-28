package entities

import (
	"time"

	"github.com/marcusyip/golang-wire-mongo/models"
)

type Access struct {
	EntityObject
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	AccessToken string    `json:"access_token"`
	ExpireAt    time.Time `json:"expire_at"`
	Provider    string    `json:"provider"`
	CreatedAt   time.Time `json:"created_at"`
}

type AccessEntity struct {
}

func (e *AccessEntity) New(m *models.Access) *Access {
	return &Access{
		EntityObject: EntityObject{"access"},
		ID:           m.ID.Hex(),
		UserID:       m.UserID.Hex(),
		AccessToken:  m.AccessToken,
		ExpireAt:     m.ExpireAt,
		Provider:     string(m.Provider),
		CreatedAt:    m.CreatedAt,
	}
}

func NewAccessEntity() *AccessEntity {
	return &AccessEntity{}
}
