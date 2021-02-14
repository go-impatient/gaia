package database

import (
	"testing"

	"gorm.io/gorm/logger"

	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"

	"github.com/go-impatient/gaia/app/conf"
	"github.com/go-impatient/gaia/internal/model"
)

type DatabaseTestSuite struct {
	suite.Suite
}

func (suite *DatabaseTestSuite) TestGetConnection() {
	cfg, err := conf.InitConfig("./../../config/config.json")
	suite.NotNil(err)
	suite.NotNil(cfg)
	sql := NewSQL()
	db := sql.DB
	suite.NotNil(db)
	suite.NotNil(GetDB())
	suite.Equal(GetDB(), db)
	suite.Nil(sql.Close())
	suite.Nil(db)
}

func (suite *DatabaseTestSuite) TestLogLevel() {
	cfg, err := conf.InitConfig("./../../config/config.json")
	suite.NotNil(err)
	suite.NotNil(cfg)
	sql := NewSQL()
	db := sql.DB

	suite.Equal(logger.Default.LogMode(logger.Silent), db.Logger)
	sql.Close()
}

func (suite *DatabaseTestSuite) TestModelAndMigrate() {
	cfg, err := conf.InitConfig("./../../config/config.json")
	suite.NotNil(err)
	suite.NotNil(cfg)
	sql := NewSQL()
	db := sql.DB

	ClearRegisteredModels()
	RegisterModel(&model.User{})
	suite.Len(models, 1)

	registeredModels := GetRegisteredModels()
	suite.Len(registeredModels, 1)
	suite.Same(models[0], registeredModels[0])

	sql.Migrate()
	ClearRegisteredModels()
	suite.Equal(0, len(models))

	defer db.Migrator().DropTable(&model.User{})

	rows, err := db.Raw("SHOW TABLES;").Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		name := ""
		if err := rows.Scan(&name); err != nil {
			panic(err)
		}
		if name == "user" {
			found = true
			break
		}
	}

	suite.True(found)
}

func (suite *DatabaseTestSuite) TestInitializers() {
	AddInitializer(func(db *gorm.DB) {
		db.Config.SkipDefaultTransaction = true
		db.Statement.Settings.Store("gorm:table_options", "ENGINE=InnoDB")
	})

	suite.Len(initializers, 2)

	cfg, err := conf.InitConfig("./../../config/config.json")
	suite.NotNil(err)
	suite.NotNil(cfg)
	sql := NewSQL()
	db := sql.DB
	suite.True(db.Config.SkipDefaultTransaction)

	suite.Nil(sql.Close())

	AddInitializer(func(db *gorm.DB) {
		db.Statement.Settings.Store("gorm:table_options", "ENGINE=InnoDB")
	})
	suite.Len(initializers, 2)

	cfg2, err2 := conf.InitConfig("./../../config/config.json")
	suite.NotNil(err2)
	suite.NotNil(cfg2)
	sql2 := NewSQL()
	db2 := sql.DB
	suite.True(db2.Config.SkipDefaultTransaction)
	val, ok := db2.Get("gorm:table_options")
	suite.True(ok)
	suite.Equal("ENGINE=InnoDB", val)

	suite.Nil(sql2.Close())

	ClearInitializers()
	suite.Empty(initializers)
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
