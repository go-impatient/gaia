package sql

import (
	"database/sql"
	"fmt"
	"strings"
)

// 数据库驱动名称
type DriverName string

const (
	SQLITE   DriverName = "sqlite"
	MYSQL    DriverName = "mysql"
	POSTGRES DriverName = "postgres"
)

func ParseDriverName(driverName string) (DriverName, error) {
	switch strings.ToLower(driverName) {
	case "sqlite":
		return SQLITE, nil
	case "mysql":
		return MYSQL, nil
	case "postgres":
		return POSTGRES, nil
	default:
		return SQLITE, fmt.Errorf("未知的数据库驱动名称: '%s', 默认名称是 'sqlite'", driverName)
	}
}

type Config interface {
	// DSN
	ConnectionString() string
	// 数据库驱动类型
	Dialect() DriverName
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	GetMaxOpenConnection() int
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	GetMaxIdleConnection() int
}

type Migrater interface {
	// 数据库迁移
	Migrate() error
}

type DBStats interface {
	// 数据库统计信息
	DBStats() sql.DBStats
}
