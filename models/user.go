package models

import "golang.org/x/oauth2"

type User struct {
	Model       `bson:",inline"`
	Email       string        `bson:"email"`
	Username    string        `bson:"username"`
	DisplayName string        `bson:"display_name"`
	AvatarURL   string        `bson:"avatar_url"`
	Provider    string        `bson:"provider"`
	OAuth2Token *oauth2.Token `bson:"oauth2_token,omitempty"`
	FacebookID  string        `bson:"facebook_id"`

	HashedPassword string `bson:"hashed_password"`
	EmailVerified  bool   `bson:"email_verified"`
}
