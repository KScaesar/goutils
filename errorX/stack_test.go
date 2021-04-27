package errorX

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:generate go test -trimpath -run=TestStacks_CountOK -v github.com/Min-Feng/goutils/errorX
func TestStacks_CountOK(t *testing.T) {
	topDefinedErr := New(10001, http.StatusInternalServerError, "system internal error")

	repo := func() error {
		infraErr := errors.New("connect fail")
		return Wrap(topDefinedErr, infraErr.Error())
	}
	rollback := func(error) error {
		return errors.New("data center bang")
	}
	useCase := func() error {
		repoErr := repo()
		if repoErr != nil {
			fixErr := rollback(repoErr)
			if fixErr != nil {
				// 遇到第二個 error 應該用 Wrap 再次包裝 保留 stack 訊息
				// 並把 次要錯誤(fixErr)的文字訊息 加入到 主要錯誤(repoErr)
				// 若想保留 fixErr stack, 可以考慮在此處加上 log
				return Wrapf(repoErr, "update entity failed: %v", fixErr)
			}

			// 如果只是單一錯誤往上傳遞 只用 WithMsg 增加文字訊息就好
			return WithMsg(repoErr, "update entity failed")
		}
		return nil
	}
	finalErr := useCase()
	stacks := Stacks(finalErr)

	for _, stack := range stacks {
		for _, frame := range stack {
			fmt.Println(frame)
		}
		fmt.Println()
	}

	var actualStackCount int
	for _, stack := range stacks {
		actualStackCount += len(stack)
	}
	expectedStackCount := 9
	assert.Equal(t, expectedStackCount, actualStackCount)
	assert.Equal(t, topDefinedErr, Cause(finalErr))
}

//go:generate go test -trimpath -run=TestStack_CountOK -v github.com/Min-Feng/goutils/errorX
func TestStack_CountOK(t *testing.T) {
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
	stack := Stack(finalErr)

	// spew.Dump(stack)

	expectedStackCount := 5
	assert.Equal(t, expectedStackCount, len(stack))
}
