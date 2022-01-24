package xLog

type TriggerKind int

//go:generate stringer -type=TriggerKind -linecomment -output=trigger_kind_string.go
const (
	TriggerKindUnknown   TriggerKind = iota // unknown
	TriggerKindHttp                         // http
	TriggerKindWebsocket                    // websocket
	TriggerKindGrpc                         // grpc
	TriggerKindMQ                           // mq
)
