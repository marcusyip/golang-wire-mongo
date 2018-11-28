package errors

import "github.com/marcusyip/golang-wire-mongo/entities"

var ent entities.ErrorEntity

var InternalServerError = ent.New(10000, "Internal server error")
var MissingAccessToken = ent.New(10100, "Missing access token")
var UnauthorizedError = ent.New(10101, "Unauthorized")
var InvalidCredential = ent.New(10101, "Invalid credential")
var UsernameAlreadyExists = ent.New(10102, "Username already exists")
var EmailAlreadyExists = ent.New(10103, "Email already exists")
