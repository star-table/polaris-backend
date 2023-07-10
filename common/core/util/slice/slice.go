package slice

import "reflect"

func IsSlice(arr interface{}) bool {
	if arr == nil {
		return true
	}
	v := reflect.ValueOf(arr)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return false
	}

	return true
}
