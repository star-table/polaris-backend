package mvc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/util/str"
)

type RequestInfo struct {
	Parameters map[string]string
	Body       string
	Ctx        context.Context
}

var BasicTypeConverter = map[string]func(v string) (interface{}, error){}

func init() {
	intFunc := func(bitSize int) func(v string) (interface{}, error) {
		return func(v string) (interface{}, error) {
			obj, err := strconv.ParseInt(v, 10, bitSize)
			if err != nil {
				return 0, err
			}
			switch bitSize {
			case 0:
				obj := int(obj)
				return &obj, nil
			case 8:
				obj := int8(obj)
				return &obj, nil
			case 16:
				obj := int16(obj)
				return &obj, nil
			case 32:
				obj := int32(obj)
				return &obj, nil
			}
			return &obj, nil
		}
	}
	floatFunc := func(bitSize int) func(v string) (interface{}, error) {
		return func(v string) (interface{}, error) {
			obj, err := strconv.ParseFloat(v, bitSize)
			if err != nil {
				return 0, err
			}
			switch bitSize {
			case 32:
				obj := float32(obj)
				return &obj, nil
			}
			return &obj, nil
		}
	}
	uintFunc := func(bitSize int) func(v string) (interface{}, error) {
		return func(v string) (interface{}, error) {
			obj, err := strconv.ParseUint(v, 10, bitSize)
			if err != nil {
				return 0, err
			}
			switch bitSize {
			case 0:
				obj := uint(obj)
				return &obj, nil
			case 8:
				obj := uint8(obj)
				return &obj, nil
			case 16:
				obj := uint16(obj)
				return &obj, nil
			case 32:
				obj := uint32(obj)
				return &obj, nil
			}
			return &obj, nil
		}
	}
	BasicTypeConverter[reflect.Bool.String()] = func(v string) (interface{}, error) {
		b, e := strconv.ParseBool(v)
		return &b, e
	}

	BasicTypeConverter[reflect.Int.String()] = intFunc(0)
	BasicTypeConverter[reflect.Int8.String()] = intFunc(8)
	BasicTypeConverter[reflect.Int16.String()] = intFunc(16)
	BasicTypeConverter[reflect.Int32.String()] = intFunc(32)
	BasicTypeConverter[reflect.Int64.String()] = intFunc(64)
	BasicTypeConverter[reflect.Uint.String()] = uintFunc(0)
	BasicTypeConverter[reflect.Uint8.String()] = uintFunc(8)
	BasicTypeConverter[reflect.Uint16.String()] = uintFunc(16)
	BasicTypeConverter[reflect.Uint32.String()] = uintFunc(32)
	BasicTypeConverter[reflect.Uint64.String()] = uintFunc(64)
	BasicTypeConverter[reflect.Float32.String()] = floatFunc(32)
	BasicTypeConverter[reflect.Float64.String()] = floatFunc(64)
	BasicTypeConverter[reflect.String.String()] = func(v string) (interface{}, error) {
		return &v, nil
	}
	BasicTypeConverter["date.Time"] = func(v string) (interface{}, error) {
		t, e := time.Parse(consts.AppTimeFormat, v)
		return &t, e
	}
	BasicTypeConverter["time.Time"] = func(v string) (interface{}, error) {
		t, e := time.Parse(consts.AppTimeFormat, v)
		return &t, e
	}
	BasicTypeConverter["types.Time"] = func(v string) (interface{}, error) {
		t, err := time.Parse(consts.AppTimeFormat, v)
		tt := types.Time(t)
		return &tt, err
	}
}

func InjectFunc(targetFunc interface{}, reqInfo RequestInfo) ([]reflect.Value, error) {
	targetType := reflect.TypeOf(targetFunc)
	if reflect.Func != targetType.Kind() {
		return nil, errors.New("target is not func")
	}
	numIn := targetType.NumIn()
	inputValues := make([]reflect.Value, numIn)
	if numIn > 0 {
		for i := 0; i < numIn; i++ {
			elem := targetType.In(i)
			isPtr := false
			//if elem.Kind() == reflect.Ptr {
			//	elem = elem.Elem()
			//	isPtr = true
			//}
			judgeIsPtr(&isPtr, &elem)
			//if elem.String() == "context.Context" {
			//	if isPtr {
			//		inputValues[i] = reflect.ValueOf(&reqInfo.Ctx)
			//	} else {
			//		inputValues[i] = reflect.ValueOf(reqInfo.Ctx)
			//	}
			//	continue
			//}
			if assemblyInputValues(&inputValues, elem, isPtr, reqInfo, i) {
				continue
			}
			if elem.Kind() == reflect.Struct {
				value, err := ParseValue(elem, isPtr, reqInfo)
				if err != nil {
					log.Error(fmt.Sprintf("%#v", err))
					return nil, err
				}
				inputValues[i] = *value
			}
		}
	}
	return reflect.ValueOf(targetFunc).Call(inputValues), nil
}

