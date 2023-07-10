package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/md5"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/temp"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"gopkg.in/fatih/set.v0"
	"upper.io/db.v3/lib/sqlbuilder"
)

const (
	AtPrefix = "@#["
	AtSuffix = "]&$"
	AtSplit  = ":"
)

var log = logger.GetDefaultLogger()

func ConvertObject(src interface{}, source interface{}) errs.SystemErrorInfo {
	err := copyer.Copy(src, source)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return nil
}

func FieldInUpdate(updateFields []string, field string) bool {
	if updateFields == nil {
		return true
	}
	bol, err := slice.Contain(updateFields, field)
	if err != nil {
		return false
	}
	return bol
}

//判断是否为bool类型（为1或2）
func IsBool(t int) bool {
	return t == 1 || t == 2
}

// InterfaceToInt64 interface 转 int64
func InterfaceToInt64(v interface{}) (int64, bool) {
	res, err := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
	if err != nil {
		return 0, false
	}
	return res, true
}

/**
 方法返回参数结束
移除的ids|新增的ids|err
*/
func GetDifMemberIds(beforeUserIds []int64, afterUserIds []int64) ([]int64, []int64) {
	beforeChangeMembersSet := set.New(set.ThreadSafe)
	for _, member := range beforeUserIds {
		beforeChangeMembersSet.Add(member)
	}
	afterChangeMembersSet := set.New(set.ThreadSafe)
	for _, member := range afterUserIds {
		afterChangeMembersSet.Add(member)
	}

	deletedMemberSet := set.Difference(beforeChangeMembersSet, afterChangeMembersSet)
	addedMemberSet := set.Difference(afterChangeMembersSet, beforeChangeMembersSet)

	deletedMemberIds := convertSetToArray(deletedMemberSet)
	addedMemberIds := convertSetToArray(addedMemberSet)

	return deletedMemberIds, addedMemberIds
}

func convertSetToArray(l set.Interface) []int64 {
	arr := make([]int64, l.Size())
	for i, id := range l.List() {
		arr[i] = id.(int64)
	}
	return arr
}

/**
 方法返回参数结束
移除的ids|新增的ids|err
*/
func GetDifMemberIdsByString(beforeUserIds []string, afterUserIds []string) ([]string, []string) {
	beforeChangeMembersSet := set.New(set.ThreadSafe)
	for _, member := range beforeUserIds {
		beforeChangeMembersSet.Add(member)
	}
	afterChangeMembersSet := set.New(set.ThreadSafe)
	for _, member := range afterUserIds {
		afterChangeMembersSet.Add(member)
	}

	deletedMemberSet := set.Difference(beforeChangeMembersSet, afterChangeMembersSet)
	addedMemberSet := set.Difference(afterChangeMembersSet, beforeChangeMembersSet)

	deletedMemberIds := convertSetToArrayString(deletedMemberSet)
	addedMemberIds := convertSetToArrayString(addedMemberSet)

	return deletedMemberIds, addedMemberIds
}

func convertSetToArrayString(l set.Interface) []string {
	arr := make([]string, l.Size())
	for i, id := range l.List() {
		arr[i] = id.(string)
	}
	return arr
}

func ReadAndWrite(dirFrom string, context map[string]interface{}, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	originSql, err := FileRead(dirFrom)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.FileReadFail, err)
	}
	resultSql, err := temp.Render(originSql, context)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.TemplateRenderError, err)
	}

	batch := strings.Split(resultSql, ";")
	for _, v := range batch {
		log.Infof("执行sql: %v", v)
		if v == "" {
			continue
		}
		_, err := tx.Exec(v)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}

	return nil
}

func FileRead(dir string) (string, error) {

	path, _ := os.Getwd()
	fmt.Println("当前路径：", path)

	fileObj, err := os.Open(dir)
	defer func() {
		if err := fileObj.Close(); err != nil {
			log.Info(strs.ObjectToString(err))
		}
	}()

	if err != nil {
		return "", err
	}
	contents, err := ioutil.ReadAll(fileObj)
	if err != nil {
		return "", err
	}

	result := strings.Replace(string(contents), "\n", "", 0)
	return result, nil
}

