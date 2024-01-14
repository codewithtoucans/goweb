package errors

type publicError struct {
	err error
	msg string
}

func NewPublicError(err error, msg string) error {
	return publicError{err: err, msg: msg}
}

func (e publicError) Error() string {
	return e.err.Error()
}

func (e publicError) Unwrap() error {
	return e.err
}

func (e publicError) Public() string {
	return e.msg
}
