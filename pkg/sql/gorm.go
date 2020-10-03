package sql

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/go-impatient/gaia/pkg/mutex"
)

// OpenConnection... 连接数据库
func OpenConnection(config Config) (*gorm.DB, error) {
	mutex.Db.Lock()
	defer mutex.Db.Unlock()

	dialect, err := getDialect(config.Dialect(), config.ConnectionString())
	if err != nil {
		log.Fatalf("Database driver failed: [%s]", err)
	}

	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: &nopLogger{},
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名, 例如: 't_user'
		},
	})

	// 尝试多连接几次, 是否能连接成功
	if err != nil || db == nil {
		for i := 1; i <= 12; i++ {
			db, err := gorm.Open(dialect, &gorm.Config{
				Logger: &nopLogger{},
			})

			if db != nil && err == nil {
				break
			}

			time.Sleep(5 * time.Second)
		}

		if err != nil || db == nil {
			log.Fatalf("Database connection failed: [%s]", err)
		}
	}

	// 数据库调优
	if sqlDB, err := db.DB(); err == nil {
		// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
		sqlDB.SetMaxIdleConns(config.GetMaxIdleConnection())
		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(config.GetMaxOpenConnection())

		// SetConnMaxLifetime 设置了连接可复用的最大时间。
		sqlDB.SetConnMaxLifetime(10 * time.Minute)

		db.Logger = db.Logger.LogMode(logger.Info)

	}

	return db, nil
}

// 获取方言
func getDialect(driver DriverName, dsn string) (gorm.Dialector, error) {
	switch driver {
	case SQLITE:
		return sqlite.Open(dsn), nil
	case MYSQL:
		return mysql.Open(dsn), nil
	case POSTGRES:
		return postgres.Open(dsn), nil
	default:
		return nil, fmt.Errorf("no database driver named `%s` found", driver)
	}
}
