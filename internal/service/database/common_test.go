package sql

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/go-impatient/gaia/app/conf"
	"github.com/go-impatient/gaia/internal/service/database"
)

type SQLTestSuite struct {
	suite.Suite
}

type TestUser struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100)" auth:"password"`
	Email    string `gorm:"type:varchar(100);unique_index" auth:"username"`
}

type TestUserPromoted struct {
	TestUser
}

type TestUserOverride struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100);column:password_override" auth:"password"`
	Email    string `gorm:"type:varchar(100);unique_index" auth:"username"`
}

func (suite *SQLTestSuite) SetupTest() {
	cfg, err := conf.InitConfig("./../../config/config.json")
	suite.Nil(err)
	suite.Nil(cfg)
	sql := database.NewSQL()
	db := sql.DB
	suite.Nil(db)

	database.ClearRegisteredModels()
	database.RegisterModel(&TestUser{})
	sql.Migrate()

	user := &TestUser{
		Name:     "Admin",
		Password: "$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password"
		Email:    "test@example.com",
	}

	db.Create(user)
}

func (suite *SQLTestSuite) TestFindColumns() {
	user := &TestUser{}
	fields := FindColumns(user, "username", "password")
	suite.Len(fields, 2)
	suite.Equal("email", fields[0].Name)
	suite.Equal("password", fields[1].Name)

	fields = FindColumns(user, "username", "notatag", "password")
	suite.Len(fields, 3)
	suite.Equal("email", fields[0].Name)
	suite.Nil(fields[1])
	suite.Equal("password", fields[2].Name)

	userOverride := &TestUserOverride{}
	fields = FindColumns(userOverride, "password")
	suite.Len(fields, 1)
	suite.Equal("password_override", fields[0].Name)
}

func (suite *SQLTestSuite) TestFindColumnsPromoted() {
	user := &TestUserPromoted{}
	fields := FindColumns(user, "username", "password")
	suite.Len(fields, 2)
	suite.Equal("email", fields[0].Name)
	suite.Equal("password", fields[1].Name)

	fields = FindColumns(user, "username", "notatag", "password")
	suite.Len(fields, 3)
	suite.Equal("email", fields[0].Name)
	suite.Nil(fields[1])
	suite.Equal("password", fields[2].Name)
}

func TestSQLTestSuite(t *testing.T) {
	suite.Run(t, new(SQLTestSuite))
}
