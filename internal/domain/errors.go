package domain

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(message string) *Error {
	return &Error{Message: message}
}

var (
	ErrInvalidDocumentNumber = NewError("document number is required")
	ErrAccountAlreadyExists  = NewError("account already exists")
)
