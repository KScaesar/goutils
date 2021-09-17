package logger

type Kind int

//go:generate stringer -type=Kind -linecomment -output=custom_field_string.go
const (
	KindUnknown      Kind = iota // unknown
	KindHTTP                     // http
	KindWebsocket                // websocket
	KindMessageQueue             // mq
	KindApplication              // app
)
