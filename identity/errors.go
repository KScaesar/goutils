package identity

import (
	"net/http"

	"github.com/KScaesar/goutils/errors"
)

var (
	ErrAuthentication = errors.New(
		2001,
		http.StatusUnauthorized,
		"invalid authentication",
	)

	ErrAuthorization = errors.New(
		2002,
		http.StatusForbidden,
		"invalid authorization",
	)
)
