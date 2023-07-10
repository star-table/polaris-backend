package str

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
)

const (
	// 定义换行符
	EOF = '\x03' // 结束字符
	NL  = '\n'
	CR  = '\r'
	EOL = '\n'
)

// CutCommentObj 按照要求解析评论内容用于计算评论长度。at 符号不算作长度。只计算评论文本。
// ref: https://github.com/suhanyujie/hello_parser_of_js
type CutCommentObj struct {
	Src     string
	SrcRune []rune
	Offset  int
	Line    int
	Col     int
	Total   int // src 总长度
	IsPeek  bool
	Ch      rune
}

func NewCommentCutter(comment string) *CutCommentObj {
	chars := []rune(comment)
	return &CutCommentObj{
		Src:     comment,
		SrcRune: chars,
		Offset:  -1,
		Total:   len(chars),
	}
}

func (obj *CutCommentObj) CutComment(limit int) string {
	cutComment := strings.Builder{}
	textCount := 0
	needOverflowText := false
	isTrue := false
	for {
		partStr := ""
		// obj.SkipWhitespace()
		strTile := obj.Peek(1)
		if strTile == string(EOF) {
			break
		} else if strTile == "<" { // at 某人字符串
			// isTrue 为 true 表示是特殊的字符串，需要忽略计数
			partStr, isTrue = obj.ReadOneAtExp()
			if !isTrue {
				textCount += len(partStr)
			}
		} else if strTile == "[" { // 附件
			// isTrue 为 true 表示是特殊的字符串，需要忽略计数
			partStr, isTrue = obj.ReadAttachmentText()
			if !isTrue {
				textCount += len(partStr)
			}
		} else {
			partStr = obj.ReadCommentText()
			textCount += len(partStr)
		}
		if textCount > limit {
			needOverflowText = true
			break
		}
		cutComment.WriteString(partStr)
	}
	if needOverflowText {
		cutComment.WriteString(consts.CardIssueChangeDescTextOverflow)
	}

	return cutComment.String()
}

// Read 读取并消耗 n 个字符
func (obj *CutCommentObj) Read(n int) string {
	var c rune
	ret := make([]rune, 0)
	offset := obj.Offset
	for i := 0; i < n; i++ {
		nextIdx := offset + 1
		if nextIdx >= obj.Total {
			c = EOF
			ret = append(ret, c)
			break
		}
		c = obj.SrcRune[nextIdx]
		offset = nextIdx

		if c == rune(CR) || c == rune(NL) {
			if c == CR && obj.SrcRune[nextIdx+1] == NL {
				offset++
			}
			if !obj.IsPeek {
				obj.Line++
				obj.Col = 0
			}
			c = EOL
		} else if !obj.IsPeek {
			obj.Col++
		}
		ret = append(ret, c)
	}
	if !obj.IsPeek {
		obj.Ch = c
		obj.Offset = offset
	}

	return string(ret)
}

// Peek 查看接下来的 n 个字符，但不消耗
func (obj *CutCommentObj) Peek(n int) string {
	obj.IsPeek = true
	str := obj.Read(n)
	obj.IsPeek = false
	return str
}

func (obj *CutCommentObj) SkipWhitespace() {
	for {
		strTile := obj.Peek(1)
		if strTile == " " || strTile == "\t" || strTile == string(EOL) {
			obj.Read(1)
			continue
		}
		break
	}
}

func (obj *CutCommentObj) ReadCommentText() string {
	consumeStr := ""
	strTile := obj.Peek(1)
	if exist, _ := slice.Contain([]string{"<", string(EOL), string(EOF)}, strTile); exist {
		// do nothing
	} else {
		strTile = obj.Read(1)
		consumeStr = strTile
	}

	return consumeStr
}

// ReadOneAtExp 读取一个 at 标签字符串
// <at id=ou_3ab7fe596cf91692218f744558ae157f></at>
func (obj *CutCommentObj) ReadOneAtExp() (consumeStr string, isAtText bool) {
	strTile := obj.Peek(48)
	// 必须匹配 `id=ou_3ab7fe596cf91692218f744558ae157f></at>`
	isMatch, _ := regexp.MatchString(`<at id=ou_[\w]+></at>`, strTile)
	if isMatch {
		isAtText = true
		consumeStr = obj.Read(48)
	} else {
		consumeStr = obj.Read(1)
	}

	return
}

// ReadAttachmentText 读取并消耗 `[附件]` 文本
func (obj *CutCommentObj) ReadAttachmentText() (consumeStr string, isAttach bool) {
	strTile := obj.Peek(4)
	// 必须匹配 `id=ou_3ab7fe596cf91692218f744558ae157f></at>`
	isMatch := strTile == "[附件]"
	if isMatch {
		isAttach = true
		consumeStr = obj.Read(4)
	} else {
		consumeStr = obj.Read(1)
	}
	return
}

// TruncateText 截断文本
// count 保留的字符个数
func TruncateText(text string, count int) string {
	if count <= 0 {
		return ""
	}
	if RuneLen(text) > count {
		text = Substr(text, 0, count) + consts.CardIssueChangeDescTextOverflow
	}

	return text
}

// GetTruncatedComment 评论字数截断，不包含@成员的字数
// 卫忠实现的 暂时没用。
func GetTruncatedComment(comment string, limitWords int, userIds []int64) string {
	actualLimits := limitWords + len(userIds)*48
	noAtComment := GetCardCommentRealWords(comment)
	noAtCommentLen := RuneLen(noAtComment)
	a := RuneLen(comment)
	fmt.Println(a)
	if RuneLen(comment) > actualLimits {
		//i := RuneLen(comment) - actualLimits
		lastStr := comment[noAtCommentLen-1+(len(userIds)-1)*48-1:]
		if !strings.Contains(lastStr, "<") {
			// 说明被截断不要的字符中没有@成员
			return comment[:actualLimits] + consts.CardIssueChangeDescTextOverflow
		} else {
			// 被截断后面有@成员的字符
			// 有两种情况 1、被截断没有需要需要显示的字符  2、需要显示的字符被截断在后面了
			return comment[:noAtCommentLen-1+(len(userIds)-1)*48-1] + consts.CardIssueChangeDescTextOverflow
		}
	}

	return comment
}

// TruncateColumnName 飞书卡片中，列名部分的截断处理
func TruncateColumnName(column lc_table.LcCommonField) string {
	displayName := column.Label
	if column.AliasTitle != "" {
		displayName = column.AliasTitle
	}
	if RuneLen(displayName) > consts.ColumnAliasLimit {
		// header 部分，6个点的省略号太长了，使用 3 个点的省略号
		displayName = Substr(displayName, 0, consts.ColumnAliasLimit) + consts.CardIssueColumnNameOverflow
	}

	// displayName += consts.GetTabCharacter(displayName)

	return displayName
}

func TruncateName(name string, limit int) string {
	if RuneLen(name) > limit {
		// header 部分，6个点的省略号太长了，使用 3 个点的省略号
		name = Substr(name, 0, limit) + consts.CardIssueColumnNameOverflow
	}

	return name
}

// 获取过滤掉@成员的 字符
func GetCardCommentRealWords(comment string) string {
	reg := regexp.MustCompile(consts.FsCardCommentPattern)
	return reg.ReplaceAllString(comment, "")
}
