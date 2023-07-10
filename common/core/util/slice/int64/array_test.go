package int64

import "testing"

// 测试 ArrayDiff 函数
func TestArrayDiff(t *testing.T) {
	arr1 := []int64{1, 2, 3}
	arr2 := []int64{1, 2, 3}
	var diff []int64
	diff = ArrayDiff(arr1, arr2)
	if len(diff) > 0 {
		t.Error("error result 001")
		return
	}
	arr1 = []int64{1, 2, 3}
	arr2 = []int64{1, 2}
	diff = ArrayDiff(arr1, arr2)
	if diff[0] != 3 {
		t.Error("error result 002")
		return
	}
	arr1 = []int64{}
	arr2 = []int64{1, 2}
	diff = ArrayDiff(arr1, arr2)
	if len(diff) != 0 {
		t.Error("error result 003")
		return
	}
	arr1 = []int64{1,2,3}
	arr2 = []int64{}
	diff = ArrayDiff(arr1, arr2)
	if len(diff) != 3 {
		t.Error("error result 004")
		return
	}
	t.Log(diff)
}

func TestInt64RemoveSomeVal(t *testing.T) {
	list1 := []int64{1,2,3,4,5,8,10}
	removeVals := []int64{2, 3, 8 , 10}
	result := Int64RemoveSomeVal(list1, removeVals)
	t.Log(result)
}

func TestInt64Intersect(t *testing.T) {
	list1 := []int64{1,2,3,4}
	list2 := []int64{2,3,4}
	res := Int64Intersect(list1, list2)
	t.Log(res)
}
