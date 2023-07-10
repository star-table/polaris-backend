package mvc
import "reflect"

func isBodyFlag(kind reflect.Kind) bool {
	return kind == reflect.Struct || kind == reflect.Map || kind == reflect.Array || kind == reflect.Slice
}
