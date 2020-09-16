package service

import (
	"database/sql"
	"github.com/go-impatient/gaia/app/conf"
	sqlx "github.com/go-impatient/gaia/pkg/sql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log"
	"time"
)

// SQL ...
type SQL struct {
	DB *gorm.DB
}

// NewSQL ...
func NewSQL() *SQL {
	dialect := conf.Config.Database.Dialect
	host := conf.Config.Database.Host
	port := conf.Config.Database.Port
	username := conf.Config.Database.Username
	password := conf.Config.Database.Password
	database := conf.Config.Database.Database
	maxIdleConnection := conf.Config.Database.MaxIdleConns
	maxOpenConnection := conf.Config.Database.MaxOpenConns
	ssl := conf.Config.Database.Ssl

	// 根据方言选择数据库配置
	driverName, _ := sqlx.ParseDriverName(dialect)
	var sqlConfig sqlx.Config
	switch driverName {
	case sqlx.SQLITE:
	case sqlx.MYSQL:
		sqlConfig = sqlx.NewMySqlConfig(host, username, password, database, port, maxIdleConnection, maxOpenConnection)
	case sqlx.POSTGRES:
		sqlConfig = sqlx.NewPostgresConfig(host, username, password, database, ssl, port, maxIdleConnection, maxOpenConnection)
	default:
		log.Panicf("Database dialect `%s` not supported \n", dialect)
	}

	// 连接数据库
	db, _ := sqlx.OpenConnection(sqlConfig)

	// 测试数据库心跳
	if err := pingDB(db); err != nil {
		log.Printf("Failed to connect database, got error %v\n", err)
	}

	return &SQL{
		DB: db,
	}
}

// DBStats ... 数据库统计信息
func (s *SQL) DBStats() sql.DBStats {
	sqlDB, err := s.DB.DB()
	if err != nil {
		log.Printf("Failed to connect database, got error %v\n", err)
	}
	return sqlDB.Stats()
}

// CloseDB ...
func (s *SQL) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		log.Printf("Disconnect from database failed: [%s]", err)
	}
	return sqlDB.Close()
}

// migrate migrates database schemas ...
func (s *SQL) Migrate(Models []interface{}) error {
	if err := s.DB.AutoMigrate(Models...); err != nil {
		return errors.Wrap(err, "auto migrate tables failed")
	}

	return nil
}

// creates necessary database tables
func (s *SQL) CreateTables(Models []interface{}) error {
	for _, model := range Models {
		if !s.DB.Migrator().HasTable(model) {
			if err := s.DB.Migrator().CreateTable(model); err != nil {
				return errors.Wrap(err, "create table failed")
			}
		}
	}

	return nil
}

func (s *SQL) DeleteTables(Models []interface{}) error {
	if err := s.DB.Migrator().DropTable(Models...); err != nil {
		return errors.Wrap(err, "delete table failed")
	}
	return nil
}

// pingDB... 数据库心跳
func pingDB(s *gorm.DB) (err error) {
	for i := 0; i < 30; i++ {
		sqlDB, err := s.DB()
		if err == nil {
			sqlDB.Ping()
		}
		time.Sleep(time.Second)
	}
	return
}
