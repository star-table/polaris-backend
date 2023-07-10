package formfacade

import (
	"fmt"

	"github.com/star-table/polaris-backend/facade"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
)

//新增任务
func LessCreateIssue(req formvo.LessCreateIssueReq) *formvo.FormCreateOneResp {
	respVo := &formvo.FormCreateOneResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

//删除任务
func LessDeleteIssue(req formvo.LessDeleteIssueReq) *formvo.LessCommonIssueResp {
	respVo := &formvo.LessCommonIssueResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values/delete", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// LessUpdateIssue 更新无码表数据，如：任务数据、项目数据
func LessUpdateIssue(req formvo.LessUpdateIssueReq) *formvo.LessUpdateIssueResp {
	respVo := &formvo.LessUpdateIssueResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	//创建人和确认人int64转string(无码需求)
	for i, m := range req.Form {
		for s, i2 := range m {
			if ok, _ := slice.Contain([]string{consts.BasicFieldCreator, consts.BasicFieldUpdator}, s); ok {
				req.Form[i][s] = fmt.Sprintf("%v", i2)
			}
		}
	}
	err := facade.Request(consts.HttpMethodPut, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// LessUpdateIssueBatchRaw 批量更新无码任务表数据
func LessUpdateIssueBatchRaw(req *formvo.LessUpdateIssueBatchReq) *formvo.LessCommonIssueResp {
	respVo := &formvo.LessCommonIssueResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values/updateBatchRaw", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

//获取任务列表
//func LessIssueList(req formvo.LessIssueListReq) *formvo.LessIssueListResp {
//	respVo := &formvo.LessIssueListResp{}
//	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values/filter", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
//	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
//	if err.Failure() {
//		respVo.Err = err
//	}
//	return respVo
//}

// LessIssueRawList 不处理的裸数据
//func LessIssueRawList(req formvo.LessIssueListReq) *formvo.LessIssueRawListResp {
//	respVo := &formvo.LessIssueRawListResp{}
//	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values/filterRaw", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
//	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
//	if err.Failure() {
//		respVo.Err = err
//	}
//	return respVo
//}

//获取任务列表
//func LessFilterStat(req formvo.LessIssueListReq) *formvo.LessFilterStatResp {
//	respVo := &formvo.LessFilterStatResp{}
//	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values/filterStat", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
//	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
//	if err.Failure() {
//		respVo.Err = err
//	}
//	return respVo
//}

//获取任务列表
//func LessFilterCustomStat(req formvo.LessIssueListReq) *formvo.LessFilterCustomStatResp {
//	respVo := &formvo.LessFilterCustomStatResp{}
//	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values/filterCustomStat", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
//	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
//	if err.Failure() {
//		respVo.Err = err
//	}
//	return respVo
//}

//回收任务
func LessRecycleIssue(req formvo.LessRecycleIssueReq) *formvo.LessRecycleIssueResp {
	respVo := &formvo.LessRecycleIssueResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values/recycle", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	err := facade.Request(consts.HttpMethodPut, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

//恢复任务
func LessRecoverIssue(req formvo.LessRecoverIssueReq) *formvo.LessCommonIssueResp {
	respVo := &formvo.LessCommonIssueResp{}
	reqUrl := fmt.Sprintf("%s/form/inner/api/v1/apps/%d/values/recover", config.GetPreUrl(consts.ServiceNameLcForm), req.AppId)
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}