//func InjectFunc1(targetFunc interface{}, reqInfo RequestInfo) ([]reflect.Value, error) {
//	targetType := reflect.TypeOf(targetFunc)
//	if reflect.Func != targetType.Kind() {
//		return nil, errors.New("target is not func")
//	}
//	numIn := targetType.NumIn()
//	inputValues := make([]reflect.Value, numIn)
//	if numIn > 0 {
//		for i := 0; i < numIn; i++ {
//			elem := targetType.In(i)
//			isPtr := false
//			if elem.Kind() == reflect.Ptr {
//				elem = elem.Elem()
//				isPtr = true
//			}
//			//judgeIsPtr(&isPtr,&elem)
//			if elem.String() == "context.Context" {
//				if isPtr {
//					inputValues[i] = reflect.ValueOf(&reqInfo.Ctx)
//				} else {
//					inputValues[i] = reflect.ValueOf(reqInfo.Ctx)
//				}
//				continue
//			}
//			//if assemblyInputValues(&inputValues,elem,isPtr,reqInfo,i){
//			//	continue
//			//}
//			if elem.Kind() == reflect.Struct {
//				value, err := ParseValue(elem, isPtr, reqInfo)
//				if err != nil {
//					log.Error(err)
//					return nil, err
//				}
//				inputValues[i] = *value
//			}
//		}
//	}
//	return reflect.ValueOf(targetFunc).Call(inputValues), nil
//}

func judgeIsPtr(isPtr *bool, elem *reflect.Type) {
	if (*elem).Kind() == reflect.Ptr {
		*elem = (*elem).Elem()
		*isPtr = true
	}
}

func assemblyInputValues(inputValues *[]reflect.Value, elem reflect.Type, isPtr bool, reqInfo RequestInfo, i int) bool {
	if elem.String() == "context.Context" {
		if isPtr {
			(*inputValues)[i] = reflect.ValueOf(&reqInfo.Ctx)
		} else {
			(*inputValues)[i] = reflect.ValueOf(reqInfo.Ctx)
		}
		return true
	}
	return false
}

func ParseValue(elem reflect.Type, isPtr bool, reqInfo RequestInfo) (*reflect.Value, error) {
	reqObj := reflect.New(elem).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := field.Type
		isPtr := false
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
			isPtr = true
		}
		var target *reflect.Value = nil
		var err error = nil

		if converter, ok := BasicTypeConverter[fieldType.String()]; ok {
			target, err = ParseQuery(getFieldName(field), isPtr, converter, reqInfo)
		} else if isBodyFlag(fieldType.Kind()) {
			target, err = ParseBody(fieldType, isPtr, reqInfo)
		}
		if err != nil {
			log.Error(fmt.Sprintf("%#v", err))
			return nil, err
		}

		if target != nil {
			reqObj.FieldByName(field.Name).Set(*target)
		}
	}
	if isPtr {
		reqObj = reqObj.Addr()
	}
	return &reqObj, nil
}

func getFieldName(field reflect.StructField) string {
	//name := field.Tag.Get("json")
	//if name == "" {
	//	return field.Name
	//}
	//return strings.Split(name, ",")[0]
	return field.Name
}

func ParseBody(elem reflect.Type, isPtr bool, reqInfo RequestInfo) (*reflect.Value, error) {
	body := reqInfo.Body
	if body == "" && isPtr {
		return nil, nil
	}

	newStrut := reflect.New(elem)
	targetInterface := newStrut.Interface()
	err := json.FromJson(body, &targetInterface)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if !isPtr {
		newStrut = newStrut.Elem()
	}
	return &newStrut, nil
}

func ParseQuery(fieldName string, isPtr bool, converter func(v string) (interface{}, error), reqInfo RequestInfo) (*reflect.Value, error) {
	paramValue := reqInfo.Parameters[str.LcFirst(fieldName)]
	//指针直接返回空
	if paramValue == "" && isPtr {
		return nil, nil
	}
	v, err := converter(paramValue)
	if err != nil {
		log.Error(fmt.Sprintf("fieldName: %s, err: %#v", fieldName, err))
		return nil, err
	}
	va := reflect.ValueOf(v)
	if !isPtr {
		va = va.Elem()
	}
	return &va, nil
}
