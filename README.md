# goutils

開發 golang 時, 放置常用的工具集

## go test

```bash
go test ./... -trimpath -count 1 -tags integration

go test ./... -trimpath -count 1
```

## context.Context 放置哪些技術元件

- traceID
- log
- transaction object (*gorm.Tx, mongo.Session)

