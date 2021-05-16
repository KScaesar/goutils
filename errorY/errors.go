package errorY

import (
	"net/http"
)

var (
	ErrSystem = New(
		1001,
		http.StatusInternalServerError,
		"system failed",
	)
	ErrInvalidParams = New(
		1002,
		http.StatusBadRequest,
		"invalid params",
	)
)
