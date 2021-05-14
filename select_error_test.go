package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/goutils/errorY"
	"github.com/Min-Feng/goutils/logY"
)

func TestSelectError(t *testing.T) {
	logY.TestingMode()

	var majorErr error
	err := repo()
	if err != nil {
		if rollbackErr := rollback(); rollbackErr != nil {
			majorErr = SelectError(nil, err, rollbackErr)
		}
	}

	expectedMsg := "[select error]: not found"
	assert.Equal(t, expectedMsg, majorErr.Error())
}

func repo() error {
	return errorY.NewPkgError("not found")
}

func rollback() error {
	return errorY.NewPkgError("rollback: connect failed")
}
