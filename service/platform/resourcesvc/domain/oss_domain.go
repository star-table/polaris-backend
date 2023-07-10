package domain

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

// GetPayFunctionLimitResourceNum 获取功能限制的资源数量。仅适用于：项目、迭代、任务
func GetPayFunctionLimitResourceNum(orgId int64, functionCode string) (int, errs.SystemErrorInfo) {
	resp := orgfacade.GetFunctionObjArrByOrg(orgvo.GetOrgFunctionConfigReq{
		OrgId: orgId,
	})
	if resp.Failure() {
		log.Errorf("[GetPayFunctionLimitResourceNum] GetFunctionObjArrByOrg err: %v", resp.Error())
		return 0, resp.Error()
	}
	payFunctionMap := GetOrgFunctionInfoMap(resp.Data.Functions)
	// map 中不存在 function 表示不支持此项功能；存在，则验证限制数量
	if payFuncInfo, ok := payFunctionMap[functionCode]; ok {
		if payFuncInfo.HasLimit {
			return payFuncInfo.Limit[0].Num, nil
		} else {
			// 无限制
			return -1, nil
		}
	}

	return 0, nil
}

func GetOrgFunctionInfoMap(functions []orgvo.FunctionLimitObj) map[string]orgvo.FunctionLimitObj {
	funcMap := make(map[string]orgvo.FunctionLimitObj, len(functions))
	for _, function := range functions {
		funcMap[function.Key] = function
	}

	return funcMap
}
