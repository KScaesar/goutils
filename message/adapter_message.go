package message

import (
	"encoding/json"

	"github.com/Min-Feng/goutils/errors"
)

// AdapterMessage 從 queue 接收資料時, 可以利用此型別進行判斷處理
type AdapterMessage struct {
	InfoBase
	Payload json.RawMessage `json:"payload"`
}

func (r *AdapterMessage) MatchTopic(topic Topic) bool {
	return r.InfoBase.MessageTopic == topic
}

func (r *AdapterMessage) UnmarshalPayload(obj interface{}) error {
	err := json.Unmarshal(r.Payload, obj)
	if err != nil {
		return errors.Wrap(errors.ErrSystem, err.Error())
	}
	return nil
}
