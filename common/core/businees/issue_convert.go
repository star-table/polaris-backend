package businees

import (
	"fmt"
	"strings"

	"github.com/star-table/common/core/types"
	"github.com/spf13/cast"

	"github.com/star-table/common/core/logger"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

var log = logger.GetDefaultLogger()

func HandleIssueVoToLc(info vo.HomeIssueInfo) map[string]interface{} {
	temp := map[string]interface{}{}
	for s, i := range info.LessData {
		temp[s] = i
	}
	//for s, i := range sysMap {
	//	temp[s] = i
	//}
	//temp["issueStatus"] = info.Issue.Status
	if dataId, ok := temp["id"]; ok {
		temp["dataId"] = dataId
	}
	temp["id"] = info.IssueID
	temp["issueId"] = info.IssueID
	temp["childsNum"] = info.ChildsNum
	temp["info"] = info

	return temp
}

func HandleIssueBoToLc(info bo.HomeIssueInfoBo) map[string]interface{} {
	temp := map[string]interface{}{}
	for s, i := range info.LessData {
		temp[s] = i
	}
	if dataId, ok := temp["id"]; ok {
		temp["dataId"] = dataId
	}
	temp["id"] = info.IssueId
	temp["issueId"] = info.IssueId
	temp["childsNum"] = info.ChildsNum
	//把确认人审核状态加进来给前端
	temp["auditorsInfo"] = info.AuditorsInfo
	temp["endTime"] = info.Issue.EndTime
	temp[consts.BasicFieldIssueStatus] = info.Issue.Status
	temp[consts.BasicFieldAuditStatus] = info.Issue.AuditStatus

	return temp
}

func FormatUserId(userId int64) string {
	return fmt.Sprintf("%s%d", consts.LcCustomFieldUserType, userId)
}

func FormatDeptId(deptId int64) string {
	return fmt.Sprintf("%s%d", consts.LcCustomFieldDeptType, deptId)
}

//["D_1"]处理成[1]
func LcDeptToDeptIds(depts []string) []int64 {
	deptIds, _, _, _ := LcInterfaceToIds(consts.LcCustomFieldDeptType, depts, true, true)
	return deptIds
}

// 将 ["D_1111", "D_1222"] 转为 [1111, 1222]
func LcDeptToDeptIdsWithError(interfaceDepts interface{}, isIgnoreErrorType ...bool) ([]int64, errs.SystemErrorInfo) {
	deptIds, _, _, err := LcInterfaceToIds(consts.LcCustomFieldDeptType, interfaceDepts, isIgnoreErrorType...)
	return deptIds, err
}

//["U_1"]处理成[1]
func LcMemberToUserIds(members []string) []int64 {
	userIds, _, _, _ := LcInterfaceToIds(consts.LcCustomFieldUserType, members, true, true)
	return userIds
}

// 将 ["U_1111", "U_1222"] 转为 [1111, 1222]
func LcMemberToUserIdsWithError(interfaceMembers interface{}, isIgnoreErrorType ...bool) ([]int64, errs.SystemErrorInfo) {
	userIds, _, _, err := LcInterfaceToIds(consts.LcCustomFieldUserType, interfaceMembers, isIgnoreErrorType...)
	return userIds, err
}

func LcParseIds(idValues interface{}) (string, []int64) {
	var prefix string
	var ids []int64
	var isSlice, hasPrefix bool

	// 先检查成员
	ids, isSlice, hasPrefix, _ = LcInterfaceToIds(consts.LcCustomFieldUserType, idValues, true, false)
	if !isSlice {
		prefix = consts.LcCustomFieldUserType
	} else if hasPrefix {
		prefix = consts.LcCustomFieldUserType
	}

	if prefix == "" {
		// 再检查部门
		prefix = consts.LcCustomFieldDeptType
		ids, _, _, _ = LcInterfaceToIds(consts.LcCustomFieldDeptType, idValues, true, false)
	}

	return prefix, ids
}

func LcInterfaceToIds(prefix string, i interface{}, flags ...bool) ([]int64, bool, bool, errs.SystemErrorInfo) {
	var (
		userIds       = make([]int64, 0)
		err           error
		isSlice       = true
		hasPrefix     = false
		isIgnoreError = false
		isIgnoreZero  = true
	)
	// 第一个flag
	if len(flags) > 0 {
		isIgnoreError = flags[0]
	}
	// 第二个flag
	if len(flags) > 1 {
		isIgnoreZero = flags[1]
	}

	switch members := i.(type) {
	case string:
		isSlice = false
		userIds, hasPrefix, err = convertIdInterfaceToIds(prefix, userIds, members, isIgnoreZero, isIgnoreError)
		if err != nil && isIgnoreError == false {
			return nil, isSlice, hasPrefix, errs.ReqParamsValidateError
		}
	case []interface{}:
		for _, id := range members {
			userIds, hasPrefix, err = convertIdInterfaceToIds(prefix, userIds, id, isIgnoreZero, isIgnoreError)
			if err != nil && isIgnoreError == false {
				return nil, isSlice, hasPrefix, errs.ReqParamsValidateError
			}
		}
	case []string:
		for _, id := range members {
			userIds, hasPrefix, err = convertIdInterfaceToIds(prefix, userIds, id, isIgnoreZero, isIgnoreError)
			if err != nil && isIgnoreError == false {
				return nil, isSlice, hasPrefix, errs.ReqParamsValidateError
			}
		}
	case nil:
		return userIds, isSlice, hasPrefix, nil
	default:
		if isIgnoreError {
			return userIds, isSlice, hasPrefix, nil
		}
		return nil, isSlice, hasPrefix, errs.ReqParamsValidateError
	}

	return userIds, isSlice, hasPrefix, nil
}

func convertIdInterfaceToIds(prefix string, ids []int64, id interface{}, isIgnoreZero, isIgnoreError bool) ([]int64, bool, error) {
	idStr := cast.ToString(id)
	hasPrefix := strings.Contains(idStr, prefix)
	replaceStr := strings.Replace(idStr, prefix, "", 1)
	idInt, err := cast.ToInt64E(replaceStr)
	if err != nil {
		if isIgnoreError {
			return ids, hasPrefix, nil
		} else {
			return nil, hasPrefix, errs.ReqParamsValidateError
		}
	} else {
		if idInt == 0 && isIgnoreZero {
			return ids, hasPrefix, nil
		} else {
			return append(ids, idInt), hasPrefix, nil
		}
	}
}

// [1] 处理成["U_1"]
func LcMemberToUserIdStrings(users []int64) []string {
	res := []string{}
	if users == nil {
		return res
	}
	for _, u := range users {
		if u == 0 {
			continue
		}
		res = append(res, fmt.Sprintf("U_%d", u))
	}
	return res
}

// 两个int64类型切片的差集
func DifferenceInt64Set(a []int64, b []int64) []int64 {
	var c []int64
	temp := map[int64]struct{}{}

	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			c = append(c, val)
		}
	}

	return c
}

func GetTimeString(t types.Time) string {
	format := t.String()
	if format == consts.BlankTime || format == consts.BlankEmptyTime {
		return consts.ShowTrendsEmptyTime
	}
	return format
}
