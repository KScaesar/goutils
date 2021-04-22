package errorX

// New 不可隨意新增 自定義錯誤, 儘量使用 Wrapf 產生 error,
// 主要會在全域使用, 以方便進行 error assert Is(),
// 在函數中使用 New 可能是不適當的用法?
func New(code int, httpCode int, description string) error {
	return &baseError{code: code, httpCode: httpCode, description: description}
}

type baseError struct {
	code        int
	httpCode    int
	description string
}

func (b *baseError) Error() string {
	return b.description
}
