package keyvalue

import "reflect"

func hasZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}

	t := reflect.TypeOf(v)
	if t == nil {
		return true
	}

	switch t.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return false
	}

	return v == reflect.Zero(t).Interface()
}

func isKeyValue(v interface{}) bool {
	if v == nil {
		return false
	}

	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}

	return t.Kind() == reflect.Map || t.Name() == "KeyValue"
}

// IsAbleToConvert check whether an interface could be able to cast to KeyValue
func IsAbleToConvert(p interface{}) bool {
	t := reflect.TypeOf(p)
	kind := t.Kind()

	switch kind {
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		return true
	}

	return true
}
