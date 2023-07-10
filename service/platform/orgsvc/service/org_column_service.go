package service

import (
	tablev1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/tablefacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
)

func CreateOrgColumn(req orgvo.CreateOrgColumnReq) (*tablev1.CreateOrgColumnReply, errs.SystemErrorInfo) {
	// 校验权限
	userAuthorityResp := userfacade.GetUserAuthority(req.OrgId, req.UserId)
	if userAuthorityResp.Failure() {
		log.Errorf("[CreateOrgColumn.GetUserAuthority] failed, orgId:%v, userId:%v, err:%v", req.OrgId, req.UserId, userAuthorityResp.Error())
		return nil, errs.BuildSystemErrorInfo(errs.UserDomainError, userAuthorityResp.Error())
	}
	isOrgOwner := userAuthorityResp.Data.IsOrgOwner
	isSysAdmin := userAuthorityResp.Data.IsSysAdmin
	isModifyAuth, _ := slice.Contain(userAuthorityResp.Data.OptAuth, consts.OperationOrgConfigModifyField)

	if !isOrgOwner && !isSysAdmin && !isModifyAuth {
		return nil, errs.NoOperationPermissionForOrgColumn
	}

	createOrgColumnResp := tablefacade.CreateOrgColumn(req)
	if createOrgColumnResp.Failure() {
		log.Errorf("[CreateOrgColumn] tablefacade.CreateOrgColumn failed, orgId:%v, userId:%v, err:%v", req.OrgId, req.UserId, createOrgColumnResp.Error())
		return nil, createOrgColumnResp.Error()
	}

	return &tablev1.CreateOrgColumnReply{}, nil
}

func GetOrgColumns(req orgvo.GetOrgColumnsReq) (*orgvo.ReadOrgColumnsReply, errs.SystemErrorInfo) {
	columnsResp := tablefacade.GetOrgColumns(req)
	if columnsResp.Failure() {
		log.Errorf("[GetOrgColumns] tablefacade.GetOrgColumns failed, orgId:%v, userId:%v, err:%v", req.OrgId, req.UserId, columnsResp.Error())
		return nil, columnsResp.Error()
	}
	return columnsResp.Data, nil
}

func DeleteOrgColumn(req orgvo.DeleteOrgColumnReq) (*tablev1.DeleteOrgColumnReply, errs.SystemErrorInfo) {
	// 校验权限
	userAuthorityResp := userfacade.GetUserAuthority(req.OrgId, req.UserId)
	if userAuthorityResp.Failure() {
		log.Errorf("[CreateOrgColumn.GetUserAuthority] failed, orgId:%v, userId:%v, err:%v", req.OrgId, req.UserId, userAuthorityResp.Error())
		return nil, errs.BuildSystemErrorInfo(errs.UserDomainError, userAuthorityResp.Error())
	}
	isOrgOwner := userAuthorityResp.Data.IsOrgOwner
	isSysAdmin := userAuthorityResp.Data.IsSysAdmin
	isModifyAuth, _ := slice.Contain(userAuthorityResp.Data.OptAuth, consts.OperationOrgConfigModifyField)

	if !isOrgOwner && !isSysAdmin && !isModifyAuth {
		return nil, errs.NoOperationPermissionForOrgColumn
	}

	deleteOrgColumnResp := tablefacade.DeleteOrgColumn(req)
	if deleteOrgColumnResp.Failure() {
		log.Errorf("[DeleteOrgColumn] tablefacade.DeleteOrgColumn failed, orgId:%v, userId:%v, err:%v", req.OrgId, req.UserId, deleteOrgColumnResp.Error())
		return nil, deleteOrgColumnResp.Error()
	}
	return deleteOrgColumnResp.Data, nil
}
