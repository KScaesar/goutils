package errorX

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapf(t *testing.T) {
	contextDescription := "open file failed"
	thirdPartyErrMsg := "file is not exist"
	topDefinedErr := New(10001, http.StatusInternalServerError, "system internal error")

	openFile := func() error {
		return errors.New(thirdPartyErrMsg)
	}
	openErr := openFile()
	wrapErr := Wrapf(topDefinedErr, "%v: %v", contextDescription, openErr)

	actualErrMsg := wrapErr.Error()
	expectedErrMsg := fmt.Sprintf("%v: %v: %v", contextDescription, openErr, topDefinedErr)
	assert.Equal(t, expectedErrMsg, actualErrMsg)
}

func TestIs(t *testing.T) {
	topDefinedErr := New(10001, http.StatusInternalServerError, "system internal error")

	repo := func() error {
		infraErr := errors.New("connect fail")
		return Wrap(topDefinedErr, infraErr.Error())
	}
	useCase := func() error {
		repoErr := repo()
		return WithMsg(repoErr, "save entity failed")
	}
	finalErr := useCase()

	actualResult := Is(finalErr, topDefinedErr)
	assert.True(t, actualResult)
}

func TestCause(t *testing.T) {
	topDefinedErr := New(10001, http.StatusInternalServerError, "system internal error")

	repo := func() error {
		infraErr := errors.New("connect fail")
		return Wrap(topDefinedErr, infraErr.Error())
	}
	useCase := func() error {
		repoErr := repo()
		return WithMsg(repoErr, "save entity failed")
	}
	api := func() error {
		ucErr := useCase()
		return Wrap(ucErr, "key=XXX")
	}
	finalErr := api()

	var result bool
	switch Cause(finalErr) {
	case topDefinedErr:
		result = true
	default:
		result = false
	}
	assert.True(t, result)
}
