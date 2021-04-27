package errorX

import (
	"net/http"
)

var (
	ErrSystemFailed = New(
		10001,
		http.StatusInternalServerError,
		"system error",
	)
	ErrInvalidParams = New(
		10002,
		http.StatusBadRequest,
		"invalid params",
	)

	ErrAuthorizeFailed = New(
		20001,
		http.StatusUnauthorized,
		"auth failed",
	)

	ErrTimeout = New(
		30001,
		http.StatusInternalServerError,
		"timeout",
	)
)
