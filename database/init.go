package database

import "time"

func init() {
	location, _ := time.LoadLocation("UTC")
	time.Local = location
}
