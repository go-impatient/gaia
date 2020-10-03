package repository

import (
	"fmt"
	"gorm.io/gorm"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/go-impatient/gaia/internal/database"
)

type WhereParam struct {
	Field   string
	Tag     string
	Prepare interface{}
}

type QueryParam struct {
	Fields     string
	Offset     int
	Limit      int
	Order      string
	Where      []WhereParam
}

// Create 数据添加
func Create(data interface{}) bool {
	db := database.Orm().Create(data)
	if err := db.Error; err != nil {
		log.Error().Msgf("mysql query error: %v, sql[%v]", err)
		return false
	}
	return true
}

// FindMulti 数据复合查询
func FindMulti(model interface{}, query QueryParam) bool {
	db := database.Orm().Offset(query.Offset)
	if query.Limit > 0 {
		db = db.Limit(query.Limit)
	}
	if query.Fields != "" {
		db = db.Select(query.Fields)
	}
	if query.Order != "" {
		db = db.Order(query.Order)
	}
	db = parseWhereParam(db, query.Where)
	db.Find(model)
	if err := db.Error; err != nil {
		log.Error().Msgf("sql query error: %s", err.Error())
		return false
	}
	return true
}

// Count 统计某条件字段的条目数
func Count(model interface{}, count *int64, query QueryParam) bool {
	db := database.Orm().Model(model)
	db = parseWhereParam(db, query.Where)
	db = db.Count(count)
	if err := db.Error; err != nil {
		log.Error().Msgf("sql query error: %s", err.Error())
		return false
	}
	return true
}

// Delete, 根据条件和 model 的值进行批量删除
func Delete(model interface{}, query QueryParam) bool {
	if len(query.Where) == 0 {
		log.Error().Msgf("sql query error: delete failed, where conditions cannot be empty")
		return false
	}
	db := database.Orm().Model(model)
	db = parseWhereParam(db, query.Where)
	db = db.Delete(model)
	if err := db.Error; err != nil {
		log.Error().Msgf("sql query error: %s", err.Error())
		return false
	}
	return true
}

// DeleteById, 对象的主键有值，会被用于构建条件, 进行删除, 没有主键会触发批量Delete
func DeleteById(model interface{}) bool {
	db := database.Orm().Model(model)
	db.Delete(model)
	if err := db.Error; err != nil {
		log.Error().Msgf("sql query error: %s", err.Error())
		return false
	}
	return true
}

// FindOne, 根据条件和 model 的值进行查询
func FindOne(model interface{}, query QueryParam) bool {
	db := database.Orm().Model(model)
	if query.Fields != "" {
		db = db.Select(query.Fields)
	}
	db = parseWhereParam(db, query.Where)
	db = db.First(model)
	if err := db.Error; err != nil{
		log.Error().Msgf("sql query error: %s", err.Error())
		return false
	}
	return true
}

// FindById, 根据ID条件查询数据
func FindById(model interface{}, id interface{}) bool {
	db := database.Orm().Model(model)
	db.First(model, id)
	if err := db.Error; err != nil{
		log.Error().Msgf("mysql query error: %s sql[%v]", err.Error())
		return false
	}
	return true
}

// UpdataById, 对象的主键有值，会被用于构建条件, 进行更新
func UpdateById(model interface{}) bool {
	db := database.Orm().Model(model)
	db = db.Updates(model)
	if err := db.Error; err != nil {
		log.Error().Msgf("mysql query error: %s, sql[%v]", err.Error())
		return false
	}
	return true
}

// Update, 根据条件和 model 的值进行批量更新
func Update(model interface{}, data interface{}, query QueryParam) bool {
	db := database.Orm().Model(model)
	db = parseWhereParam(db, query.Where)
	db = db.Updates(data)
	if err := db.Error; err != nil {
		log.Error().Msgf("mysql query error: %s, sql[%v]", err.Error())
		return false
	}
	return true
}

// parseWhereParam, Where 条件拼装
func parseWhereParam(db *gorm.DB, where []WhereParam) *gorm.DB {
	if len(where) == 0 {
		return db
	}
	var (
		plain []string
		prepare []interface{}
	)
	for _, w := range where {
		tag := w.Tag
		if tag == "" {
			tag = "="
		}
		var plainFmt string
		switch tag {
		case "IN":
			plainFmt = fmt.Sprintf("%s IN (?)", w.Field)
		default:
			plainFmt = fmt.Sprintf("%s %s ?", w.Field, tag)
		}
		plain = append(plain, plainFmt)
		prepare = append(prepare, w.Prepare)
	}
	return db.Where(strings.Join(plain, " AND "), prepare...)
}
