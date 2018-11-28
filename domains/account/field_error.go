package account

import "strings"

type FieldErrors []error

func (es FieldErrors) Error() string {
	var errs []string
	for _, e := range es {
		errs = append(errs, e.Error())
	}
	return strings.Join(errs, ";")
}

func (es FieldErrors) Errors() []error {
	return es
}

type FieldError struct {
	Name string
	Err  error
}

func (err FieldError) Error() string {
	return err.Name + ": " + err.Err.Error()
}
