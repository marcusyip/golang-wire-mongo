package validator

import "github.com/asaskevich/govalidator"

const NameCharacters = "a-z A-Z _"

func IsName(str string) bool {
	return rxName.MatchString(str)
}

func init() {
	govalidator.TagMap["modelname"] = govalidator.Validator(IsName)
}
