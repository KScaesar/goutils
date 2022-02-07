package goutils

import (
	"bytes"
	"database/sql/driver"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Min-Feng/goutils/errors"
)

func init() {
	time.Local = time.UTC
}

const MyTimeFormat = "2006-01-02 15:04:05 -07:00"

var defaultTimeSpec = timeSpec{
	list: []string{
		MyTimeFormat,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006-01-02 15:04:05Z07:00",
		time.RFC3339,
		"2006-01-02T15:04",
		"2006-01-02T15:04:05",
	},
}

type timeSpec struct {
	mu   sync.Mutex
	list []string
}

func (v *timeSpec) appendToHead(layout string) {
	v.mu.Lock()
	newSlice := make([]string, 1, len(v.list)+1)
	newSlice[0] = layout
	newSlice = append(newSlice, v.list...)
	v.list = newSlice
	v.mu.Unlock()
}

func RegisterTimeSpec(layout string) {
	defaultTimeSpec.appendToHead(layout)
}

func TimeParse(layout string, utc bool) (t time.Time, err error) {
	for _, spec := range defaultTimeSpec.list {
		t, err = time.Parse(spec, layout)
		if err == nil {
			if utc {
				return t.UTC(), nil
			}
			return t, nil
		}
	}

	err = errors.Wrap(errors.ErrSystem, err.Error())
	return
}

type Time time.Time

func (t *Time) Scan(src interface{}) error {
	switch s := src.(type) {
	case time.Time:
		*t = Time(s)

	case []byte:
		stdTime, err := TimeParse(string(s), true)
		if err != nil {
			return err
		}
		*t = Time(stdTime)

	case string:
		stdTime, err := TimeParse(s, true)
		if err != nil {
			return err
		}
		*t = Time(stdTime)
	}

	return nil
}

func (t Time) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	timeString := string(bytes.Trim(data, `"`))

	if timeString == "null" || timeString == "" || timeString == "nil" {
		return nil
	}

	std, err := TimeParse(timeString, true)
	*t = Time(std)
	return err
}

func (t *Time) UnmarshalBSONValue(b bsontype.Type, bytes []byte) error {
	rv := bson.RawValue{Type: b, Value: bytes}
	*t = Time(rv.Time().UTC())
	return nil
}

func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	targetTime := primitive.NewDateTimeFromTime(time.Time(t))
	return bson.MarshalValue(targetTime)
}

func (t Time) String() string {
	return t.ProtoType().Format(MyTimeFormat)
}

func (t *Time) UnmarshalText(text []byte) error {
	std, err := TimeParse(string(text), true)
	*t = Time(std)
	return err
}

func (t Time) MarshalText() (text []byte, err error) {
	return []byte(t.ProtoType().Format(MyTimeFormat)), nil
}

func (t *Time) ProtoType() time.Time {
	return time.Time(*t)
}
