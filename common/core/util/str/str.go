package str

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/temp"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/axgle/mahonia"
	"github.com/shopspring/decimal"
)

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

//start：正数 - 在字符串的指定位置开始,超出字符串长度强制把start变为字符串长度
//       负数 - 在从字符串结尾的指定位置开始
//       0 - 在字符串中的第一个字符处开始
//length:正数 - 从 start 参数所在的位置返回
//       负数 - 从字符串末端返回
func Substr(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	runeStr := []rune(str)
	len_str := len(runeStr)

	if start < 0 {
		start = len_str + start
	}
	if start > len_str {
		start = len_str
	}
	end := start + length
	if end > len_str {
		end = len_str
	}
	if length < 0 {
		end = len_str + length
	}
	if start > end {
		start, end = end, start
	}
	return string(runeStr[start:end])
}

func UrlParse(url string) (string, string) {
	var host, path string
	split := strings.Split(url, "/")
	if split[0] != "http:" && split[0] != "https:" {
		return host, url
	}
	host = strings.Join(split[:3], "/")
	path = "/" + strings.Join(split[3:], "/")
	return host, path
}

func ParseOssKey(url string) string {
	_, path := UrlParse(url)
	if strings.HasPrefix(path, "/") {
		return path[1:strs.Len(path)]
	}
	return path
}

