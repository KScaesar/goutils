package errorY

import (
	"net/http"
)

var (
	ErrSystem = New(
		10001,
		http.StatusInternalServerError,
		"system failed",
	)
	ErrInvalidParams = New(
		10002,
		http.StatusBadRequest,
		"invalid params",
	)
)
