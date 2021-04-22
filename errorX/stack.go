package errorX

// Stacks 取出所有 Wrap 時, 記錄的 stack
func Stacks(err error) []StackTrace {
	pgkStacks := getPKGStacks(err)
	length := len(pgkStacks)
	last := length - 1
	stacks := make([]StackTrace, length)

	for i, pgkStack := range pgkStacks {
		stack := transformPGKStack(pgkStack)
		stacks[last-i] = stack // 為了讓最根本的錯誤, 顯示在第一個
	}
	return stacks
}

// 取得 pkg/errors 的 stack 型別
// 目標是 拿到所有 Wrap 時, 記錄的 stack
func getPKGStacks(err error) []pkgErrStack {
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

func transformPGKStack(iStack pkgErrStack) StackTrace {
	if iStack == nil {
		return nil
	}

	defaultBuff := 10
	stack := make([]Frame, 0, defaultBuff)

	pgkStack := iStack.StackTrace()
	last := len(pgkStack) - 1

	// 第 0 個 pgkStack 不需要, 因為是重新包裝過的 Wrap func
	// pkg/error 沒有提供 runtime.Caller skip 參數
	// 所以 i 從 1 開始
	for i := 1; i <= last; i++ {
		frame, _ := pgkStack[i].MarshalText()
		stack = append(stack, Frame(frame))
	}
	return stack
}

// Stack 只會找到最後一個 Wrap 時, 所記錄的 stack
func Stack(err error) StackTrace {
	pkgStack := getPKGStack(err)
	stack := transformPGKStack(pkgStack)
	return stack
}

// 取得 pkg/errors 的 stack 型別
// 目標是 拿到最後一個 Wrap 時, 所記錄的 stack
func getPKGStack(err error) pkgErrStack {
	for err != nil {
		pgkStack, ok := err.(pkgErrStack)
		if ok {
			return pgkStack
		}

		cause, ok := err.(pkgErrCause)
		if !ok {
			return nil
		}
		err = cause.Unwrap()
	}
	return nil
}
