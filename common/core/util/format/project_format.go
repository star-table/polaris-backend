package format

import (
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/polaris-backend/common/core/util"
	"regexp"
)

//项目名
func VerifyProjectNameFormat(input string) bool {
	reg := regexp.MustCompile(ProjectNamePattern)
	return reg.MatchString(input)
}

//项目前缀编号
func VerifyProjectPreviousCodeFormat(input string) bool {
	reg := regexp.MustCompile(ProjectPreviousCodePattern)
	return reg.MatchString(input)
}

//项目简介
func VerifyProjectRemarkFormat(input string) bool {
	reg := regexp.MustCompile(ProjectRemarkPattern)
	return reg.MatchString(input)
}

//任务名
func VerifyIssueNameFormat(input string) bool {
	reg := regexp.MustCompile(IssueNamePattern)
	return reg.MatchString(input)
}

//任务简介
func VerifyIssueRemarkFormat(input string) bool {
	//reg := regexp.MustCompile(IssueRemarkPattern)
	//return reg.MatchString(input)
	str := util.RenderCommentContentToMarkDown(input, true)
	inputLen := strs.Len(str)
	if inputLen > 10000 {
		return false
	} else {
		return true
	}
}

//任务评论
func VerifyIssueCommenFormat(input string) bool {
	//先判断是否为空
	if len(input) == 0 {
		return false
	}
	//去除@后的长度可以为空
	//reg := regexp.MustCompile(IssueCommenPattern)
	//str := util.RenderCommentContentToMarkDown(input, true)
	//inputLen := strs.Len(str)

	// 现在文本框+@成员 字符 限制2000个
	reg := regexp.MustCompile(IssueCommentMemberPattern)
	allString := reg.ReplaceAllString(input, "$1")
	charLen := []rune(allString)

	if len(charLen) > 2000 {
		return false
	} else {
		return true
	}
}

//任务栏名字
func VerifyProjectObjectTypeNameFormat(input string) bool {
	reg := regexp.MustCompile(ProjectObjectTypeNamePattern)
	return reg.MatchString(input)
}

// 表名
func VerifyTableNameFormat(name string) bool {
	reg := regexp.MustCompile(ProjectTableNamePattern)
	return reg.MatchString(name)
}

//项目公告
func VerifyProjectNoticeFormat(input string) bool {
	inputLen := strs.Len(input)
	if inputLen > 2000 {
		return false
	} else {
		return true
	}
}

//字段判断
func VerifySqlFieldFormat(input string) bool {
	reg := regexp.MustCompile(SqlFieldPattern)
	return reg.MatchString(input)
}
