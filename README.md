# introduction

依據 clean 架構 及 ddd 的思想  
放置開發 golang 常用的 技術元件

定義技術元件的抽象(interface)，藉此隔離實作  
且進行單元測試時, 不需要依賴外部 IO

定義抽象 同時可以將 開發程式的專注力  
放在 領域知識 而不是 技術細節

## 重點套件

- [database](./database)
- [errors](./errors)
- [message](./message)

## go test

```bash
go test ./... -trimpath -count 1 -tags integration

go test ./... -trimpath -count 1
```

## context.Context 隱藏哪些技術元件

- traceID
- log
- transaction object
    1. *gorm.Tx
    2. mongo.Session

