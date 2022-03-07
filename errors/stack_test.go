package errors_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KScaesar/goutils/errors"
)

//go:generate go test -trimpath -run=TestStacks -v github.com/KScaesar/goutils/errors
func TestStacks(t *testing.T) {
	topDefinedErr := errors.New(10001, http.StatusInternalServerError, "system internal error")

	repo := func() error {
		infraErr := errors.NewStdError("connect fail")
		return errors.Wrap(topDefinedErr, infraErr.Error())
	}
	rollback := func() error {
		return errors.NewStdError("data center bang")
	}
	useCase := func() error {
		repoErr := repo()
		if repoErr != nil {
			rollbackErr := rollback()
			if rollbackErr != nil {
				// 遇到第二個 error 或許可以用 Wrap 再次包裝 保留 stack 訊息
				// 並把 次要錯誤(rollbackErr)的文字訊息 加入到 主要錯誤(repoErr)
				// 若想保留 rollbackErr stack, 可以考慮在此處加上 log
				return errors.Wrap(repoErr, "update entity failed: %v", rollbackErr)
			}

			// 如果只是單一錯誤往上傳遞 只用 WithMsg 增加文字訊息就好
			return errors.WrapMessage(repoErr, "update entity failed")
		}
		return nil
	}
	finalErr := useCase()
	stacks := errors.Stacks(finalErr)

	expectedStacks := []errors.Stack{
		{
			"github.com/KScaesar/goutils/errors_test.TestStacks.func1 github.com/KScaesar/goutils/errors_test/stack_test.go:18 ",
			"github.com/KScaesar/goutils/errors_test.TestStacks.func3 github.com/KScaesar/goutils/errors_test/stack_test.go:24 ",
			"github.com/KScaesar/goutils/errors_test.TestStacks github.com/KScaesar/goutils/errors_test/stack_test.go:39 ",
			"testing.tRunner testing/testing.go:1123 ",
			"runtime.goexit runtime/asm_amd64.s:1374 ",
		},
		{
			"github.com/KScaesar/goutils/errors_test.TestStacks.func3 github.com/KScaesar/goutils/errors_test/stack_test.go:31 ",
			"github.com/KScaesar/goutils/errors_test.TestStacks github.com/KScaesar/goutils/errors_test/stack_test.go:39 ",
			"testing.tRunner testing/testing.go:1123 ",
			"runtime.goexit runtime/asm_amd64.s:1374 ",
		},
	}

	for i := 0; i < 3; i++ {
		assert.Equal(t, expectedStacks[0][i], stacks[0][i])
	}
	for i := 0; i < 2; i++ {
		assert.Equal(t, expectedStacks[1][i], stacks[1][i])
	}
	assert.Len(t, expectedStacks, len(stacks))
	assert.Len(t, expectedStacks[0], len(stacks[0]))
}
