package clog

import (
	"encoding/json"
	"reflect"

	"go.uber.org/zap"
)

func RedactSensitive(v any) any {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return nil
		}
		return RedactSensitive(rv.Elem().Interface())

	case reflect.Struct:
		result := make(map[string]any)
		for i := range rv.NumField() {
			field := rt.Field(i)
			if !rv.Field(i).CanInterface() {
				continue
			}
			value := rv.Field(i).Interface()

			if field.Tag.Get("sensitive") == "true" {
				result[field.Name] = "<REDACTED>"
			} else {
				result[field.Name] = RedactSensitive(value)
			}
		}
		return result

	case reflect.Slice, reflect.Array:
		safeArr := make([]any, rv.Len())
		for i := range rv.Len() {
			safeArr[i] = RedactSensitive(rv.Index(i).Interface())
		}
		return safeArr

	case reflect.Map:
		safeMap := make(map[string]any)
		for _, key := range rv.MapKeys() {
			safeMap[key.String()] = RedactSensitive(rv.MapIndex(key).Interface())
		}
		return safeMap

	default:
		return v
	}
}

func MarshalSafeJSON(v any) string {
	safe := RedactSensitive(v)
	b, _ := json.Marshal(safe)
	return string(b)
}

func SafeAny(key string, value any) zap.Field {
	safeVal := RedactSensitive(value)
	return zap.Any(key, safeVal)
}

func SafeObject(key string, value any) zap.Field {
	safeVal := RedactSensitive(value)
	return zap.Reflect(key, safeVal)
}
