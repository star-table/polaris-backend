package convert

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/star-table/common/core/util/json"
	"github.com/spf13/cast"
)

func ObjectToMap(obj interface{}) map[string]interface{} {
	jsonStr := json.ToJsonIgnoreError(obj)
	data := map[string]interface{}{}
	_ = json.FromJson(jsonStr, &data)
	return data
}

// ToIntSliceE casts an interface to a []int type.
func ToInt64SliceE(i interface{}) ([]int64, error) {
	if i == nil {
		return []int64{}, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
	}

	switch v := i.(type) {
	case []int64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToInt64E(s.Index(j).Interface())
			if err != nil {
				return []int64{}, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []int64{}, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
	}
}

// ToIntSlice casts an interface to a []int type.
func ToInt64Slice(i interface{}) []int64 {
	v, _ := ToInt64SliceE(i)
	return v
}

// unsafe
func UnsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// unsafe
func UnsafeStringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
