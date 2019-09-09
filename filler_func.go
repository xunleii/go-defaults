package defaults

import (
	"strconv"
	"time"
)

func boolFiller(field *FieldData) {
	value := false
	if field.TagValue == "true" {
		value = true
	}
	field.Value.SetBool(value)
}

func stringFiller(field *FieldData) {
	field.Value.SetString(field.TagValue)
}

func intFiller(field *FieldData) {
	value, _ := strconv.ParseInt(field.TagValue, 10, 64)
	field.Value.SetInt(value)
}

func uintFiller(field *FieldData) {
	value, _ := strconv.ParseUint(field.TagValue, 10, 64)
	field.Value.SetUint(value)
}

func floatFiller(field *FieldData) {
	value, _ := strconv.ParseFloat(field.TagValue, 64)
	field.Value.SetFloat(value)
}

func durationFiller(field *FieldData) {
	value, _ := time.ParseDuration(field.TagValue)
	field.Value.SetInt(int64(value))
}

func ignoreFiller(*FieldData) {}
