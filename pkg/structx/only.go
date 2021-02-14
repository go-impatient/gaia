package xstruct

import (
	"fmt"
	"reflect"

	"github.com/go-impatient/gaia/pkg/xslice"
)

// Only extracts the requested field from the given map[string] or structure and
// returns a map[string]interface{} containing only those values.
//
// For example:
//  type Model struct {
//    Field string
//    Num   int
//    Slice []float64
//  }
//  model := Model{
// 	  Field: "value",
// 	  Num:   42,
// 	  Slice: []float64{3, 6, 9},
//  }
//  res := Only(model, "Field", "Slice")
//
// Result:
//  map[string]interface{}{
// 	  "Field": "value",
// 	  "Slice": []float64{3, 6, 9},
//  }
//
// In case of conflicting fields (if a promoted field has the same name as a parent's
// struct field), the higher level field is kept.
func Only(data interface{}, fields ...string) map[string]interface{} {
	result := make(map[string]interface{}, len(fields))
	t := reflect.TypeOf(data)
	value := reflect.ValueOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		value = value.Elem()
	}

	switch t.Kind() {
	case reflect.Map:
		if t.Key().Kind() != reflect.String {
			panic(fmt.Errorf("helper.Only only supports map[string] and structures, %s given", t.String()))
		}
		for _, k := range value.MapKeys() {
			name := k.String()
			if xslice.ContainsStr(fields, name) {
				result[name] = value.MapIndex(k).Interface()
			}
		}
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			field := value.Field(i)
			fieldType := t.Field(i)
			name := fieldType.Name
			if field.Kind() == reflect.Struct && fieldType.Anonymous {
				for k, v := range Only(field.Interface(), fields...) {
					// Check if fields are conflicting
					// Highest level fields have priority
					if _, ok := result[k]; !ok {
						result[k] = v
					}
				}
			} else if xslice.ContainsStr(fields, name) {
				result[name] = value.Field(i).Interface()
			}
		}
	default:
		panic(fmt.Errorf("only supports map[string] and structures, %s given", t.Kind()))
	}

	return result
}
