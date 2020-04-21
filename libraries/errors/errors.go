package errors

type Error interface {
	Text() string
	ErrorCode() int
	StackTrace() string
}

// New returns a new error.
func New(err error) Error {
	return InternalError(err)
}
