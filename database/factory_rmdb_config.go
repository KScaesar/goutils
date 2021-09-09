package database

type RMDBConfig struct {
	User        string
	Password    string
	Host        string
	Port        string
	Database    string
	MaxConn     int
	MaxIdleConn int
}

func (c *RMDBConfig) setDefaultValue() {
	if c.MaxConn <= 0 {
		const default_ = 8
		c.MaxConn = default_
	}

	if c.MaxIdleConn <= 0 {
		const default_ = 4
		c.MaxIdleConn = default_
	}
}
