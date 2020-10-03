package database

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/go-impatient/gaia/app/conf"
	orm "github.com/go-impatient/gaia/pkg/sql"
)

var (
	defaultOrm *gorm.DB
	ormMap sync.Map
)

// AsDefault alias for "default"
const AsDefault = "default"

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
	driverName, _ := orm.ParseDriverName(dialect)
	var sqlConfig orm.Config
	switch driverName {
	case orm.SQLITE:
	case orm.MYSQL:
		sqlConfig = orm.NewMySqlConfig(host, username, password, database, port, maxIdleConnection, maxOpenConnection)
	case orm.POSTGRES:
		sqlConfig = orm.NewPostgresConfig(host, username, password, database, ssl, port, maxIdleConnection, maxOpenConnection)
	default:
		log.Panicf("Database dialect `%s` not supported \n", dialect)
	}

	// 连接数据库
	db, _ := orm.OpenConnection(sqlConfig)

	// 测试数据库心跳
	if err := pingDB(db); err != nil {
		log.Printf("Failed to connect database, got error %v\n", err)
	}

	defaultOrm = db
	ormMap.Store("default", db)

	return &SQL{
		DB: db,
	}
}

// Orm returns an orm's db.
func Orm(name ...string) *gorm.DB {
	if len(name) == 0 || name[0] == AsDefault {
		if defaultOrm == nil {
			log.Panicf("Invalid db `%s` \n", AsDefault)
		}
		return defaultOrm
	}

	v, ok := ormMap.Load(name[0])
	if ! ok {
		log.Panicf("Invalid db `%s` \n", name[0])
	}

	return v.(*gorm.DB)
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