func CountStrByGBK(str string) int {
	new := mahonia.NewEncoder("gbk").ConvertString(str)
	return len(new)
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//将p标签换成换行符
	re, _ = regexp.Compile("\\<p[^>]*>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("</p>")
	src = re.ReplaceAllString(src, "\n")
	// 将br标签换成换行符
	re, _ = regexp.Compile("<br/>")
	src = re.ReplaceAllString(src, "\n")
	// 替换img为（url）
	re, _ = regexp.Compile("\\<img[^>]*src=(\\'|\")(.*?)[^>]*>")
	for _, s := range re.FindAllString(src, -1) {
		re, _ = regexp.Compile("src=(\\'|\")([^\\'\"]*)(\\'|\")")
		url := re.FindAllString(s, -1)
		if len(url) == 1 {
			src = strings.ReplaceAll(src, s, "[附件]("+url[0][5:len(url[0])-1]+")")
		}
	}
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "")
	//去除连续的换行符
	//re, _ = regexp.Compile("\\s{2,}")
	//src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

// 去除评论@的标签
func TrimComment(str string) string {
	re := regexp.MustCompile("<at id=[\\s\\S]*></at>")
	return re.ReplaceAllString(str, "")
}

// 元素为 string 的 slice 版本的 InArray 函数
func InArray(needle string, strSlice []string) bool {
	hasFound := false
	if len(strSlice) < 1 {
		return hasFound
	}
	for _, val := range strSlice {
		if needle == val {
			hasFound = true
			break
		}
	}
	return hasFound
}

// 用于拼接飞书得 at 标签
func RenderAtSomeoneStr(openId, userName string) string {
	tmp := `<at id="{{.OpenId}}">@{{.UserName}}</at>`
	params := map[string]interface{}{
		"OpenId":   openId,
		"UserName": userName,
	}
	result, _ := temp.Render(tmp, params)
	return result
}

func ToString(dest interface{}) string {
	var key string
	if dest == nil {
		return key
	}
	switch dest.(type) {
	case float64:
		key = decimal.NewFromFloat(dest.(float64)).String()
	case *float64:
		key = decimal.NewFromFloat(*dest.(*float64)).String()
	case float32:
		key = decimal.NewFromFloat32(dest.(float32)).String()
	case *float32:
		key = decimal.NewFromFloat32(*dest.(*float32)).String()
	case int:
		key = strconv.Itoa(dest.(int))
	case *int:
		key = strconv.Itoa(*dest.(*int))
	case uint:
		key = strconv.Itoa(int(dest.(uint)))
	case *uint:
		key = strconv.Itoa(int(*dest.(*uint)))
	case int8:
		key = strconv.Itoa(int(dest.(int8)))
	case *int8:
		key = strconv.Itoa(int(*dest.(*int8)))
	case uint8:
		key = strconv.Itoa(int(dest.(uint8)))
	case *uint8:
		key = strconv.Itoa(int(*dest.(*uint8)))
	case int16:
		key = strconv.Itoa(int(dest.(int16)))
	case *int16:
		key = strconv.Itoa(int(*dest.(*int16)))
	case uint16:
		key = strconv.Itoa(int(dest.(uint16)))
	case *uint16:
		key = strconv.Itoa(int(*dest.(*uint16)))
	case int32:
		key = strconv.Itoa(int(dest.(int32)))
	case *int32:
		key = strconv.Itoa(int(*dest.(*int32)))
	case uint32:
		key = strconv.Itoa(int(dest.(uint32)))
	case *uint32:
		key = strconv.Itoa(int(*dest.(*uint32)))
	case int64:
		key = strconv.FormatInt(dest.(int64), 10)
	case *int64:
		key = strconv.FormatInt(*dest.(*int64), 10)
	case uint64:
		key = strconv.FormatUint(dest.(uint64), 10)
	case *uint64:
		key = strconv.FormatUint(*dest.(*uint64), 10)
	case string:
		key = dest.(string)
	case *string:
		key = *dest.(*string)
	case []byte:
		key = string(dest.([]byte))
	case *[]byte:
		key = string(*dest.(*[]byte))
	case bool:
		if dest.(bool) {
			key = "true"
		} else {
			key = "false"
		}
	case *bool:
		if *dest.(*bool) {
			key = "true"
		} else {
			key = "false"
		}
	default:
		newValue, _ := json.Marshal(dest)
		key = string(newValue)
	}
	return key
}

func ToInt64(dest interface{}) (int64, error) {
	var key int64
	switch dest.(type) {
	case float64:
		key = int64(dest.(float64))
	case float32:
		key = int64(dest.(float32))
	case int:
		key = int64(dest.(int))
	case uint:
		key = int64(dest.(uint))
	case int8:
		key = int64(dest.(int8))
	case uint8:
		key = int64(dest.(uint8))
	case int16:
		key = int64(dest.(int16))
	case uint16:
		key = int64(dest.(uint16))
	case int32:
		key = int64(dest.(int32))
	case uint32:
		key = int64(dest.(uint32))
	case int64:
		key = int64(dest.(int64))
	case uint64:
		key = int64(dest.(uint64))
	case string:
		return strconv.ParseInt(dest.(string), 10, 64)
	default:
		return 0, errs.TypeConvertError
	}
	return key, nil
}

// ArrayDiff 对比，获取在 arr1 中存在，但在 arr2 中不存在的元素集合
func ArrayDiff(arr1, arr2 []string) (diffArr []string) {
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

func Last(str string, n int) string {
	if strs.Len(str) <= n {
		return str
	}
	return str[strs.Len(str)-n:]
}

// Int64Implode 将切片内容拼接成字符串
func Int64Implode(list []int64, glue string) string {
	strList := make([]string, len(list))
	for i, item := range list {
		strList[i] = strconv.FormatInt(item, 10)
	}
	return strings.Join(strList, glue)
}

// ToPtr 字符串转指针
func ToPtr(s string) *string {
	return &s
}

func RuneLen(s string) int {
	s1 := []rune(s)
	return len(s1)
}

// ReplaceWhiteSpaceCharToSpace 将空白字符替换为一个空格
func ReplaceWhiteSpaceCharToSpace(str string) string {
	str = strings.ReplaceAll(str, "  ", " ")
	str = strings.ReplaceAll(str, "\t", " ")
	str = strings.ReplaceAll(str, "\n", " ")

	return str
}

// CheckStrInArray 判断字符串是否字符串切片中
func CheckStrInArray(strArray []string, target string) bool {
	if strArray == nil {
		return true
	}
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}

// DeleteSliceElement 删除指定字符串元素
func DeleteSliceElement(array []string, elem string) []string {
	tmp := make([]string, 0, len(array))
	for _, v := range array {
		if v != elem {
			tmp = append(tmp, v)
		}
	}
	return tmp
}
