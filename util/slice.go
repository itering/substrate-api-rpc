package util

import (
	"encoding/json"
	"fmt"
)

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func PrintJson(j interface{}) {
	var val string
	switch i := j.(type) {
	case string:
		val = i
	case []byte:
		val = string(i)
	case error:
		val = i.Error()
	default:
		b, _ := json.MarshalIndent(i, "", "  ")
		val = string(b)
	}
	fmt.Println(val)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
