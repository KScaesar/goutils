package goutils

import (
	"bytes"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"

	"github.com/Min-Feng/goutils/errors"
)

func init() {
	time.Local = time.UTC
}

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

func TimeParse(timeLayout string, utc bool) (t time.Time, err error) {
	for _, spec := range timeSpec {
		t, err = time.Parse(spec, timeLayout)
		if err == nil {
			if utc {
				return t.UTC(), nil
			}
			return t, nil
		}
	}
	return time.Time{}, errors.Wrap(errors.ErrSystem, err.Error())
}

var MyTimeFormat = "2006-01-02 15:04:05 -07:00"

type Time time.Time

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

func (t Time) String() string {
	return t.Unwrap().Format(MyTimeFormat)
}

func (t *Time) UnmarshalText(text []byte) error {
	std, err := TimeParse(string(text), true)
	*t = Time(std)
	return err
}

func (t Time) MarshalText() (text []byte, err error) {
	return []byte(t.Unwrap().Format(MyTimeFormat)), nil
}

func (t Time) Unwrap() time.Time {
	return time.Time(t)
}
