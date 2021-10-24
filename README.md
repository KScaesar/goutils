# introduction

依據 clean 架構 及 ddd 的思想  
放置開發 golang 常用的基礎技術元件  

定義抽象依賴，隔離實作  
期望開發程式的專注力  
放在 領域知識 而不是技術細節  

## go test

```bash
go test ./... -trimpath -count 1 -tags integration

go test ./... -trimpath -count 1
```

## context.Context 放置哪些技術元件

- traceID
- log
- transaction object (*gorm.Tx, mongo.Session)