func WriteWithIoutil(name, content string) {
	data := []byte(content)
	if ioutil.WriteFile(name, data, 0644) == nil {
		log.Info("写入文件成功:" + content)
	}
}

func PageOption(page *int, size *int) (pageA uint, sizeA uint) {
	pageA = uint(0)
	sizeA = uint(0)
	if page != nil && size != nil && *page > 0 && *size > 0 {
		pageA = uint(*page)
		sizeA = uint(*size)
	}
	return
}

//评论提及规则: @#[姓名:open_id:user_id]&$
func RenderCommentContentToMarkDown(content string, onlyComment bool) string {
	return RenderCommentContentToMarkDownWithOpenIdMap(content, onlyComment, map[string]string{})
}

//评论 带替换
func RenderCommentContentToMarkDownWithOpenIdMap(content string, onlyComment bool, openIdMap map[string]string) string {
	contentLen := strs.Len(content)
	if contentLen == 0 {
		return content
	}
	var buffer bytes.Buffer

	var subBuffer bytes.Buffer
	status := 1

	chars := []rune(content)

	for i := 0; i < contentLen; i++ {
		c := chars[i]
		if status == 1 {
			if isMatch, nextIndex := renderCommentMatchPrefix(i, chars, contentLen); isMatch {
				i = nextIndex
				status = 2
				continue
			}
		} else if status == 2 {
			if isMatch, nextIndex := renderCommentMatchSuffix(i, chars, contentLen); isMatch {
				str := subBuffer.String()
				strs := strings.Split(str, AtSplit)
				if len(strs) >= 2 {
					if !onlyComment {
						if openId, ok := openIdMap[strs[1]]; ok {
							buffer.WriteString("<at id=" + openId + "></at>")
						} else {
							buffer.WriteString("<at id=" + strs[1] + "></at>")
						}
					} else {
						buffer.WriteString("")
					}
				}
				subBuffer.Reset()
				i = nextIndex
				status = 1
			} else {
				subBuffer.WriteString(string(c))
			}
			continue
		}
		buffer.WriteString(string(c))
	}
	return buffer.String()
}

// 将评论转成原始格式，输入@#[xxx:open_id:user_id]&$ 你好， 输出@xxx 你好
func RenderCommentContentToOrigin(content string) string {
	contentLen := strs.Len(content)
	if contentLen == 0 {
		return content
	}
	var buffer bytes.Buffer
	var subBuffer bytes.Buffer
	status := 1
	chars := []rune(content)

	for i := 0; i < contentLen; i++ {
		c := chars[i]
		if status == 1 {
			if isMatch, nextIndex := renderCommentMatchPrefix(i, chars, contentLen); isMatch {
				i = nextIndex
				status = 2
				continue
			}
		} else if status == 2 {
			if isMatch, nextIndex := renderCommentMatchSuffix(i, chars, contentLen); isMatch {
				str := subBuffer.String()
				strs := strings.Split(str, AtSplit)
				if len(strs) >= 0 {
					buffer.WriteString("@" + strs[0] + " ")
				}
				subBuffer.Reset()
				i = nextIndex
				status = 1
			} else {
				subBuffer.WriteString(string(c))
			}
			continue
		}
		buffer.WriteString(string(c))
	}
	return buffer.String()
}

