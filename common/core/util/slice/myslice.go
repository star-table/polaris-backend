package slice

import (
	"strconv"

	"github.com/spf13/cast"
)

// 切片元素的删除
func SliceRemove(list *[]interface{}, index int) {
	//if index >= len(list) || (index + 1) > len(list) {
	//	return list
	//}
	*list = append((*list)[:index], (*list)[index+1:])
}

// Int64Slice2IfSlice int64 切片转 interface 切片
func Int64Slice2IfSlice(int64Arr []int64) []interface{} {
	ifValArr := make([]interface{}, len(int64Arr))
	for i, val := range int64Arr {
		ifValArr[i] = val
	}

	return ifValArr
}

// Int64ToStringSlice int64 转 string 切片
func Int64ToStringSlice(arr []int64) []string {
	resArr := make([]string, 0, len(arr))
	for _, ele := range arr {
		resArr = append(resArr, strconv.FormatInt(ele, 10))
	}

	return resArr
}

func StringToInt64Slice(arr []string) []int64 {
	resArr := make([]int64, 0, len(arr))
	for _, ele := range arr {
		resArr = append(resArr, cast.ToInt64(ele))
		//tmpVal, _ := strconv.ParseInt(ele, 10, 64)
		//if tmpVal > 0 {
		//	resArr = append(resArr, tmpVal)
		//}
	}

	return resArr
}

// Int64ArrayDiff 对比，获取在 arr1 中存在，但在 arr2 中不存在的元素集合
func Int64ArrayDiff(arr1, arr2 []int64) (diffArr []int64) {
	if len(arr2) < 1 || len(arr1) < 1 {
		diffArr = arr1
		return
	}
	for i := 0; i < len(arr1); i++ {
		item := arr1[i]
		isIn := false
		for j := 0; j < len(arr2); j++ {
			if item == arr2[j] {
				isIn = true
				break
			}
		}
		if !isIn {
			diffArr = append(diffArr, item)
		}
	}

	return diffArr
}
