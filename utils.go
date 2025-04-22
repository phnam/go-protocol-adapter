package sdk

import (
	"encoding/json"
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(i interface{}) string {
	str := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.Split(str, "/")
	if len(parts) > 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	return str
}

func ConvertToObjectSlice[T any](jsonString string) []T {
	var result []T
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil
	}
	return result
}
