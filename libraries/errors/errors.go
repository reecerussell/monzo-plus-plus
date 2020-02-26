package errors

type Error interface {
	Text() string
	ErrorCode() int
	StackTrace() string
}
