# readme

## 注意事項

zerolog 大多數的方法 (method) 都是 value 語意  
不影響之前的 log 資訊

少部分 pointer 語意 method, 如下:

- func (l *Logger) WithContext(ctx context.Context) context.Context
- func (l *Logger) UpdateContext(update func(c Context) Context)

## level

只使用四種等級

- "debug"
- "info"
- "error"
- "panic"