func GetCommentAtUserIds(content string) []int64 {
	userIds := make([]int64, 0)
	contentLen := strs.Len(content)
	if contentLen == 0 {
		return userIds
	}
	var buffer bytes.Buffer

	var subBuffer bytes.Buffer
	status := 1

	chars := []rune(content)

	for i := 0; i < contentLen; i++ {
		c := chars[i]
		if status == 1 {
			if isMatch, nextIndex := renderCommentMatchPrefix(i, chars, contentLen); isMatch {
				i = nextIndex
				status = 2
				continue
			}
		} else if status == 2 {
			if isMatch, nextIndex := renderCommentMatchSuffix(i, chars, contentLen); isMatch {
				str := subBuffer.String()
				strs := strings.Split(str, AtSplit)
				if len(strs) >= 2 {
					userId, err := strconv.ParseInt(strs[1], 10, 64)
					if err == nil {
						userIds = append(userIds, userId)
					}
				}
				subBuffer.Reset()
				i = nextIndex
				status = 1
			} else {
				subBuffer.WriteString(string(c))
			}
			continue
		}
		buffer.WriteString(string(c))
	}
	if len(userIds) > 0 {
		userIds = slice.SliceUniqueInt64(userIds)
	}
	return userIds
}

//是否匹配前缀，并返回新下标
func renderCommentMatchPrefix(index int, content []rune, len int) (bool, int) {
	prefixLen := strs.Len(AtPrefix)
	limit := index + prefixLen - 1
	if limit < len {
		sub := content[index : limit+1]
		if string(sub) == AtPrefix {
			return true, limit
		}
	}
	return false, -1
}

//是否匹配后缀，并返回新下标
func renderCommentMatchSuffix(index int, content []rune, len int) (bool, int) {
	prefixLen := strs.Len(AtSuffix)
	limit := index + prefixLen - 1
	if limit < len {
		sub := content[index : limit+1]
		if string(sub) == AtSuffix {
			return true, limit
		}
	}
	return false, -1
}

//拼接url
func JointUrl(host, path string) string {
	if strings.HasSuffix(host, "/") && strings.HasPrefix(path, "/") {
		return host + path[1:strs.Len(path)]
	} else if !strings.HasSuffix(host, "/") && !strings.HasPrefix(path, "/") {
		return host + "/" + path
	} else {
		return host + path
	}
}

//获取文件后缀
func ParseFileSuffix(fileName string) string {
	suffix := ""
	if strings.Index(fileName, ".") != -1 {
		suffixSplit := strings.Split(fileName, ".")
		suffix = suffixSplit[len(suffixSplit)-1]
	}
	return suffix
}

//获取文件名
func ParseFileName(path string) string {
	nameSplit := strings.Split(path, "/")
	return nameSplit[len(nameSplit)-1]
}

//植入文件名
func ModifyFileName(path string, str string) string {
	fileName := ParseFileName(path)
	filePath := ""
	suffix := ""
	if strings.Index(fileName, ".") != -1 {
		suffixSplit := strings.Split(path, ".")
		for i := 0; i < len(suffixSplit)-1; i++ {
			filePath += suffixSplit[i] + "."
		}
		suffix = suffixSplit[len(suffixSplit)-1]
	}
	if filePath != "" {
		filePath = filePath[0 : strs.Len(filePath)-1]
	}
	return filePath + str + "." + suffix
}

func GetCompressedPath(path string, typ int) string {
	if typ == consts.OssResource {
		return path + "?x-oss-process=style/thumbnail_001"
	} else if typ == consts.LocalResource {
		return ModifyFileName(path, "_compressed")
	}
	return path
}

func GetOssKeyInfo(key string) bo.OssKeyInfo {
	segments := strings.Split(key, "/")
	ossKeyInfo := bo.OssKeyInfo{}
	for _, segment := range segments {
		if strings.HasPrefix(segment, consts.OssKeySegmentOrg) {
			id, err := strconv.ParseInt(segment[strs.Len(consts.OssKeySegmentOrg):strs.Len(segment)], 10, 64)
			if err != nil {
				continue
			}
			ossKeyInfo.OrgId = id
		} else if strings.HasPrefix(segment, consts.OssKeySegmentProject) {
			id, err := strconv.ParseInt(segment[strs.Len(consts.OssKeySegmentProject):strs.Len(segment)], 10, 64)
			if err != nil {
				continue
			}
			ossKeyInfo.ProjectId = id
		} else if strings.HasPrefix(segment, consts.OssKeySegmentIssue) {
			id, err := strconv.ParseInt(segment[strs.Len(consts.OssKeySegmentIssue):strs.Len(segment)], 10, 64)
			if err != nil {
				continue
			}
			ossKeyInfo.IssueId = id
		}
	}
	return ossKeyInfo
}

