package identity

import (
	"net/http"

	"github.com/Min-Feng/goutils/errorY"
)

var (
	ErrAuthentication = errorY.New(
		2001,
		http.StatusUnauthorized,
		"invalid authentication",
	)

	ErrAuthorization = errorY.New(
		2002,
		http.StatusForbidden,
		"invalid authorization",
	)
)
