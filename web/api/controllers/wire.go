package controllers

import "github.com/google/go-cloud/wire"

var WireSet = wire.NewSet(
	NewRegisterController,
	NewAccountController,
	NewAccessController,
	NewOauthFacebookController,
)
