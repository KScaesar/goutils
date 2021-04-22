package errorX

func Code(err error) int {
	base, ok := Cause(err).(*baseError)
	if !ok {
		return -1
	}
	return base.code
}

func HTTPCode(err error) int {
	base, ok := Cause(err).(*baseError)
	if !ok {
		return -1
	}
	return base.httpCode
}

// SimpleInfo 只取第一次 Wrap 的資訊, 後續 WithMsg 的訊息不會存在
func SimpleInfo(err error) string {
	msgErr := err
	for err != nil {
		cause, ok := err.(pkgErrCause)
		if !ok {
			return msgErr.Error()
		}
		msgErr = err
		err = cause.Unwrap()
	}
	return "no error, this all OK"
}

// Description 自定義錯誤 文字描述
func Description(err error) string {
	return Cause(err).Error()
}
