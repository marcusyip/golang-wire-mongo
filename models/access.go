package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type AccessProvider string

var (
	AccessProviderGuest AccessProvider = "guest"
	// Google oauth
	AccessProviderGoogle   AccessProvider = "google"
	AccessProviderFacebook AccessProvider = "facebook"
	// Admin API
	AccessProviderAdmin AccessProvider = "admin"
)

type Access struct {
	Model        `bson:",inline"`
	AccessToken  string            `bson:"access_token"`
	RefreshToken string            `bson:"refresh_token"`
	UserID       objectid.ObjectID `bson:"user_id,omitempty"`
	ExpireAt     time.Time         `bson:"expire_at"`
	Provider     AccessProvider    `bson:"provider"`

	User *User `bson:"-"`
}
