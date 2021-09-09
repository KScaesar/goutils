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

func (c *RMDBConfig) MaxConn_() int {
	if c.MaxConn <= 0 {
		const defaultSize = 8
		return defaultSize
	}
	return c.MaxConn
}

func (c *RMDBConfig) MaxIdleConn_() int {
	if c.MaxConn <= 0 {
		const defaultSize = 4
		return defaultSize
	}
	return c.MaxIdleConn
}
