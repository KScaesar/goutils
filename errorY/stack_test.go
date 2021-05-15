package errorY_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/goutils/errorY"
)

//go:generate go test -trimpath -run=TestStacks -v github.com/Min-Feng/goutils/errorY
func TestStacks(t *testing.T) {
	topDefinedErr := errorY.New(10001, http.StatusInternalServerError, "system internal error")

	repo := func() error {
		infraErr := errorY.NewStdError("connect fail")
		return errorY.Wrap(topDefinedErr, infraErr.Error())
	}
	rollback := func(error) error {
		return errorY.NewStdError("data center bang")
	}
	useCase := func() error {
		repoErr := repo()
		if repoErr != nil {
			fixErr := rollback(repoErr)
			if fixErr != nil {
				// 遇到第二個 error 或許可以用 Wrap 再次包裝 保留 stack 訊息
				// 並把 次要錯誤(fixErr)的文字訊息 加入到 主要錯誤(repoErr)
				// 若想保留 fixErr stack, 可以考慮在此處加上 log
				return errorY.Wrap(repoErr, "update entity failed: %v", fixErr)
			}

			// 如果只是單一錯誤往上傳遞 只用 WithMsg 增加文字訊息就好
			return errorY.WrapMessage(repoErr, "update entity failed")
		}
		return nil
	}
	finalErr := useCase()
	stacks := errorY.Stacks(finalErr)

	expectedStacks := []errorY.StackTrace{
		{
			"github.com/Min-Feng/goutils/errorY_test.TestStacks.func1 github.com/Min-Feng/goutils/errorY_test/stack_test.go:18 ",
			"github.com/Min-Feng/goutils/errorY_test.TestStacks.func3 github.com/Min-Feng/goutils/errorY_test/stack_test.go:24 ",
			"github.com/Min-Feng/goutils/errorY_test.TestStacks github.com/Min-Feng/goutils/errorY_test/stack_test.go:39 ",
			"testing.tRunner testing/testing.go:1123 ",
			"runtime.goexit runtime/asm_amd64.s:1374 ",
		},
		{
			"github.com/Min-Feng/goutils/errorY_test.TestStacks.func3 github.com/Min-Feng/goutils/errorY_test/stack_test.go:31 ",
			"github.com/Min-Feng/goutils/errorY_test.TestStacks github.com/Min-Feng/goutils/errorY_test/stack_test.go:39 ",
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
