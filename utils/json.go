package utils

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func MustJsonMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// formatStruct takes any struct and formats it as a map[string]any.
func PrepareTidyJSON(input interface{}) (any, error) {
	return formatValue(reflect.ValueOf(input)), nil
}

// formatValue handles the conversion of values based on their kind.
func formatValue(v reflect.Value) any {
	switch v.Kind() {
	case reflect.Pointer:
		// Check if the pointer is nil
		if v.IsNil() {
			return nil
		}
		// Does this pointer have a method named String() string?
		if m := v.MethodByName("String"); m.IsValid() {
			// If so, call it and return the result
			return m.Call(nil)[0].Interface()
		}

		return formatValue(v.Elem())
	case reflect.Struct:
		return formatStructValue(v)
	case reflect.Slice, reflect.Array:
		return formatSliceValue(v)
	case reflect.String:
		return formatStringValue(v.String())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64, reflect.Bool:
		return v.Interface()
	default:
		// For other kinds, return as is or add more cases as needed.
		return v.Interface()
	}
}

// formatStructValue formats struct values.
func formatStructValue(v reflect.Value) map[string]any {
	result := make(map[string]any)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		// check if field is exported
		if !field.IsExported() {
			continue
		}

		key := tidyFieldName(field)
		// fmt.Println("formatStructValue", key, v.Field(i).Type().String())
		result[key] = formatValue(v.Field(i))
		// clean up empty fields
		if result[key] == nil {
			delete(result, key)
		}
	}
	return result
}

// formatSliceValue formats slice and array values.
func formatSliceValue(v reflect.Value) any {
	// check if this a slice of bytes
	if v.Type().Elem().Kind() == reflect.Uint8 {
		return hexutil.Encode(v.Bytes())
	}

	result := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = formatValue(v.Index(i))
	}
	return result
}

// formatStringValue formats string values, converting base64 to hex if applicable.
func formatStringValue(s string) string {
	if data, err := base64.StdEncoding.DecodeString(s); err == nil {
		// Successfully decoded from base64, now convert to hex
		return hex.EncodeToString(data)
	}
	// Return the original string if it's not base64
	return s
}

func tidyFieldName(field reflect.StructField) string {
	// if json tag is present, use that
	if tag, ok := field.Tag.Lookup("json"); ok {
		return strings.Split(tag, ",")[0]
	}
	// otherwise return name in camel case
	if len(field.Name) > 0 {
		return strings.ToLower(field.Name[:1]) + field.Name[1:]
	}

	return strings.ToLower(field.Name)
}
