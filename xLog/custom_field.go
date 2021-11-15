package xLog

type TriggerKind int

//go:generate stringer -type=TriggerKind -linecomment -output=custom_field_string.go
const (
	TriggerKindUnknown       TriggerKind = iota // unknown
	TriggerKindHttp                             // http
	TriggerKindWebsocket                        // websocket
	TriggerKindGrpc                             // grpc
	TriggerKindMessageBroker                    // mq
)
