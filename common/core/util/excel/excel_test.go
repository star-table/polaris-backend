package excel

import (
	"fmt"
	"github.com/star-table/common/core/util/slice"
	"testing"
)

func TestGenerateCSVFromXLSXFile2(t *testing.T) {
	a := []int64{1, 2}
	fmt.Println(slice.Contain(a, 1))
}

type OneField struct {
	HeaderName string
	HandleFunc func() (interface{}, error)
}

// 从 excel 中解析出数据，组装成 map 方式。导入。
func TestGenerateCSVFromXLSXFile3(t *testing.T) {
	//sortMap := NewSortMap()
	//sortMap.Insert("issueType", "任务类型")
	//sortMap.Insert("issueTitle", "任务标题")
	//sortMap.Insert("owner", "负责人")
	//sortMap.Insert("isParent", "父子类型")
	//sortMap.Insert("priority", "优先级")
	//sortMap.Insert("demandType", "需求类型")
	//sortMap.Insert("demandSource", "需求来源")
	////sortMap.Insert("bugType", "缺陷类型")
	////sortMap.Insert("bugProperty", "严重程度")
	//sortMap.Insert("issueStatus", "任务状态")
	//sortMap.Insert("desc", "任务描述")
	//sortMap.Insert("startTime", "任务开始时间")
	//sortMap.Insert("endTime", "任务结束时间")
	//sortMap.Insert("tag", "标签")
	//sortMap.Insert("follower", "关注人")
	//sortMap.Insert("createTime", "任务创建时间")
	//sortMap.Insert("createUser", "任务创建者")
	//sortMap.Insert("codeId", "编号ID")
	//sortMap.Insert("completeTime", "实际完成时间")
	//for _, key := range sortMap.KeyVec {
	//	if val, ok := sortMap.DataMap[key]; ok {
	//		fmt.Println(val)
	//	}
	//}
	// input: headerMap, dataMapList
	// 通过表头，匹配 sheet 中的列
	// 将一行数据转换为 map，以进行下一步处理
	// 为了后续可扩展性，采用一个字段，使用一个闭包处理，此外会有默认的闭包处理

}

func issueTypeHandle() (interface{}, error) {
	return 1, nil
}

func defaultHandle() (interface{}, error) {
	return 1, nil
}
