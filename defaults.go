package defaults

import (
	"reflect"
)

type Func func(field *FieldData)
type FieldData struct {
	Value    reflect.Value
	TagValue string
}

func SetDefaults(v interface{}) {
	value := reflect.ValueOf(v).Elem()
	fields := fieldsFromValue(value)

	for _, field := range fields {
		if isEmpty(field) {
			setDefaults(field)
		}
	}
}

func AddCustomDefault(p reflect.Type, fnc Func) {
	defaultWithType[p.String()] = fnc
}

func fieldsFromValue(value reflect.Value) []*FieldData {
	vtype := value.Type()

	switch vtype.Kind() {
	case reflect.Ptr:
		return fieldsFromValue(value.Elem())
	case reflect.Array, reflect.Slice:
		var fields []*FieldData

		for i := 0; i < value.Len(); i++ {
			subFields := fieldsFromValue(value.Index(i))
			if subFields != nil {
				fields = append(fields, subFields...)
			}
		}
		return fields
	case reflect.Map:
		var fields []*FieldData

		for _, key := range value.MapKeys() {
			subFields := fieldsFromValue(value.MapIndex(key))
			if subFields != nil {
				fields = append(fields, subFields...)
			}
		}
		return fields
	case reflect.Struct:
		var fields []*FieldData

		for i := 0; i < value.NumField(); i++ {
			subValue := value.Field(i)
			tag := vtype.Field(i).Tag.Get("default")

			switch subValue.Kind() {
			case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
				subFields := fieldsFromValue(subValue)
				if subFields != nil {
					fields = append(fields, subFields...)
				}
			case reflect.Bool,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64,
				reflect.String:
				if tag != "" && subValue.CanSet() {
					fields = append(fields, &FieldData{
						Value:    subValue,
						TagValue: tag,
					})
				}
			}
		}
		return fields
	default:
		return nil
	}
}

func setDefaults(field *FieldData) {
	fnc, exists := defaultWithType[field.Value.Type().String()]
	if exists {
		fnc(field)
		return
	}

	fnc, exists = defaultWithKind[field.Value.Kind()]
	if exists {
		fnc(field)
		return
	}
}

func isEmpty(field *FieldData) bool {
	switch field.Value.Kind() {
	case reflect.Bool:
		return !field.Value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return field.Value.Float() == .0
	case reflect.String:
		return field.Value.String() == ""
	default:
		return false
	}
}

var defaultWithKind = map[reflect.Kind]Func{
	reflect.Invalid:       ignoreFiller,
	reflect.Bool:          boolFiller,
	reflect.Int:           intFiller,
	reflect.Int8:          intFiller,
	reflect.Int16:         intFiller,
	reflect.Int32:         intFiller,
	reflect.Int64:         intFiller,
	reflect.Uint:          uintFiller,
	reflect.Uint8:         uintFiller,
	reflect.Uint16:        uintFiller,
	reflect.Uint32:        uintFiller,
	reflect.Uint64:        uintFiller,
	reflect.Uintptr:       ignoreFiller,
	reflect.Float32:       floatFiller,
	reflect.Float64:       floatFiller,
	reflect.Complex64:     ignoreFiller,
	reflect.Complex128:    ignoreFiller,
	reflect.Array:         ignoreFiller,
	reflect.Chan:          ignoreFiller,
	reflect.Func:          ignoreFiller,
	reflect.Interface:     ignoreFiller,
	reflect.Map:           ignoreFiller,
	reflect.Ptr:           ignoreFiller,
	reflect.Slice:         ignoreFiller,
	reflect.String:        stringFiller,
	reflect.Struct:        ignoreFiller,
	reflect.UnsafePointer: ignoreFiller,
}
// todo: check custom structure like time.Time
var defaultWithType = map[string]Func{
	"time.Duration": durationFiller,
}
