package oauth2

import "github.com/google/go-cloud/wire"

var OAuth2Set = wire.NewSet(
	NewManager,
	NewFacebookProvider,
)
