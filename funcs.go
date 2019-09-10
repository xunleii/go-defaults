package defaults

import (
	"reflect"
	"strconv"
	"time"
)

func defaultBool(field *FieldData) {
	value := false
	if field.TagValue == "true" {
		value = true
	}
	field.Value.SetBool(value)
}

func defaultInt(field *FieldData) {
	value, _ := strconv.ParseInt(field.TagValue, 10, 64)
	field.Value.SetInt(value)
}

func defaultUint(field *FieldData) {
	value, _ := strconv.ParseUint(field.TagValue, 10, 64)
	field.Value.SetUint(value)
}

func defaultFloat(field *FieldData) {
	value, _ := strconv.ParseFloat(field.TagValue, 64)
	field.Value.SetFloat(value)
}

func defaultDuration(field *FieldData) {
	value, _ := time.ParseDuration(field.TagValue)
	field.Value.SetInt(int64(value))
}

func defaultString(field *FieldData) { field.Value.SetString(field.TagValue) }
func defaultBytes(field *FieldData)  { field.Value.Set(reflect.ValueOf([]byte(field.TagValue))) }
