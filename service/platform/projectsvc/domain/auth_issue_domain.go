package domain

import (
	"fmt"
	"strings"

	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo/appauth"
	"github.com/spf13/cast"
)

// 根据APP AUTH信息校验OPT权限、FIELD权限
func CheckAppOptFieldAuth(payLevel int, tableId int64, operation string, checkAuthFields []string, appAuth appauth.GetAppAuthData) bool {
	// 非付费组织不校验字段权限
	if payLevel != consts.PayLevelStandard {
		if !checkFieldAuth(tableId, operation, checkAuthFields, appAuth) {
			return false
		}
	}

	// 校验操作权限
	optAuth := expandOptAuth(appAuth.OptAuth)
	access, _ := slice.Contain(optAuth, operation)
	return access
}

// 将形如：`["Permission.Pro.Tag-Create,Modify,Delete"]` 的权限数组转换为 operation（格式是：`Permission.Pro.Tag.Create`） 一致的格式。
func expandOptAuth(optAuthArr []string) []string {
	opList := []string{}
	for _, item := range optAuthArr {
		infos := strings.Split(item, "-")
		if len(infos) > 1 {
			opPrev := infos[0]
			opSuffixArr := strings.Split(infos[1], ",")
			for _, oneSuffix := range opSuffixArr {
				opList = append(opList, fmt.Sprintf("%s.%s", opPrev, oneSuffix))
			}
		} else {
			opList = append(opList, item)
		}
	}
	return opList
}

func checkFieldAuth(tableId int64, operation string, checkAuthFields []string, appAuth appauth.GetAppAuthData) bool {
	tableIdStr := cast.ToString(tableId)
	for _, s := range checkAuthFields {
		// 这是系统定义的一些，不需要校验
		// id和issueId在无码里会剔除掉，不会更新，所以不需要校验
		if ok, _ := slice.Contain([]string{"lessUpdateIssueReq", consts.BasicFieldId, consts.BasicFieldIssueId}, s); ok {
			continue
		}
		if operation == consts.OperationProIssue4Create {
			if ok, _ := slice.Contain([]string{
				consts.BasicFieldOwnerId,
				consts.BasicFieldFollowerIds,
				consts.BasicFieldProjectObjectTypeId,
				consts.BasicFieldIssueStatus,
				consts.BasicFieldPriority,
				consts.BasicFieldCreator,
				consts.BasicFieldUpdator,
				consts.BasicFieldIterationId,
				consts.BasicFieldOwnerId,
				consts.BasicFieldProjectId,
				consts.BasicFieldParentId,
			}, s); ok {
				continue
			}
		}
		if !appAuth.HasFieldWriteAuth(tableIdStr, s) {
			log.Infof("[checkFieldAuth] 无权字段：%s", s)
			return false
		}
	}
	return true
}
