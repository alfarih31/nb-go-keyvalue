// Package keyvalue provide an ability to working with Key-Value pair in golang `map[string]interface{} types`
package keyvalue

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type KeyValue map[string]interface{}

// JSON return `map[string]string` as represent of JSON object
func (k KeyValue) JSON() map[string]string {
	o := map[string]string{}
	for key, val := range k {
		o[key] = fmt.Sprintf("%v", val)
	}

	return o
}

// ToMap return `map[string]interface{}` as purpose to cast back KeyValue
func (k KeyValue) ToMap() map[string]interface{} {
	o := map[string]interface{}{}
	for key, val := range k {
		o[key] = val
	}

	return o
}

// AssignTo will assign Value from a KeyValue to other KeyValue with corresponding `Key`. This method is an opposite of Assign and return value by `in-place` operation.
func (k KeyValue) AssignTo(target KeyValue, replaceExist ...bool) {
	rExist := false
	if len(replaceExist) > 0 {
		rExist = replaceExist[0]
	}

	for key, val := range k {
		targetValue, exist := target[key]

		// Recursive assign
		vval := reflect.TypeOf(val)
		vtval := reflect.TypeOf(targetValue)

		if vval != nil && vtval != nil {
			if vval.Kind() == reflect.Map && vtval.Kind() == reflect.Map || (vval.Name() == "KeyValue" && vtval.Name() == "KeyValue") {
				sourceKvVal, _ := FromStruct(val)
				targetKvVal, _ := FromStruct(targetValue)

				sourceKvVal.AssignTo(targetKvVal, rExist)

				target[key] = targetKvVal
				continue
			}
		}

		// Check is target is zero value
		isZero := true
		if vtval != nil {
			// If kind is `slice` or `Array` or `Map` then always decide as not zero
			if vtval.Kind() == reflect.Slice || vtval.Kind() == reflect.Array || vtval.Kind() == reflect.Map {
				isZero = false
			} else {
				isZero = targetValue == reflect.Zero(vtval).Interface()
			}
		}

		// If exist but don't replace exist then continue
		if exist && !rExist && !isZero {
			continue
		}

		target[key] = val
	}
}

// Assign will assign Value by other KeyValue. This method inspired by JavaScript `Object.assign()` method
// Reference: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/assign
func (k KeyValue) Assign(source KeyValue, replaceExist ...bool) {
	rExist := false
	if len(replaceExist) > 0 {
		rExist = replaceExist[0]
	}

	for key, val := range source {
		existingValue, exist := k[key]

		// Recursive assign
		vval := reflect.TypeOf(val)
		veval := reflect.TypeOf(existingValue)

		if vval != nil && veval != nil {
			if vval.Kind() == reflect.Map && veval.Kind() == reflect.Map || (vval.Name() == "KeyValue" && veval.Name() == "KeyValue") {
				sourceKvVal, _ := FromStruct(val)
				existingKvVal, _ := FromStruct(existingValue)

				existingKvVal.Assign(sourceKvVal, rExist)

				k[key] = existingKvVal
				continue
			}
		}

		isZero := true
		if veval != nil {
			// If kind is `slice` or `Array` or `Map` then always decide as non zero
			if veval.Kind() == reflect.Slice || veval.Kind() == reflect.Array || veval.Kind() == reflect.Map || veval.Name() == "KeyValue" {
				isZero = false
			} else {
				isZero = existingValue == reflect.Zero(veval).Interface()
			}
		}

		// If exist & not zero value but don't replace exist then continue
		if exist && !rExist && !isZero {
			continue
		}

		k[key] = val
	}
}

// Keys return Array of Keys of KeyValue
func (k KeyValue) Keys() []string {
	var keys []string
	for key := range k {
		keys = append(keys, key)
	}
	return keys
}

// Values return Array of Values of KeyValue
func (k KeyValue) Values() []interface{} {
	var values []interface{}
	for _, val := range k {
		values = append(values, val)
	}

	return values
}

// String method will format the KeyValue in JSON format when calling such methods `string(KeyValue)` or fmt.Println(KeyValue)`
func (k KeyValue) String() string {
	j, _ := json.Marshal(k)
	return string(j)
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

// Unmarshal KeyValue to a Struct
func (k KeyValue) Unmarshal(i interface{}) error {
	err := json.Unmarshal([]byte(k.String()), i)

	return err
}

func structToMap(strct interface{}) (map[string]interface{}, error) {
	j, e := json.Marshal(strct)

	if e != nil {
		return nil, e
	}

	var t map[string]interface{}

	e = json.Unmarshal(j, &t)

	if e != nil {
		return nil, e
	}

	return t, nil
}

// FromStruct create a KeyValue from Struct
func FromStruct(strct interface{}) (KeyValue, error) {
	if reflect.TypeOf(strct).Name() == "KeyValue" {
		return strct.(KeyValue), nil
	}

	if !IsAbleToConvert(strct) {
		return nil, errors.New("cannot convert")
	}

	mapString, err := structToMap(strct)

	if err != nil {
		return nil, err
	}

	kv := KeyValue{}
	for key, val := range mapString {
		if val != nil {
			t := reflect.TypeOf(val)
			if t.Kind() == reflect.Map {
				kv[key], err = FromStruct(val)
				if err != nil {
					return nil, err
				}
				continue
			}
		}

		kv[key] = val
	}

	return kv, nil
}
