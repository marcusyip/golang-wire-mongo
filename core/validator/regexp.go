package validator

import "regexp"

const (
	Name string = "^[a-zA-Z0-9_]+$"
)

var (
	rxName = regexp.MustCompile(Name)
)
