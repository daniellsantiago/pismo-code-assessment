package domain

var (
	ErrInvalidDocumentNumber = NewError("document number is required")
	ErrAccountAlreadyExists  = NewError("account already exists")
	ErrAccountNotFound       = NewError("account was not found")
)

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(message string) *Error {
	return &Error{Message: message}
}
