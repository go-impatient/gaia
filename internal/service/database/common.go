package sql

import (
	"reflect"
	"strings"

	"github.com/go-impatient/gaia/internal/service/database"

	"github.com/go-impatient/gaia/pkg/xslice"
)

// Column matches a column name with a struct field.
type Column struct {
	Name  string
	Field *reflect.StructField
}

// FindColumns in given struct. A field matches if it has a "auth" tag with the given value.
// Returns a slice of found fields, ordered as the input "fields" slice.
// If the nth field is not found, the nth value of the returned slice will be nil.
//
// Promoted fields are matched as well.
//
// Given the following struct and "username", "notatag", "password":
//  type TestUser struct {
// 		gorm.Model
// 		Name     string `gorm:"type:varchar(100)"`
// 		Password string `gorm:"type:varchar(100)" auth:"password"`
// 		Email    string `gorm:"type:varchar(100);unique_index" auth:"username"`
//  }
//
// The result will be the "Email" field, "nil" and the "Password" field.
func FindColumns(strct interface{}, fields ...string) []*Column {
	length := len(fields)
	result := make([]*Column, length)

	value := reflect.ValueOf(strct)
	t := reflect.TypeOf(strct)
	if t.Kind() == reflect.Ptr {
		value = value.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := value.Field(i)
		fieldType := t.Field(i)
		if field.Kind() == reflect.Struct && fieldType.Anonymous {
			// Check promoted fields recursively
			for i, v := range FindColumns(field.Interface(), fields...) {
				if v != nil {
					result[i] = v
				}
			}
			continue
		}

		tag := fieldType.Tag.Get("auth")
		if index := xslice.IndexOf(fields, tag); index != -1 {
			result[index] = &Column{
				Name:  columnName(&fieldType),
				Field: &fieldType,
			}
		}
	}

	return result
}

func columnName(field *reflect.StructField) string {
	for _, t := range strings.Split(field.Tag.Get("gorm"), ";") { // Check for gorm column name override
		if strings.HasPrefix(t, "column") {
			v := strings.Split(t, ":")
			return strings.TrimSpace(v[1])
		}
	}

	return database.GetDB().NamingStrategy.ColumnName("", field.Name)
}
