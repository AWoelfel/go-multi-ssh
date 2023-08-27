package assert

import (
	"reflect"
	"testing"
)

func AssertObjectsEqual(t *testing.T, obj1, obj2 interface{}) {
	value1 := reflect.ValueOf(obj1)
	value2 := reflect.ValueOf(obj2)

	if value1.Type() != value2.Type() {
		t.Errorf("Type mismatch: %v != %v", value1.Type(), value2.Type())
		return
	}

	switch value1.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value1.Int() != value2.Int() {
			t.Errorf("Int mismatch: %v != %v", value1.Int(), value2.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if value1.Uint() != value2.Uint() {
			t.Errorf("Uint mismatch: %v != %v", value1.Uint(), value2.Uint())
		}
	case reflect.Float32, reflect.Float64:
		if value1.Float() != value2.Float() {
			t.Errorf("Float mismatch: %v != %v", value1.Float(), value2.Float())
		}
	case reflect.Bool:
		if value1.Bool() != value2.Bool() {
			t.Errorf("Bool mismatch: %v != %v", value1.Bool(), value2.Bool())
		}
	case reflect.String:
		if value1.String() != value2.String() {
			t.Errorf("String mismatch: %v != %v", value1.String(), value2.String())
		}
	case reflect.Slice:
		if value1.Len() != value2.Len() {
			t.Errorf("Slice length mismatch: %v != %v", value1.Len(), value2.Len())
			return
		}
		for i := 0; i < value1.Len(); i++ {
			AssertObjectsEqual(t, value1.Index(i).Interface(), value2.Index(i).Interface())
		}
	case reflect.Struct:
		for i := 0; i < value1.NumField(); i++ {
			AssertObjectsEqual(t, value1.Field(i).Interface(), value2.Field(i).Interface())
		}
	case reflect.Pointer:
		if value1.Pointer() != value2.Pointer() {
			t.Errorf("Pointer mismatch: %v != %v", value1.Pointer(), value2.Pointer())
			return
		}
	default:
		t.Errorf("Unsupported type: %v", value1.Kind())
	}
}
