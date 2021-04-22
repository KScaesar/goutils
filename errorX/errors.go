package errorX

import (
	"net/http"
)

var (
	ErrSystemFailed = New(
		10001,
		http.StatusInternalServerError,
		"system internal error",
	)
	ErrInvalidParams = New(
		10002,
		http.StatusBadRequest,
		"invalid params",
	)
	ErrNotFound = New(
		10003,
		http.StatusNotFound,
		"not found",
	)

	ErrAuthorizeFailed = New(
		20001,
		http.StatusUnauthorized,
		"auth failed",
	)
	ErrTokenParseFailed = New(
		20002,
		http.StatusUnauthorized,
		"token parse failed",
	)
	ErrTokenExpired = New(
		20003,
		http.StatusUnauthorized,
		"token expired",
	)

	ErrProviderTimeout = New(
		30001,
		http.StatusInternalServerError,
		"timeout when connect to provider",
	)
	ErrProviderNotWork = New(
		30002,
		http.StatusInternalServerError,
		"provider service failed",
	)
)
