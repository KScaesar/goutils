package errors

type customError interface {
	rootCause()
	Code() int
	HttpStatus() int
}
