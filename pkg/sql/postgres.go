package sql

import (
	"fmt"
	"runtime"
)

type PostgresConfig struct {
	Host              string
	Port              int
	Username          string
	Password          string
	Database          string
	SSL               bool
	MaxIdleConnection int
	MaxOpenConnection int
}

func NewPostgresConfig(host, userName, password, database string, ssl bool, port, maxIdleConnection, maxOpenConnection int) Config {
	return &PostgresConfig{
		Host:              host,
		Port:              port,
		Username:          userName,
		Password:          password,
		Database:          database,
		SSL:               ssl,
		MaxIdleConnection: maxIdleConnection,
		MaxOpenConnection: maxOpenConnection,
	}
}

func (c *PostgresConfig) ConnectionString() string {
	sslMode := "enable"
	if !c.SSL {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s binary_parameters=yes",
		c.Host,
		c.Port,
		c.Username,
		c.Database,
		c.Password,
		sslMode,
	)
}

func (c *PostgresConfig) Dialect() DriverName {
	return POSTGRES
}

func (c *PostgresConfig) GetMaxOpenConnection() int {
	limit := c.MaxOpenConnection

	if limit <= 0 {
		limit = (runtime.NumCPU() * 2) + 16
	}

	if limit > 1024 {
		limit = 1024
	}

	return limit
}

func (c *PostgresConfig) GetMaxIdleConnection() int {
	limit := c.MaxIdleConnection

	if limit <= 0 {
		limit = runtime.NumCPU() + 8
	}

	if limit > c.GetMaxOpenConnection() {
		limit = c.GetMaxOpenConnection()
	}

	return limit
}
