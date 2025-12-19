package domain

type ErrorKind int

const (
	KindValidation ErrorKind = iota
	KindNotFound
)

type Error struct {
	Kind    ErrorKind
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrInvalidDocumentNumber = &Error{KindValidation, "document number is required"}
	ErrAccountAlreadyExists  = &Error{KindValidation, "account already exists"}
	ErrAccountNotFound       = &Error{KindNotFound, "account was not found"}
	ErrInvalidOperationType  = &Error{KindValidation, "invalid operation type"}
	ErrInvalidAmount         = &Error{KindValidation, "amount must be greater than zero"}
)
