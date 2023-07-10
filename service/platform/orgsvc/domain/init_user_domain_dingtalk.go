package domain

//import (
//	"github.com/star-table/common/core/util/copyer"
//	"github.com/star-table/common/core/util/md5"
//	"github.com/star-table/common/core/util/uuid"
//	"github.com/star-table/common/library/db/mysql"
//	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
//	"github.com/star-table/polaris-backend/common/core/consts"
//	"github.com/star-table/polaris-backend/common/core/errs"
//	"github.com/star-table/polaris-backend/common/core/util/pinyin"
//	"github.com/star-table/polaris-backend/common/extra/dingtalk"
//	"github.com/star-table/polaris-backend/common/model/bo"
//	"github.com/star-table/polaris-backend/facade/idfacade"
//	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
//	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
//	"upper.io/db.v3/lib/sqlbuilder"
//)
//
//func InitDingTalkUser(orgId int64, corpId string, openId string, tx sqlbuilder.Tx) (*bo.UserInfoBo, errs.SystemErrorInfo) {
//	topDep, err1 := GetTopDepartmentInfo(orgId)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, err1
//	}
//
//	client, err2 := dingtalk.GetDingTalkClientRest(corpId)
//	if err2 != nil {
//		log.Error(err2)
//		return nil, errs.BuildSystemErrorInfo(errs.DingTalkClientError, err2)
//	}
//
//	userIdResp, err2 := client.GetUserIdByUnionId(openId)
//	if err2 != nil {
//		log.Error(err2)
//		return nil, errs.DingTalkOpenApiCallError
//	}
//	if userIdResp.ErrCode != 0 {
//		log.Error(userIdResp.ErrMsg)
//		return nil, errs.DingTalkOpenApiCallError
//	}
//
//	userId := userIdResp.UserId
//
//	userDetailResp, err2 := client.GetUserDetail(userId, nil)
//	if err2 != nil {
//		log.Error(err2)
//		return nil, errs.DingTalkOpenApiCallError
//	}
//	if userDetailResp.ErrCode != 0 {
//		log.Error(userDetailResp.ErrMsg)
//		return nil, errs.DingTalkOpenApiCallError
//	}
//
//	user := userDetailResp.UserList
//
//	userNativeId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUser)
//	if idErr != nil {
//		log.Error(idErr)
//		return nil, idErr
//	}
//	userOutInfoId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOutInfo)
//	if idErr != nil {
//		log.Error(idErr)
//		return nil, idErr
//	}
//	userConfigId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserConfig)
//	if idErr != nil {
//		log.Error(idErr)
//		return nil, idErr
//	}
//	userOrgId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOrganization)
//	if idErr != nil {
//		log.Error(idErr)
//		return nil, idErr
//	}
//	userDepId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserDepartment)
//	if idErr != nil {
//		log.Error(idErr)
//		return nil, idErr
//	}
//	userPo := assemblyDingTalkUserInfo(orgId, user)
//	userPo.Id = userNativeId
//	userOutPo := assemblyDingTalkUserOutInfo(orgId, userNativeId, user)
//	userOutPo.Id = userOutInfoId
//	userConfigPo := assemblyOrgUserConfigInfo(orgId, userConfigId, userNativeId)
//	userConfigPo.Id = userConfigId
//	userOrgRelationPo := assemblyUserOrgRelationInfo(orgId, userOrgId, userNativeId)
//	userOrgRelationPo.Id = userOrgId
//	userDepRelationPo := assemblyUserDepRelationInfo(orgId, userNativeId, topDep.Id)
//	userDepRelationPo.Id = userDepId
//
//	dbErr := mysql.TransInsert(tx, &userPo)
//	if dbErr != nil {
//		log.Error(dbErr)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//	dbErr = mysql.TransInsert(tx, &userOutPo)
//	if dbErr != nil {
//		log.Error(dbErr)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//	dbErr = mysql.TransInsert(tx, &userConfigPo)
//	if dbErr != nil {
//		log.Error(dbErr)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//	dbErr = mysql.TransInsert(tx, &userOrgRelationPo)
//	if dbErr != nil {
//		log.Error(dbErr)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//	dbErr = mysql.TransInsert(tx, &userDepRelationPo)
//	if dbErr != nil {
//		log.Error(dbErr)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//
//	userBo := &bo.UserInfoBo{}
//	_ = copyer.Copy(userPo, userBo)
//	return userBo, nil
//}
//
//func assemblyDingTalkUserInfo(orgId int64, duser sdk.UserList) po.PpmOrgUser {
//	phoneNumber := ""
//	sourceChannel := sdk_const.SourceChannelDingTalk
//	sourcePlatform := sdk_const.SourceChannelDingTalk
//	name := duser.Name
//
//	pwd := uuid.NewUuid()
//	salt := uuid.NewUuid()
//	pwd = md5.Md5V(salt + pwd)
//	userPo := &po.PpmOrgUser{
//		OrgId:              orgId,
//		Name:               name,
//		NamePinyin:         pinyin.ConvertToPinyin(name),
//		Avatar:             duser.Avatar,
//		LoginName:          phoneNumber, //
//		LoginNameEditCount: 0,
//		Email:              "",
//		Mobile:             phoneNumber,
//		Password:           pwd,
//		PasswordSalt:       salt,
//		SourceChannel:      sourceChannel,
//		SourcePlatform:     sourcePlatform,
//	}
//	return *userPo
//}
//
//func assemblyDingTalkUserOutInfo(orgId int64, userId int64, duser sdk.UserList) po.PpmOrgUserOutInfo {
//	pwd := uuid.NewUuid()
//	salt := uuid.NewUuid()
//	pwd = md5.Md5V(salt + pwd)
//	userOutInfo := &po.PpmOrgUserOutInfo{}
//	userOutInfo.UserId = userId
//	userOutInfo.OrgId = orgId
//	userOutInfo.OutOrgUserId = duser.UserId
//	userOutInfo.OutUserId = duser.UnionId
//	userOutInfo.IsDelete = consts.AppIsNoDelete
//	userOutInfo.Status = consts.AppStatusEnable
//	userOutInfo.SourceChannel = sdk_const.SourceChannelDingTalk
//	userOutInfo.Name = duser.Name
//	userOutInfo.Avatar = duser.Avatar
//	userOutInfo.JobNumber = duser.JobNumber
//
//	return *userOutInfo
//}
