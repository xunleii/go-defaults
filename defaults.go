package defaults

import (
	"reflect"
)

// Func is a function that fill any struct field from the tag value.
type Func func(field *FieldData)

// FieldData contains the value and the tag value of a struct field.
type FieldData struct {
	Value    reflect.Value
	TagValue string
}

// SetDefaults applies the default values to the struct object, the struct type must have
// the StructTag with name "default" and the directed value.
//
// Usage
//     type ExampleBasic struct {
//         Foo bool   `default:"true"`
//         Bar string `default:"33"`
//         Qux int8
//     }
//
//      foo := &ExampleBasic{}
//      SetDefaults(foo)
func SetDefaults(v interface{}) {
	value := reflect.ValueOf(v).Elem()
	fields := fieldsFromValue(value)

	for _, field := range fields {
		if isEmpty(field) {
			setDefaults(field)
		}
	}
}

// RegisterCustomDefault bind a custom Func with a specific type allowing to set default value
// for user defined type.
func RegisterCustomDefault(p reflect.Type, fnc Func) { defaultWithType[p.String()] = fnc }

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
				_, exists := defaultWithType[subValue.Type().String()]
				if exists {
					fields = append(fields, &FieldData{
						Value:    subValue,
						TagValue: tag,
					})
					continue
				}

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
	case reflect.Complex64, reflect.Complex128:
		return field.Value.Complex() == 0i
	case reflect.String:
		return field.Value.String() == ""
	case reflect.Array, reflect.Slice, reflect.Map:
		return field.Value.Len() == 0
	case reflect.Chan, reflect.Interface, reflect.Ptr, reflect.Struct, reflect.Func:
		return field.Value.IsNil()
	default:
		return false
	}
}

var defaultWithKind = map[reflect.Kind]Func{
	reflect.Bool:    defaultBool,
	reflect.Int:     defaultInt,
	reflect.Int8:    defaultInt,
	reflect.Int16:   defaultInt,
	reflect.Int32:   defaultInt,
	reflect.Int64:   defaultInt,
	reflect.Uint:    defaultUint,
	reflect.Uint8:   defaultUint,
	reflect.Uint16:  defaultUint,
	reflect.Uint32:  defaultUint,
	reflect.Uint64:  defaultUint,
	reflect.Float32: defaultFloat,
	reflect.Float64: defaultFloat,
	reflect.String:  defaultString,
}

var defaultWithType = map[string]Func{
	"time.Duration": defaultDuration,
	"[]uint8":       defaultBytes,
}
