package sortedmap

import "github.com/pkg/errors"

// SortMap 从 excel 获取数据时，都是以数组的方式，这样有一个不足之处。当我们要对 excel 解析时，格式发生变化，对应的解析过程也需要改变，并且是以索引的方式取值，并解析时，降低了可维护性
// 例如开始时间是「开始时间」字段为第 9 列，首次开发功能时，我们解析取值是 `row[8]`，但下次迭代时，我们可能需要改变之前的字段列，导致解析过程需要变化。
// 因此优化为使用 map 的处理方式。增强关联性，便于处理。SortMap 是类似于一种有序 map 的数据结构。
// 一行数据存储在一个 SortMap 实例中，达到以 map 的方式处理 excel 数据。
type SortMap struct {
	DataMap           map[string]interface{}
	KeyVec            []string
	InsertValueHandle func(value interface{}) interface{}
}

// SortMap 构造函数。
func NewSortMap() *SortMap {
	return &SortMap{
		DataMap:           map[string]interface{}{},
		KeyVec:            []string{},
		InsertValueHandle: nil,
	}
}

// 想 SortMap 中插入一个键值对。
func (m *SortMap) Insert(key string, value interface{}) {
	if m.InsertValueHandle != nil {
		value = m.InsertValueHandle(value)
	}
	if _, ok := m.DataMap[key]; ok {
		m.DataMap[key] = value
	} else {
		m.DataMap[key] = value
		m.KeyVec = append(m.KeyVec, key)
	}
}

// 对插入的值做一些特殊的处理，例如：根据名称做中英文转换
func (m *SortMap) SetInsertValueHandler(f func(value interface{}) interface{}) {
	m.InsertValueHandle = f
}

// 通过键，获取值
func (m *SortMap) GetByKey(key string, defaultVal string) string {
	resValue := defaultVal
	if val, ok := m.DataMap[key]; ok {
		resValue = val.(string)
		if len(resValue) < 1 {
			resValue = defaultVal
		}
	}
	return resValue
}

// 给定一个已知的 SortMap，通过 index，获取 key
func (m *SortMap) GetKeyByIndex(index int) (string, error) {
	if index >= len(m.KeyVec) {
		return "", errors.New("索引超出范围")
	}
	return m.KeyVec[index], nil
}

// 返回一个 map，键为数据的键，值为对应列的索引，从 0 开始。
func (m *SortMap) GetMapKeyByKey() map[string]int {
	m1 := make(map[string]int, len(m.KeyVec))
	for index, keyVal := range m.KeyVec {
		m1[keyVal] = index
	}
	return m1
}

func (m *SortMap) Len() int {
	return len(m.KeyVec)
}