func CheckIssueCommentLen(comment string) errs.SystemErrorInfo {
	commentLen := strs.Len(comment)
	if commentLen <= 0 || commentLen > 500 {
		return errs.IssueCommentLenError
	}
	return nil
}

func CheckIssueTitleLen(str string) errs.SystemErrorInfo {
	strLen := strs.Len(str)
	if strLen <= 0 || strLen > 50 {
		return errs.IssueTitleError
	}
	return nil
}

func CheckIssueRemarkLen(str string) errs.SystemErrorInfo {
	strLen := strs.Len(str)
	if strLen > 500 {
		return errs.IssueRemarkLenError
	}
	return nil
}

func CheckUserNameLen(str string) errs.SystemErrorInfo {
	strLen := strs.Len(str)
	if strLen <= 0 || strLen > 64 {
		return errs.UserNameLenError
	}
	return nil
}

//表情转换
func UnicodeEmojiCodeFilter(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			//u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			//ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

func PwdEncrypt(cleartextPassword, salt string) string {
	return md5.Md5V(cleartextPassword + "+" + salt)
}

func PwdEncryptForLesscodeAccoutLogin(loginNamePlusInputPwd, salt string) string {
	return md5.Md5V(loginNamePlusInputPwd + "erp" + salt)
}

//获取MQTT表channel
func GetMQTTTableChannel(orgId, appId, tableId int64) string {
	return temp.RenderIgnoreError(consts.MQTTChannelTable, map[string]interface{}{
		consts.MQTTChannelKeyOrg:   orgId,
		consts.MQTTChannelKeyApp:   appId,
		consts.MQTTChannelKeyTable: tableId,
	})
}

//获取MQTT项目channel
func GetMQTTProjectChannel(orgId, projectId int64) string {
	return temp.RenderIgnoreError(consts.MQTTChannelProject, map[string]interface{}{
		consts.MQTTChannelKeyOrg:     orgId,
		consts.MQTTChannelKeyProject: projectId,
	})
}

//获取MQTT app channel
func GetMQTTAppChannel(orgId, appId int64) string {
	return temp.RenderIgnoreError(consts.MQTTChannelApp, map[string]interface{}{
		consts.MQTTChannelKeyOrg: orgId,
		consts.MQTTChannelKeyApp: appId,
	})
}

//获取MQTT组织channel
func GetMQTTOrgChannel(orgId int64) string {
	return temp.RenderIgnoreError(consts.MQTTChannelOrg, map[string]interface{}{
		consts.MQTTChannelKeyOrg: orgId,
	})
}

//获取MQTT组织channel
func GetMQTTOrgRootChannel(orgId int64) string {
	return temp.RenderIgnoreError(consts.MQTTChannelOrgRoot, map[string]interface{}{
		consts.MQTTChannelKeyOrg: orgId,
	})
}

//获取MQTT用户channel
func GetMQTTUserChannel(orgId, userId int64) string {
	return temp.RenderIgnoreError(consts.MQTTChannelUser, map[string]interface{}{
		consts.MQTTChannelKeyOrg:  orgId,
		consts.MQTTChannelKeyUser: userId,
	})
}

func RoleOperationCodesMatch(operationCode string, operationCodes string) bool {
	codes := strings.Split(operationCodes, "|")
	if len(codes) == 1 {
		return codes[0] == operationCode
	}
	for _, code := range codes {
		if code == "("+operationCode+")" {
			return true
		}
	}
	return false
}

// GetCustomEnv 第三方定制客户（独立部署客户）的环境变量，
func GetCustomEnv() string {
	return os.Getenv(consts.RunCustomEnvKey)
}

// IsChinaMobileEnv 应用当前运行环境是否是“中国移动”的独立部署环境
func IsChinaMobileEnv() bool {
	return GetCustomEnv() == consts.CustomEnvChinaMobile
}
