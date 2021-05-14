package errorY

import "bytes"

// Stacks 取出所有 Wrap 時, 記錄的 stack
func Stacks(err error) []StackTrace {
	pgkStacks := getPkgStacks(err)
	length := len(pgkStacks)
	last := length - 1
	stacks := make([]StackTrace, length)

	for i, pgkStack := range pgkStacks {
		stack := transformPgkStack(pgkStack)
		stacks[last-i] = stack // 為了讓最初的錯誤, 顯示在第一個, 每次 errorY.Wrap() 就會新增 stack 訊息
	}
	return stacks
}

// 取得 pkg/errors 的 stack 型別
// 目標是 拿到所有 Wrap 時, 記錄的 stack
func getPkgStacks(err error) []pkgErrStack {
	defaultBuff := 3
	pgkStacks := make([]pkgErrStack, 0, defaultBuff)
	for err != nil {
		pgkStack, ok := err.(pkgErrStack)
		if ok {
			pgkStacks = append(pgkStacks, pgkStack)
		}

		cause, ok := err.(pkgErrCause)
		if !ok {
			return pgkStacks
		}
		err = cause.Unwrap()
	}
	return nil
}

func transformPgkStack(iStack pkgErrStack) StackTrace {
	if iStack == nil {
		return make([]Frame, 0)
	}

	defaultBuff := 20
	stack := make([]Frame, 0, defaultBuff)

	pgkStack := iStack.StackTrace()
	last := len(pgkStack) - 1

	// 第 0 個 pgkStack 不需要, 因為是重新包裝過的 Wrap func
	// pkg/error 沒有提供 runtime.Caller skip 參數
	// 所以 i 從 1 開始
	for i := 1; i <= last; i++ {
		frame, _ := pgkStack[i].MarshalText()
		if isExistGinNext(frame) {
			continue
		}
		frame = append(frame, ' ')
		stack = append(stack, Frame(frame))
	}
	return stack
}

var filterGinNext = []byte("gin.(*Context).Next")

func isExistGinNext(frame []byte) bool {
	return bytes.Contains(frame, filterGinNext)
}
