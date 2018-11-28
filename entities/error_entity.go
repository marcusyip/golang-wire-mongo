package entities

type Error struct {
	EntityObject
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorEntity struct {
}

func (e *ErrorEntity) New(code int, message string) *Error {
	return &Error{
		EntityObject: EntityObject{"error"},
		Code:         code,
		Message:      message,
	}
}

func NewErrorEntity() *ErrorEntity {
	return &ErrorEntity{}
}
