package database

import "fmt"

type RMDBConfig struct {
	User        string `configs:"user"`
	Password    string `configs:"password"`
	Host        string `configs:"host"`
	Port        string `configs:"port"`
	Database    string `configs:"database"`
	MaxConn     int    `configs:"max_conn"`
	MaxIdleConn int    `configs:"max_idle_conn"`
}

func (c *RMDBConfig) MysqlDSN() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
}
