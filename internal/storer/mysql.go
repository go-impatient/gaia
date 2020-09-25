package storer

import (
	"database/sql"
	"log"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/go-impatient/gaia/pkg/dig"
)

func init() {
	dig.DigProvide(NewDB)
}

// SQL ...
type SQL struct {
	DB *gorm.DB
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

// NewDB ...
func NewDB() *SQL {
	// 连接数据库
	db := connectDB()

	// 测试数据库心跳
	if err := pingDB(db); err != nil {
		log.Printf("Failed to connect database, got error %v\n", err)
	}

	return &SQL{
		DB: db,
	}
}

// connectDB... 连接数据库, 配置数据库
func connectDB() *gorm.DB {
	var err error
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection failed: [%s]", err)
	} else {
		// 数据库调优
		if sqlDB, err := db.DB(); err == nil {
			// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
			sqlDB.SetMaxIdleConns(10)

			// SetMaxOpenConns 设置打开数据库连接的最大数量。
			sqlDB.SetMaxOpenConns(100)

			// SetConnMaxLifetime 设置了连接可复用的最大时间。
			sqlDB.SetConnMaxLifetime(time.Hour)

			db.Logger = db.Logger.LogMode(logger.Info)
		} else {
			log.Fatalf("Database connection failed: [%s]", err)
		}
	}
	return db
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
