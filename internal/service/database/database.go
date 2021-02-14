package database

import (
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/go-impatient/gaia/app/conf"
	"github.com/go-impatient/gaia/pkg/sql"
)

// Initializer is a function meant to modify a connection settings
// at the global scope when it's created.
//
// Use `db.InstantSet()` and not `db.Set()`, since the latter clones
// the gorm.DB instance instead of modifying it.
type Initializer func(*gorm.DB)

var (
	defaultSQL *gorm.DB
	sqlMap     sync.Map

	mu           sync.Mutex
	models       []interface{}
	initializers []Initializer
)

// AsDefault alias for "default"
const AsDefault = "default"

// SQL ...
type SQL struct {
	DB *gorm.DB
}

// NewSQL ...
func NewSQL() *SQL {
	mode := conf.AppConfig.Mode
	dialect := conf.DBConfig.Dialect
	host := conf.DBConfig.Host
	port := conf.DBConfig.Port
	username := conf.DBConfig.Username
	password := conf.DBConfig.Password
	database := conf.DBConfig.Database
	maxIdleConnection := conf.DBConfig.MaxIdleConns
	maxOpenConnection := conf.DBConfig.MaxOpenConns
	ssl := conf.DBConfig.Ssl

	// 根据方言选择数据库配置
	driverName, _ := sql.ParseDriverName(dialect)
	var sqlConfig sql.SQLConfig
	switch driverName {
	case sql.SQLITE:
		sqlConfig = sql.NewSqliteConfig(mode, database, false)
	case sql.MYSQL:
		sqlConfig = sql.NewMySqlConfig(mode, host, username, password, database, port, maxIdleConnection, maxOpenConnection)
	case sql.POSTGRES:
		sqlConfig = sql.NewPostgresConfig(mode, host, username, password, database, ssl, port, maxIdleConnection, maxOpenConnection)
	default:
		log.Panicf("Database dialect `%s` not supported \n", dialect)
	}

	// 连接数据库
	db, err := sql.OpenConnection(sqlConfig)
	if err != nil {
		panic(err)
	}

	// Initializer functions are meant to modify a connection settings
	for _, initializer := range initializers {
		initializer(db)
	}

	// 测试数据库心跳
	if err := pingDB(db); err != nil {
		log.Printf("Failed to connect database, got error %v\n", err)
	}

	defaultSQL = db
	sqlMap.Store("default", db)

	return &SQL{
		DB: db,
	}
}

// GetDB 全局的SQL服务
func GetDB(name ...string) *gorm.DB {
	if len(name) == 0 || name[0] == AsDefault {
		if defaultSQL == nil {
			log.Panicf("Invalid db `%s` \n", AsDefault)
		}
		return defaultSQL
	}

	v, ok := sqlMap.Load(name[0])
	if !ok {
		log.Panicf("Invalid db `%s` \n", name[0])
	}

	return v.(*gorm.DB)
}

// AddInitializer adds a database connection initializer function.
// Initializer functions are meant to modify a connection settings
// at the global scope when it's created.
//
// Initializer functions are called in order, meaning that functions
// added last can override settings defined by previous ones.
func AddInitializer(initializer Initializer) {
	initializers = append(initializers, initializer)
}

// ClearInitializers remove all database connection initializer functions.
func ClearInitializers() {
	initializers = []Initializer{}
}

// Close the database connections if they exist.
func (s *SQL) Close() error {
	mu.Lock()
	defer mu.Unlock()

	sqlDB, err := s.DB.DB()
	if err != nil {
		log.Printf("Disconnect from database failed: [%s]", err)
	}
	return sqlDB.Close()
}

// RegisterModel registers a model for auto-migration.
// When writing a model file, you should always register it in the init() function.
//  func init() {
//		database.RegisterModel(&MyModel{})
//  }
func RegisterModel(model interface{}) {
	models = append(models, model)
}

// GetRegisteredModels get the registered models.
// The returned slice is a copy of the original, so it
// cannot be modified.
func GetRegisteredModels() []interface{} {
	return append(make([]interface{}, 0, len(models)), models...)
}

// ClearRegisteredModels unregister all models.
func ClearRegisteredModels() {
	models = []interface{}{}
}

// Migrate migrates all registered models.
func (s *SQL) Migrate() error {
	if err := s.DB.AutoMigrate(models...); err != nil {
		return errors.Wrap(err, "bcrypt migrate tables failed")
	}

	return nil
}

// creates necessary database tables
func (s *SQL) CreateTables() error {
	for _, model := range models {
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
