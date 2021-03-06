package errors

import (
	"sync"
)

// Stacks 取出所有 Wrap 時, 記錄的 stack
func Stacks(err error) []Stack {
	pgkStacks := getPkgStacks(err)
	length := len(pgkStacks)
	last := length - 1
	stacks := make([]Stack, length)

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

func transformPgkStack(iStack pkgErrStack) Stack {
	if iStack == nil {
		return make([]Frame, 0)
	}

	defaultBuff := 20
	stack := make(Stack, 0, defaultBuff)

	pgkStack := iStack.StackTrace()
	last := len(pgkStack) - 1

	// 第 0 個 pgkStack 不需要, 因為是重新包裝過的 Wrap func
	// pkg/error 沒有提供 runtime.Caller skip 參數
	// 所以 i 從 1 開始
	for i := 1; i <= last; i++ {
		frame, _ := pgkStack[i].MarshalText()
		if defaultChecker.canRemove(frame) {
			continue
		}
		frame = append(frame, ' ')
		stack = append(stack, Frame(frame))
	}
	return stack
}

var defaultChecker removeFrameChecker

func RegisterFrameFilter(filter ...FrameFilter) {
	defaultChecker.addRule(filter...)
}

// FrameFilter return true 表示 error stack 要過濾 這個 frame, 不希望它出現
type FrameFilter func(frame []byte) bool

type removeFrameChecker struct {
	mu      sync.RWMutex
	filters []FrameFilter
}

func (checker *removeFrameChecker) addRule(f ...FrameFilter) {
	checker.mu.Lock()
	defer checker.mu.Unlock()
	checker.filters = append(checker.filters, f...)
}

// canRemove return true means that the frame can be removed,
// the function does not perform deletion, it is query semantics
func (checker *removeFrameChecker) canRemove(frame []byte) bool {
	checker.mu.RLock()
	defer checker.mu.RUnlock()

	for _, filter := range checker.filters {
		if filter(frame) {
			return true
		}
	}
	return false
}
