package errors

import (
	"net/http"
)

var (
	ErrSystem = New(
		1001,
		http.StatusInternalServerError,
		"system failed",
	)
	ErrCoding = New(
		1002,
		http.StatusInternalServerError,
		"backend developer does not correctly understand how the function is used",
	)
	ErrInvalidParams = New(
		1003,
		http.StatusBadRequest,
		"invalid parameter",
	)
	ErrNotFound = New(
		10004,
		http.StatusNotFound,
		"not found",
	)
)
