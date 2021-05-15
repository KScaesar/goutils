# value interface

## fmt
https://pkg.go.dev/fmt#Stringer.String
```
type Stringer interface {
	String() string
}
```

## sql
https://golang.org/pkg/database/sql/driver/#Valuer
```
type ValueConverter interface {
    // ConvertValue converts a value to a driver Value.
    ConvertValue(v interface{}) (Value, error)
}
```

https://pkg.go.dev/database/sql#Scanner.Scan
```
type Scanner interface {
	// Scan assigns a value from a database driver.
	//
	// The src value will be of one of the following types:
	//
	//    int64
	//    float64
	//    bool
	//    []byte
	//    string
	//    time.Time
	//    nil - for NULL values
	//
	// An error should be returned if the value cannot be stored
	// without loss of information.
	//
	// Reference types such as []byte are only valid until the next call to Scan
	// and should not be retained. Their underlying memory is owned by the driver.
	// If retention is necessary, copy their values before the next call to Scan.
	Scan(src interface{}) error
}
```

## json
https://golang.org/pkg/encoding/json/#Marshaler
```
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}
```

https://golang.org/pkg/encoding/json/#Unmarshaler
```
type Unmarshaler interface {
    UnmarshalJSON([]byte) error
}
```

## encoding
https://golang.org/pkg/encoding/#TextMarshaler
```
type TextMarshaler interface {
    MarshalText() (text []byte, err error)
}
```

https://golang.org/pkg/encoding/#TextUnmarshaler
```
type TextUnmarshaler interface {
    UnmarshalText(text []byte) error
}
```