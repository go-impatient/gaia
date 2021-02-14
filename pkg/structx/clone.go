package xstruct

import (
	"reflect"
)

// Clone, 拷贝一份结构体
func Clone(src, dst interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		value := srcVal.Field(i)
		name := srcVal.Type().Field(i).Name

		dstValue := dstVal.FieldByName(name)
		if dstValue.IsValid() == false {
			continue
		}
		dstValue.Set(value)
	}
}
