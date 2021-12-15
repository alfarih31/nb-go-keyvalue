package keyvalue

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type KeyValue map[string]interface{}

func (k KeyValue) JSON() map[string]string {
	o := map[string]string{}
	for key, val := range k {
		o[key] = fmt.Sprintf("%v", val)
	}

	return o
}

func (k KeyValue) ToMap() map[string]interface{} {
	o := map[string]interface{}{}
	for key, val := range k {
		o[key] = val
	}

	return o
}

func (k KeyValue) AssignTo(target KeyValue, replaceExist ...bool) {
	rExist := false
	if len(replaceExist) > 0 {
		rExist = replaceExist[0]
	}

	for key, val := range k {
		targetValue, exist := target[key]

		// Recursive assignTo
		if reflect.ValueOf(val).Kind() == reflect.Map && reflect.ValueOf(targetValue).Kind() == reflect.Map {
			sourceKvVal, _ := FromStruct(val)
			targetKvVal, _ := FromStruct(targetValue)

			sourceKvVal.AssignTo(targetKvVal, rExist)

			target[key] = targetKvVal
			return
		}

		// Check is target is zero value
		t := reflect.TypeOf(targetValue)
		isZero := true
		if t != nil {
			// If kind is `slice` or `Array` or `Map` then always decide as zero
			if t.Kind() == reflect.Slice || t.Kind() == reflect.Array || t.Kind() == reflect.Map {
				isZero = false
			} else {
				isZero = targetValue == reflect.Zero(t).Interface()
			}
		}

		// If exist but don't replace exist then continue
		if exist && !rExist && !isZero {
			continue
		}

		target[key] = val
	}
}

func (k KeyValue) Assign(source KeyValue, replaceExist ...bool) {
	rExist := false
	if len(replaceExist) > 0 {
		rExist = replaceExist[0]
	}

	for key, val := range source {
		existingValue, exist := k[key]

		// Recursive assign
		if reflect.ValueOf(val).Kind() == reflect.Map && reflect.ValueOf(existingValue).Kind() == reflect.Map {
			sourceKvVal, _ := FromStruct(val)
			existingKvVal, _ := FromStruct(existingValue)

			existingKvVal.Assign(sourceKvVal, rExist)

			k[key] = existingKvVal
			return
		}

		// Check existing is zero value
		t := reflect.TypeOf(existingValue)
		isZero := true
		if t != nil {
			// If kind is `slice` or `Array` or `Map` then always decide as zero
			if t.Kind() == reflect.Slice || t.Kind() == reflect.Array || t.Kind() == reflect.Map {
				isZero = false
			} else {
				isZero = existingValue == reflect.Zero(t).Interface()
			}
		}

		// If exist & not zero value but don't replace exist then continue
		if exist && !rExist && !isZero {
			continue
		}

		k[key] = val
	}
}

func (k KeyValue) Keys() []string {
	var keys []string
	for key := range k {
		keys = append(keys, key)
	}
	return keys
}

func (k KeyValue) Values() []interface{} {
	var values []interface{}
	for _, val := range k {
		values = append(values, val)
	}

	return values
}

func (k KeyValue) String() string {
	j, _ := json.Marshal(k)
	return string(j)
}

func IsAbleToConvert(p interface{}) bool {
	t := reflect.TypeOf(p)
	name := t.Name()
	kind := t.Kind()

	if name == "KeyValue" {
		return true
	}

	switch kind {
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		return true
	}

	return true
}

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

func FromStruct(strct interface{}) (KeyValue, error) {
	if !IsAbleToConvert(strct) {
		return nil, errors.New("cannot convert")
	}

	mapString, err := structToMap(strct)

	if err != nil {
		return nil, err
	}

	kv := KeyValue{}
	for key, val := range mapString {
		kv[key] = val
	}

	return kv, nil
}
