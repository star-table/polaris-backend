package int64

import (
	"strconv"
	"strings"
)

// 查找在 arr1 中，但不在 arr2 中的元素
func ArrayDiff(arr1, arr2 []int64) (diffArr []int64) {
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

// 去重，并且不会打乱顺序
func ArrayUnique(arr1 []int64) []int64 {
	resArr := []int64{}
	uniqueMap := map[int64]bool{}
	for _, one := range arr1 {
		if _, ok := uniqueMap[one]; ok {
			continue
		} else {
			uniqueMap[one] = true
			resArr = append(resArr, one)
		}
	}
	return resArr
}

func InArray(val int64, list []int64) bool {
	res := false
	for _, item := range list {
		if item == val {
			res = true
			break
		}
	}
	return res
}

// Int64Intersect 两个 int64 数组切片的交集
func Int64Intersect(list1, list2 []int64) []int64 {
	uniqueMap := map[int64]bool{}
	for _, ele := range list1 {
		uniqueMap[ele] = true
	}
	result := make([]int64, 0)
	for _, ele := range list2 {
		if _, ok := uniqueMap[ele]; ok {
			result = append(result, ele)
		}
	}
	return result
}

// StringArrIntersect 两个 string 数组切片的交集
func StringArrIntersect(list1, list2 []string) []string {
	uniqueMap := map[string]bool{}
	for _, ele := range list1 {
		uniqueMap[ele] = true
	}
	result := make([]string, 0)
	for _, ele := range list2 {
		if _, ok := uniqueMap[ele]; ok {
			result = append(result, ele)
		}
	}
	return result
}

// 移除 haystack 中的一些元素
func Int64RemoveSomeVal(haystack, removeVals []int64) []int64 {
	if len(removeVals) < 1 {
		return haystack
	}
	result := make([]int64, 0)
	for _, ele := range haystack {
		if InArray(ele, removeVals) {
			continue
		}
		result = append(result, ele)
	}
	return result
}

// Int64Explode 将切片内容拼接成字符串
func Int64Explode(list []int64, glue string) string {
	strList := make([]string, 0)
	for _, item := range list {
		strList = append(strList, strconv.FormatInt(item, 10))
	}
	return strings.Join(strList, glue)
}

// Int64ArrToInterfaceArr int64 类型的切片转为 interface{} 切片
func Int64ArrToInterfaceArr(int64Arr []int64) []interface{} {
	newArr := make([]interface{}, len(int64Arr))
	for _, item := range int64Arr {
		newArr = append(newArr, item)
	}

	return newArr
}

// CompareSliceInt64 判断 两个int64类型的切片是否有不同
func CompareSliceInt64(list1, list2 []int64) bool {
	if len(list1) != len(list2) {
		return false
	}
	if (list1 == nil) != (list2 == nil) {
		return false
	}
	for k, v := range list1 {
		if v != list2[k] {
			return false
		}
	}
	return true
}

// CompareSliceAddDelInt64 判断 两个int64类型的切片是否有不同，获得新增和删除的集合，顺序无所谓
func CompareSliceAddDelInt64(new, old []int64) (same bool, add []int64, del []int64) {
	for _, v1 := range new {
		found := false
		for _, v2 := range old {
			if v1 == v2 {
				found = true
				break
			}
		}
		if !found {
			add = append(add, v1)
		}
	}
	for _, v1 := range old {
		found := false
		for _, v2 := range new {
			if v1 == v2 {
				found = true
				break
			}
		}
		if !found {
			del = append(del, v1)
		}
	}
	same = len(add) == 0 && len(del) == 0
	return
}

// CompareSliceAddDelString 判断 两个 string 类型的切片是否有不同，获得新增和删除的集合，顺序无所谓
func CompareSliceAddDelString(new, old []string) (same bool, add []string, del []string) {
	for _, v1 := range new {
		found := false
		for _, v2 := range old {
			if v1 == v2 {
				found = true
				break
			}
		}
		if !found {
			add = append(add, v1)
		}
	}
	for _, v1 := range old {
		found := false
		for _, v2 := range new {
			if v1 == v2 {
				found = true
				break
			}
		}
		if !found {
			del = append(del, v1)
		}
	}
	same = len(add) == 0 && len(del) == 0

	return
}

// GetSearchedIndexArr 从 int64 数组中，查找一个数，返回所在索引。如果未找到，返回 -1
func GetSearchedIndexArr(list []int64, needle int64) int {
	for index, item := range list {
		if needle == item {
			return index
		}
	}
	return -1
}
