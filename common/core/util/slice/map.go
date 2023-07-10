package slice

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/db/mysql"
)

type JsonDiff struct {
	HasDiff   bool
	Result    string
	ResultMap []ChangeField
}

type ChangeField struct {
	Field    string
	OldVaule interface{}
	NewValue interface{}
}

func JsonCompare(left, right map[string]interface{}) JsonDiff {
	diff := &JsonDiff{HasDiff: false, Result: ""}
	jsonDiffDict(left, right, 1, diff)
	return *diff
}

func marshal(j interface{}) string {
	value, _ := json.Marshal(j)
	return string(value)
}

func jsonDiffDict(json1, json2 map[string]interface{}, depth int, diff *JsonDiff) {
	blank := strings.Repeat(" ", (2 * (depth - 1)))
	longBlank := strings.Repeat(" ", (2 * (depth)))
	diff.Result = diff.Result + "\n" + blank + "{"
	for key, value := range json1 {
		quotedKey := fmt.Sprintf("\"%s\"", key)
		if _, ok := json2[key]; ok {
			switch value.(type) {
			case map[string]interface{}:
				if _, ok2 := json2[key].(map[string]interface{}); !ok2 {
					assemblyDiffDictResult(diff, quotedKey, blank, key, json2, value)
				} else {
					diff.Result = diff.Result + "\n" + longBlank + quotedKey + ": "
					jsonDiffDict(value.(map[string]interface{}), json2[key].(map[string]interface{}), depth+1, diff)
				}
			case []interface{}:
				diff.Result = diff.Result + "\n" + longBlank + quotedKey + ": "
				compareDiffDictOfInterface(diff, quotedKey, blank, key, json2, value, depth)
			default:
				compareDiffDictOfDefault(diff, quotedKey, blank, key, longBlank, json2, value)
			}
		} else {
			diff.HasDiff = true
			diff.Result = diff.Result + "\n-" + blank + quotedKey + ": " + marshal(value)
		}
		diff.Result = diff.Result + ","
	}
	//for key, value := range json2 {
	//	if _, ok := json1[key]; !ok {
	//		diff.HasDiff = true
	//		diff.Result = diff.Result + "\n+" + blank + "\"" + key + "\"" + ": " + marshal(value) + ","
	//	}
	//}
	compareJson2KeyWithOutJson1(diff, blank, json1, json2)
	diff.Result = diff.Result + "\n" + blank + "}"
}

func compareJson2KeyWithOutJson1(diff *JsonDiff, blank string, json1, json2 map[string]interface{}) {
	for key, value := range json2 {
		if _, ok := json1[key]; !ok {
			diff.HasDiff = true
			diff.Result = diff.Result + "\n+" + blank + "\"" + key + "\"" + ": " + marshal(value) + ","
		}
	}
}

func compareDiffDictOfDefault(diff *JsonDiff, quotedKey, blank, key, longBlank string, json2 map[string]interface{}, value interface{}) {
	if !reflect.DeepEqual(value, json2[key]) {
		assemblyDiffDictResult(diff, quotedKey, blank, key, json2, value)
	} else {
		diff.Result = diff.Result + "\n" + longBlank + quotedKey + ": " + marshal(value)
	}
}

func compareDiffDictOfInterface(diff *JsonDiff, quotedKey, blank, key string, json2 map[string]interface{}, value interface{}, depth int) {
	if _, ok2 := json2[key].([]interface{}); !ok2 {
		assemblyDiffDictResult(diff, quotedKey, blank, key, json2, value)
	} else {
		jsonDiffList(value.([]interface{}), json2[key].([]interface{}), depth+1, diff)
	}
}

func assemblyDiffDictResult(diff *JsonDiff, quotedKey, blank, key string, json2 map[string]interface{}, value interface{}) {

	diff.HasDiff = true
	diff.Result = diff.Result + "\n-" + blank + quotedKey + ": " + marshal(value) + ","
	diff.Result = diff.Result + "\n+" + blank + quotedKey + ": " + marshal(json2[key])
	diff.ResultMap = append(diff.ResultMap, ChangeField{
		Field:    quotedKey,
		OldVaule: marshal(value),
		NewValue: marshal(json2[key]),
	})
}

func jsonDiffList(json1, json2 []interface{}, depth int, diff *JsonDiff) {
	blank := strings.Repeat(" ", (2 * (depth - 1)))
	longBlank := strings.Repeat(" ", (2 * (depth)))
	diff.Result = diff.Result + "\n" + blank + "["
	size := len(json1)
	if size > len(json2) {
		size = len(json2)
	}
	for i := 0; i < size; i++ {
		switch json1[i].(type) {
		case map[string]interface{}:
			compareDiffListOfMap(diff, blank, i, json1, json2, depth)
		case []interface{}:
			compareDiffListOfInterface(diff, blank, i, json1, json2, depth)
		default:
			compareDiffListOfDefault(diff, blank, longBlank, i, json1, json2)
		}
		diff.Result = diff.Result + ","
	}
	for i := size; i < len(json1); i++ {
		diff.HasDiff = true
		diff.Result = diff.Result + "\n-" + blank + marshal(json1[i])
		diff.Result = diff.Result + ","
	}
	for i := size; i < len(json2); i++ {
		diff.HasDiff = true
		diff.Result = diff.Result + "\n+" + blank + marshal(json2[i])
		diff.Result = diff.Result + ","
	}
	diff.Result = diff.Result + "\n" + blank + "]"
}

func assemblyDiffListResult(diff *JsonDiff, blank string, i int, json1, json2 []interface{}) {
	diff.HasDiff = true
	diff.Result = diff.Result + "\n-" + blank + marshal(json1[i]) + ","
	diff.Result = diff.Result + "\n+" + blank + marshal(json2[i])
}

func compareDiffListOfMap(diff *JsonDiff, blank string, i int, json1, json2 []interface{}, depth int) {
	if _, ok := json2[i].(map[string]interface{}); ok {
		jsonDiffDict(json1[i].(map[string]interface{}), json2[i].(map[string]interface{}), depth+1, diff)
	} else {
		assemblyDiffListResult(diff, blank, i, json1, json2)
	}
}

func compareDiffListOfInterface(diff *JsonDiff, blank string, i int, json1, json2 []interface{}, depth int) {
	if _, ok2 := json2[i].([]interface{}); !ok2 {
		assemblyDiffListResult(diff, blank, i, json1, json2)
	} else {
		jsonDiffList(json1[i].([]interface{}), json2[i].([]interface{}), depth+1, diff)
	}
}

func compareDiffListOfDefault(diff *JsonDiff, blank, longBlank string, i int, json1, json2 []interface{}) {
	if !reflect.DeepEqual(json1[i], json2[i]) {
		assemblyDiffListResult(diff, blank, i, json1, json2)
	} else {
		diff.Result = diff.Result + "\n" + longBlank + marshal(json1[i])
	}
}

func CopyUpd(data mysql.Upd) mysql.Upd {
	result := mysql.Upd{}
	for s, i := range data {
		result[s] = i
	}
	return result
}

func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func CaseCamelCopy(data map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{}, 0)
	for s, i := range data {
		key := Lcfirst(Case2Camel(s))
		res[key] = i
	}

	return res
}
