package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ExitIfNotNil exits the process with status 1 and error with the error
// passed to the function is not nil
func ExitIfNotNil(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

// MapGet takes in the map keys - each key goes one level deeper in the map
func MapGet(mp map[string]interface{}, key ...string) (interface{}, bool) {
	if mp == nil {
		return nil, false
	}

	if len(key) == 0 {
		return mp, true
	}

	if len(key) == 1 {
		val, ok := mp[key[0]]
		return val, ok
	}

	val, ok := mp[key[0]]
	if !ok {
		return mp, false
	}

	switch v := val.(type) {
	case map[string]interface{}:
		return MapGet(v, key[1:]...)
	case []interface{}:
		// Check if we can find key in the nested structure
		if len(key) < 2 {
			return mp, false
		}

		// Check if the key[1] is of type uint, if it is then
		keyNum, err := strconv.Atoi(key[1])
		if err != nil {
			return mp, false
		}

		if keyNum >= len(v) {
			return mp, false
		}

		valMapM, ok := v[keyNum].(map[string]interface{})
		if !ok {
			return mp, false
		}

		return MapGet(valMapM, key[2:]...)
	case []map[string]interface{}:
		// Check if we can find key in the nested structure
		if len(key) < 2 {
			return mp, false
		}

		// Check if the key[1] is of type uint, if it is then
		keyNum, err := strconv.Atoi(key[1])
		if err != nil {
			return mp, false
		}

		if keyNum >= len(v) {
			return mp, false
		}

		return MapGet(v[keyNum], key[2:]...)
	}

	return mp, true
}

// MapSet takes in the map that needs to be manipulated, the value that needs to
// be assgined to be assigned and the key - each key goes one level deeper in the map
func MapSet(mp map[string]interface{}, value interface{}, key ...string) {
	var _mapSet func(map[string]interface{}, interface{}, ...string) map[string]interface{}

	_mapSet = func(mp map[string]interface{}, value interface{}, key ...string) map[string]interface{} {
		if mp == nil {
			return nil
		}

		if len(key) == 0 {
			return mp
		}

		if len(key) == 1 {
			mp[key[0]] = value
			return mp
		}

		val, ok := mp[key[0]]
		if !ok {
			internal := map[string]interface{}{}

			mp[key[0]] = internal
			val = internal
		}

		switch v := val.(type) {
		case map[string]interface{}:
			mp[key[0]] = _mapSet(v, value, key[1:]...)
			return mp
		case []interface{}:
			// Check if we can find key in the nested structure
			if len(key) < 2 {
				return mp
			}

			// Check if the key[1] is of type uint, if it is then
			keyNum, err := strconv.Atoi(key[1])
			if err != nil {
				return mp
			}

			if keyNum >= len(v) {
				return mp
			}

			valMapM, ok := v[keyNum].(map[string]interface{})
			if !ok {
				return mp
			}

			v[keyNum] = _mapSet(valMapM, value, key[2:]...)

			mp[key[0]] = v

			return mp
		case []map[string]interface{}:
			// Check if we can find key in the nested structure
			if len(key) < 2 {
				return mp
			}

			// Check if the key[1] is of type uint, if it is then
			keyNum, err := strconv.Atoi(key[1])
			if err != nil {
				return mp
			}

			if keyNum >= len(v) {
				return mp
			}

			v[keyNum] = _mapSet(v[keyNum], value, key[2:]...)

			mp[key[0]] = v

			return mp
		}

		return mp
	}

	_mapSet(mp, value, key...)
}

func MapFromYAML(yamlByt []byte) (map[string]interface{}, error) {
	var data map[string]interface{}

	err := yaml.Unmarshal(yamlByt, &data)
	if err != nil {
		return nil, err
	}

	return RecursiveCastMapStringInterfaceToMapStringInterface(data), nil
}

// RecursiveCastMapStringInterfaceToMapStringInterface will convert a
// map[string]interface{} recursively => map[string]interface{}
func RecursiveCastMapStringInterfaceToMapStringInterface(in map[string]interface{}) map[string]interface{} {
	res := ConvertMapInterfaceMapString(in)
	out, ok := res.(map[string]interface{})
	if !ok {
		logrus.Debug("failed to cast")
	}

	return out
}

// ConvertMapInterfaceMapString converts map[interface{}]interface{} => map[string]interface{}
//
// It will also convert []interface{} => []string
func ConvertMapInterfaceMapString(v interface{}) interface{} {
	switch x := v.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v2 := range x {
			switch k2 := k.(type) {
			case string:
				m[k2] = ConvertMapInterfaceMapString(v2)
			default:
				m[fmt.Sprint(k)] = ConvertMapInterfaceMapString(v2)
			}
		}
		v = m
	case []interface{}:
		for i, v2 := range x {
			x[i] = ConvertMapInterfaceMapString(v2)
		}
	case map[string]interface{}:
		for k, v2 := range x {
			x[k] = ConvertMapInterfaceMapString(v2)
		}
	}

	return v
}

// Any takes any data and converts it into another data
//
// Internally it converts typ1 to json byt and then convert it back
// to typ2 hence it is RREQUIRED for both of the types to be encodable
// in JSON or else the function will return error
func AnyToAny(typ1, typ2 interface{}) error {
	byt, err := json.Marshal(typ1)
	if err != nil {
		return errors.New("failed to convert")
	}

	if err := json.Unmarshal(byt, &typ2); err != nil {
		return errors.New("failed to convert")
	}

	return nil
}
