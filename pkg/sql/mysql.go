package sql

import (
	"fmt"
	"runtime"
)

type MySqlConfig struct {
	Host              string
	Port              int
	Username          string
	Password          string
	Database          string
	MaxIdleConnection int
	MaxOpenConnection int
}

func NewMySqlConfig(host, userName, password, database string, port, maxIdleConnection, maxOpenConnection int) Config {
	return &MySqlConfig{
		Host:              host,
		Port:              port,
		Username:          userName,
		Password:          password,
		Database:          database,
		MaxIdleConnection: maxIdleConnection,
		MaxOpenConnection: maxOpenConnection,
	}
}

func (c *MySqlConfig) ConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

func (c *MySqlConfig) Dialect() DriverName {
	return MYSQL
}

func (c *MySqlConfig) GetMaxOpenConnection() int {
	limit := c.MaxOpenConnection

	if limit <= 0 {
		limit = (runtime.NumCPU() * 2) + 16
	}

	if limit > 1024 {
		limit = 1024
	}

	return limit
}

func (c *MySqlConfig) GetMaxIdleConnection() int {
	limit := c.MaxIdleConnection

	if limit <= 0 {
		limit = runtime.NumCPU() + 8
	}

	if limit > c.GetMaxOpenConnection() {
		limit = c.GetMaxOpenConnection()
	}

	return limit
}
