package goutils

import (
	"bytes"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"

	"github.com/Min-Feng/goutils/errors"
)

var timeSpec = []string{
	MyTimeFormat,
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
	"2006-01-02 15:04:05Z07:00",
	time.RFC3339,
	"2006-01-02T15:04",
	"2006-01-02T15:04:05",
}

func TimeParse(timeLayout string) (t time.Time, err error) {
	for _, spec := range timeSpec {
		t, err = time.ParseInLocation(spec, timeLayout, time.Local)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.Wrap(errors.ErrSystem, err.Error())
}

const MyTimeFormat = "2006-01-02 15:04:05 -07:00"

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	timeString := string(bytes.Trim(data, `"`))

	if timeString == "null" || timeString == "" || timeString == "nil" {
		return nil
	}

	std, err := TimeParse(timeString)
	*t = Time(std)
	return err
}

func (t *Time) UnmarshalBSONValue(b bsontype.Type, bytes []byte) error {
	rv := bson.RawValue{Type: b, Value: bytes}
	*t = Time(rv.Time())
	return nil
}

func (t Time) String() string {
	return t.ProtoType().String()
}

func (t *Time) UnmarshalText(text []byte) error {
	std, err := TimeParse(string(text))
	*t = Time(std)
	return err
}

func (t Time) MarshalText() (text []byte, err error) {
	return []byte(t.ProtoType().Format(MyTimeFormat)), nil
}

func (t Time) ProtoType() time.Time {
	return time.Time(t)
}
