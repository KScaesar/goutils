package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode(t *testing.T) {
	definedCode := 9478
	definedErr := New(definedCode, http.StatusBadGateway, "not match ip")
	fooErr := Wrap(definedErr, "foo")
	barErr := WrapMessage(fooErr, "bar")

	assert.Equal(t, definedCode, Code(barErr))
	assert.Equal(t, http.StatusBadGateway, HttpStatus(barErr))
}

func TestSimpleInfo(t *testing.T) {
	definedCode := 9478
	definedErr := New(definedCode, http.StatusBadGateway, "not match ip")
	repoErr := Wrap(definedErr, "repo save failed")
	fooErr := Wrap(repoErr, "foo failed")

	actualMsg := SimpleInfo(fooErr)
	expectedMsg := "repo save failed: not match ip"
	assert.Equal(t, expectedMsg, actualMsg)
}
