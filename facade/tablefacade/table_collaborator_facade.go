package tablefacade

import (
	"fmt"

	v1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade"
)

func CheckIsAppCollaborator(orgId, userId, appId, checkUserId int64) *projectvo.CheckIsAppCollaboratorReply {
	respVo := &projectvo.CheckIsAppCollaboratorReply{Data: &v1.CheckIsAppCollaboratorReply{}}
	reqUrl := fmt.Sprintf("%s%s/rows/isAppCollaborator", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	req := &v1.CheckIsAppCollaboratorRequest{
		AppId:  appId,
		UserId: checkUserId,
	}
	err := facade.RequestWithCommonHeader(orgId, userId, consts.HttpMethodPost, reqUrl, nil, nil, req, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetUserAppCollaboratorRoles(orgId, userId, appId int64) *projectvo.GetUserAppCollaboratorRolesReply {
	respVo := &projectvo.GetUserAppCollaboratorRolesReply{Data: &v1.GetUserAppCollaboratorRolesReply{}}
	reqUrl := fmt.Sprintf("%s%s/rows/userAppCollaboratorRoles", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	req := &v1.GetUserAppCollaboratorRolesRequest{
		AppId:  appId,
		UserId: userId,
	}
	err := facade.RequestWithCommonHeader(orgId, userId, consts.HttpMethodPost, reqUrl, nil, nil, req, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetAppCollaboratorRoles(orgId, userId, appId int64) *projectvo.GetAppCollaboratorRolesReply {
	respVo := &projectvo.GetAppCollaboratorRolesReply{Data: &v1.GetAppCollaboratorRolesReply{}}
	reqUrl := fmt.Sprintf("%s%s/rows/appCollaboratorRoles", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	req := &v1.GetAppCollaboratorRolesRequest{
		AppId: appId,
	}
	err := facade.RequestWithCommonHeader(orgId, userId, consts.HttpMethodPost, reqUrl, nil, nil, req, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetDataCollaborators(orgId, userId int64, dataIds []int64) *projectvo.GetDataCollaboratorsReply {
	respVo := &projectvo.GetDataCollaboratorsReply{Data: &v1.GetDataCollaboratorsReply{}}
	reqUrl := fmt.Sprintf("%s%s/rows/dataCollaborators", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	req := &v1.GetDataCollaboratorsRequest{
		DataIds: dataIds,
	}
	err := facade.RequestWithCommonHeader(orgId, userId, consts.HttpMethodPost, reqUrl, nil, nil, req, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}
