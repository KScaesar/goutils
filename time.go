package goutils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"

	"github.com/Min-Feng/goutils/errors"
)

var timeSpec = []string{
	"2006-01-02 15:04:05",
	"2006-01-02",
	"2006-01-02 15:04",
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02T15:04",
	"2006-01-02T15:04:05",
	time.RFC3339,
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

const TimeFormat = "2006-01-02 15:04:05 -07:00"

type TimeViewModel time.Time

func (t *TimeViewModel) UnmarshalBSONValue(b bsontype.Type, bytes []byte) error {
	rv := bson.RawValue{Type: b, Value: bytes}
	*t = TimeViewModel(rv.Time())
	return nil
}

func (t TimeViewModel) String() string {
	return t.ProtoType().String()
}

func (t TimeViewModel) MarshalText() (text []byte, err error) {
	const defaultFormant = "2006-01-02 15:04:05 -07:00"
	return []byte(t.ProtoType().Format(defaultFormant)), nil
}

func (t TimeViewModel) ProtoType() time.Time {
	return time.Time(t)
}
