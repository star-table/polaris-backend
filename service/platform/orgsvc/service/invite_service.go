package service

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/uuid"
	consts2 "github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/core/util/rand"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"github.com/nyaruka/phonenumbers"
	"github.com/tealeg/xlsx/v2"
)

func GetInviteCode(currentUserId int64, orgId int64, sourcePlatform string) (*orgvo.GetInviteCodeRespVoData, errs.SystemErrorInfo) {
	//用户角色权限校验
	authErr := AuthOrgRole(orgId, currentUserId, consts2.RoleOperationPathOrgUser, consts2.OperationOrgInviteUserInvite)
	if authErr != nil {
		log.Error(authErr)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, authErr)
	}
	inviteInfo := bo.InviteInfoBo{
		InviterId:      currentUserId,
		OrgId:          orgId,
		SourcePlatform: sourcePlatform,
	}
	// 产品、测试：如果是飞书来源，则不能在极星生成邀请码。
	//forbidGenInviteCodeSourceList := []string{sdk_const.SourceChannelFeishu}
	//if exist, _ := slice.Contain(forbidGenInviteCodeSourceList, sourcePlatform); exist {
	//	return nil, errs.ForbidInviteSourcePlatform
	//}

	inviteCode := rand.RandomInviteCode(uuid.NewUuid() + strconv.FormatInt(currentUserId, 10) + sourcePlatform)
	err := domain.SetUserInviteCodeInfo(inviteCode, inviteInfo)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, err
	}
	return &orgvo.GetInviteCodeRespVoData{InviteCode: inviteCode, Expire: consts.CacheUserInviteCodeExpire}, nil
}

// InviteUserByPhones 邀请-手机号批量邀请
// 邀请链接通过前端传入（实际上是后端生成的）
func InviteUserByPhones(orgId, userId int64, input orgvo.InviteUserByPhonesReqVoData) (*orgvo.InviteUserByPhonesRespVoData, errs.SystemErrorInfo) {
	errList := make([]orgvo.InviteUserByPhonesRespVoDataErrItem, 0)
	successList := make([]orgvo.InviteUserByPhonesRespVoDataSuccessItem, 0)
	result := &orgvo.InviteUserByPhonesRespVoData{
		ErrorList:   errList,
		SuccessList: successList,
	}
	orgBo, err := GetOrgInfoBo(orgId)
	if err != nil {
		return result, err
	}
	phones := input.Phones
	if len(phones) < 1 {
		err := errs.ParamError
		log.Errorf("Phones 参数错误：%s", err)
		return nil, err
	}
	if len(input.InviteCode) < 3 {
		err := errs.ParamError
		log.Errorf("InviteCode 参数错误: %s", err)
		return nil, err
	}
	// 检查手机号是否已经存在系统内
	loginPhoneList := make([]string, len(phones))
	for i, phone := range phones {
		loginPhoneList[i] = fmt.Sprintf("%s-%s", phone.Origin, phone.Number)
	}
	userBoList, err := domain.GetExistUserByPhones(orgId, loginPhoneList)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	userBoMap := make(map[string]bo.UserInfoBo, len(userBoList))
	for _, userBo := range userBoList {
		userBoMap[userBo.Mobile] = userBo
	}
	// 键值对，键是手机号，值是发送短信对应的模板参数。{phoneNum: tplParam}
	couldSendPhoneInfoMap := make(map[string]map[string]string, 0)
	inviteCode := input.InviteCode
	smsType := consts2.AuthCodeTypeInviteUserByPhones
	// 向所有手机号发送邀请短信
	for _, phone := range phones {
		// 组装手机号
		phoneNumber := fmt.Sprintf("%s-%s", phone.Origin, phone.Number)
		if phone.Origin == "+86" {
			if !format.VerifyChinaPhoneWithoutAreaFormat(phone.Number) {
				errList = append(errList, orgvo.InviteUserByPhonesRespVoDataErrItem{
					Number: phone.Number,
					Reason: consts2.ImportUserErrOfPhoneFormatErr,
				})
				continue
			}
		}
		// 检查是否已存在
		if _, ok := userBoMap[phoneNumber]; ok {
			errList = append(errList, orgvo.InviteUserByPhonesRespVoDataErrItem{
				Number: phone.Number,
				Reason: consts2.ImportUserErrOfPhoneExistErr,
			})
			continue
		}

		// 发短信频率限制检查
		limitErr := domain.CheckSMSLoginCodeFreezeTime(smsType, 1, phone.Number)
		if limitErr != nil {
			errList = append(errList, orgvo.InviteUserByPhonesRespVoDataErrItem{
				Number: phone.Number,
				Reason: consts2.ImportUserErrOfPhoneSendLimitErr,
			})
			continue
		}

		// 在白名单中，不发送短信
		if !IsInWhiteList(phone.Number) {
			tplParam := map[string]string{
				consts2.SMSParamsNameOrgName:    orgBo.Name,
				consts2.SMSParamsNameInviteCode: inviteCode,
			}
			couldSendPhoneInfoMap[phone.Number] = tplParam
		}
		successList = append(successList, orgvo.InviteUserByPhonesRespVoDataSuccessItem{
			Number: phone.Number,
		})
	}
	// 协程池处理异步发送短信
	if len(couldSendPhoneInfoMap) > 0 {
		//go func() {
		//	defer func() {
		//		if r := recover(); r != nil {
		//			log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
		//		}
		//	}()
		//	pool := go_pool.NewGoPool(4)
		//	for phoneNum, tplParam := range couldSendPhoneInfoMap {
		//		pool.NewTask(func() {
		//			sendErr := sendSmsByTplWithParam(smsType, phoneNum, tplParam)
		//			if sendErr != nil {
		//				log.Error(sendErr)
		//			}
		//
		//			// 记录已发送，缓存 1min
		//			// 参数 addressType，1：手机号。根据 service/platform/orgsvc/service/send_auth_code_service.go:31。
		//			setFreezeErr := domain.SetSMSLoginCodeFreezeTime(smsType, 1, phoneNum, 1)
		//			if setFreezeErr != nil {
		//				//这里不要影响主流程
		//				log.Error(setFreezeErr)
		//			}
		//		})
		//	}
		//}()
	}
	result.ErrorList = errList
	result.SuccessList = successList

	return result, nil
}

// ExportInviteTemplate 邀请-下载导入成员模板
func ExportInviteTemplate(orgId, userId int64) (*orgvo.ExportInviteTemplateRespVoData, errs.SystemErrorInfo) {
	result := &orgvo.ExportInviteTemplateRespVoData{
		Url: "https://attachments.startable.cn/org_14396/project_29713/issue_911145/resource/2021/12/17/0352589654914a1caf046e38d1846ad81639729525600.xlsx?attname=%E6%88%90%E5%91%98%E5%AF%BC%E5%85%A5%E6%A8%A1%E6%9D%BF.xlsx",
	}
	//relatePath, fileDir, fileName, err := GetExportFileInfo(orgId, "成员导入模板.xlsx")
	//if err != nil {
	//	log.Error(err)
	//	return result, err
	//}
	//GenExcelFileForInviteTemplate(fileDir+"/"+fileName)
	//result.Url = config.GetOSSConfig().LocalDomain + relatePath + "/" + fileName

	return result, nil
}

// GenExcelFileForImportMemberErr 邀请-导入成员时，失败数据的生成和导出
func GenExcelFileForImportMemberErr(fullFileName string, errList []orgvo.ImportMembersResultErrItem) errs.SystemErrorInfo {
	// 两列：姓名 | 手机号码 | 错误原因
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}
	row = sheet.AddRow()

	cell = row.AddCell()
	cell.Value = "姓名"
	cell = row.AddCell()
	cell.Value = "国家代码"
	cell = row.AddCell()
	cell.Value = "手机号码"
	cell = row.AddCell()
	cell.Value = "错误原因"

	// 错误数据列表
	for _, info := range errList {
		row = sheet.AddRow()

		cell = row.AddCell()
		cell.Value = info.Name
		cell = row.AddCell()
		cell.Value = info.Region
		cell = row.AddCell()
		cell.Value = info.Phone
		cell = row.AddCell()
		cell.Value = info.Reason
	}

	err = file.Save(fullFileName)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.InviteImportTplGenExcelErr, err)
	}

	return nil
}

// GetExportFileInfo 获取导出的文件相关信息
// fileName string 导出的文件名，如：“成员导入模板.xlsx”
func GetExportFileInfo(orgId int64, fileName string) (relatePath, fileDir, outFileName string, err errs.SystemErrorInfo) {
	curMonthlyDate := time.Now().Format("200601")
	relatePath = "/import_member/org_" + strconv.FormatInt(orgId, 10) + "/date_" + curMonthlyDate
	// log.Info(config.GetOSSConfig())
	fileDir = strings.TrimRight(config.GetOSSConfig().RootPath, "/") + relatePath
	mkdirErr := os.MkdirAll(fileDir, 0777)
	if mkdirErr != nil {
		log.Error(mkdirErr)
		return "", "", "", errs.BuildSystemErrorInfo(errs.IssueDomainError, mkdirErr)
	}
	outFileName = fileName
	return relatePath, fileDir, outFileName, nil
}

// ImportMembers 邀请-excel 导入组织成员
func ImportMembers(orgId, userId int64, input *orgvo.ImportMembersReqVoData) (*orgvo.ImportMembersRespVoData, errs.SystemErrorInfo) {
	result := &orgvo.ImportMembersRespVoData{
		ErrCount:     0,
		SuccessCount: 0,
		SkipCount:    0,
		ErrExcelUrl:  "",
	}
	successCount := 0
	skipPhones := []string{}
	errList := make([]orgvo.ImportMembersResultErrItem, 0)
	// 直接获取导入的用户列表数据
	importList := input.ImportUserList
	if len(importList) < 1 {
		return result, nil
	}
	if len(importList) > 200 {
		err := errs.ImportUserNumTooMany
		log.Error(err)
		return result, err
	}
	// 查询已经存在的用户
	userPhoneList := make([]string, len(importList))
	for i, user := range importList { // 目前，在表中，login_name 和 mobile 字段存的格式是 `+86-15010011001`
		phoneNumber := fmt.Sprintf("%s-%s", user.Origin, user.Phone)
		userPhoneList[i] = phoneNumber
	}
	existUserBoMap := make(map[string]bo.UserInfoBo, 0)
	existUserBoList, err := domain.GetExistUserByPhones(orgId, userPhoneList)
	if err != nil && err.Code() != errs.UserNotExist.Code() {
		log.Error(err)
		return result, err
	}
	// 保存有异常的手机号
	errPhones := make([]string, 0, 10)
	existUserPhones := make([]string, len(existUserBoList))
	for _, userBo := range existUserBoList {
		errPhones = append(errPhones, userBo.Mobile)
		existUserPhones = append(existUserPhones, userBo.Mobile)
		existUserBoMap[userBo.Mobile] = userBo
	}
	// 检查手机号重复出现（表单内重复）的情况。
	newList := make([]orgvo.ImportMembersReqVoDataUserItem, 0)
	checkExistMap := make(map[string]bool, 0)
	for _, user := range importList {
		if _, exist := checkExistMap[user.Phone]; exist {
			errList = append(errList, orgvo.ImportMembersResultErrItem{
				Name:   user.Name,
				Region: user.Origin,
				Phone:  user.Phone,
				Reason: consts2.ImportUserErrOfPhoneMultipleErr,
			})
			phoneNumber := fmt.Sprintf("%s-%s", user.Origin, user.Phone)
			errPhones = append(errPhones, phoneNumber)
			continue
		} else {
			checkExistMap[user.Phone] = true
			newList = append(newList, user)
		}
	}
	toImportLoginPhonesMap := make(map[string]orgvo.ImportMembersReqVoDataUserItem, 0)
	filteredLoginNames := make([]string, 0)
	importList = newList
	// 直接通过信息注册用户
	for i, user := range importList {
		// 目前，在表中，login_name 和 mobile 字段存的格式是 `+86-15010011001`
		phoneNumber := fmt.Sprintf("%s-%s", user.Origin, user.Phone)
		phoneInvalid, errNum := phonenumbers.Parse(user.Origin+user.Phone, "")
		if errNum != nil || !phonenumbers.IsValidNumber(phoneInvalid) {
			errList = append(errList, orgvo.ImportMembersResultErrItem{
				Name:   user.Name,
				Region: user.Origin,
				Phone:  user.Phone,
				Reason: consts2.ImportUserErrOfPhoneFormatErr,
			})
			errPhones = append(errPhones, phoneNumber)
			continue
		}

		// 校验名字，名字如果为空，则跳过
		if user.Name == "" {
			// 根据手机号生成一个默认的用户名
			importList[i].Name = domain.GenDefaultNameByPhone(user.Phone)
			user.Name = importList[i].Name
		}
		if !format.VerifyUserNameFormat(user.Name) {
			errList = append(errList, orgvo.ImportMembersResultErrItem{
				Name:   user.Name,
				Region: user.Origin,
				Phone:  user.Phone,
				Reason: errs.UserNameLenError.Message(),
			})
			errPhones = append(errPhones, phoneNumber)
			continue
		}

		// 已经存在的，跳过
		if exist, _ := slice.Contain(existUserPhones, phoneNumber); exist {
			if _, ok := existUserBoMap[phoneNumber]; ok {
				//errList = append(errList, orgvo.ImportMembersResultErrItem{
				//	Name:   user.Name,
				//	Region: user.Origin,
				//	Phone:  user.Phone,
				//	Reason: consts2.ImportUserErrOfPhoneExistErr,
				//})
				//errPhones = append(errPhones, user.Phone)
				skipPhones = append(skipPhones, phoneNumber)
			}
			continue
		}
		if exist, _ := slice.Contain(errPhones, phoneNumber); exist {
			continue
		}
		log.Infof("[ImportMembers] 用户%s未注册，phone: %s，开始注册....", user.Name, user.Phone)
		filteredLoginNames = append(filteredLoginNames, phoneNumber)
		toImportLoginPhonesMap[phoneNumber] = user
	}

	// 有两种可能：1.user user_org 都不存在手机号对应的记录；2.user 已存在，但是 user_org 不存在，此时需要增加组织的关联
	filteredLoginNames = slice.SliceUniqueString(filteredLoginNames)
	needRegisterLoginNames, needCreateOrgUserRelationUserIdsMap, needResetCheckStatusUserIdsMap, err := domain.DetectUserInfoInUser(orgId, filteredLoginNames)
	if err != nil {
		log.Error(err)
		return result, err
	}

	needCreateOrgUserRelationPhones := []string{}
	for _, phone := range needCreateOrgUserRelationUserIdsMap {
		needCreateOrgUserRelationPhones = append(needCreateOrgUserRelationPhones, phone)
	}

	needResetCheckStatusPhones := []string{}
	for _, phone := range needResetCheckStatusUserIdsMap {
		needResetCheckStatusPhones = append(needResetCheckStatusPhones, phone)
	}

	notSkipPhones := []string{}
	notSkipPhones = append(notSkipPhones, needCreateOrgUserRelationPhones...)
	notSkipPhones = append(notSkipPhones, needRegisterLoginNames...)
	notSkipPhones = append(notSkipPhones, needResetCheckStatusPhones...)

	//needUpdateNameUserIds := []int64{}

	if len(needRegisterLoginNames) > 0 {
		for _, loginName := range needRegisterLoginNames {
			phoneNumber := loginName
			user := orgvo.ImportMembersReqVoDataUserItem{}
			if tmpUser, ok := toImportLoginPhonesMap[phoneNumber]; ok {
				user = tmpUser
			}
			tmpOrgBo, err := domain.UserRegister(bo.UserSMSRegisterInfo{
				PhoneNumber:    phoneNumber,
				SourceChannel:  consts2.AppSourceChannelWeb,
				SourcePlatform: consts2.AppSourceChannelWeb,
				Name:           user.Name,
				InviteCode:     "",
				MobileRegion:   user.Origin,
			})
			if err != nil {
				errList = append(errList, orgvo.ImportMembersResultErrItem{
					Name:   user.Name,
					Region: user.Origin,
					Phone:  user.Phone,
					Reason: err.Message(),
				})
				errPhones = append(errPhones, phoneNumber)
				continue
			}
			// 更新新注册的用户的信息
			err = ImportMemberRegisterHandleHook(orgvo.ImportUserInfo{
				OperateUid:    userId,
				OrgId:         orgId,
				SourceChannel: consts2.AppSourceChannelWeb, // 只有 web 来源的组织可以执行该导入功能
			}, *tmpOrgBo)
			if err != nil {
				errList = append(errList, orgvo.ImportMembersResultErrItem{
					Name:   user.Name,
					Region: user.Origin,
					Phone:  user.Phone,
					Reason: err.Message(),
				})
				errPhones = append(errPhones, phoneNumber)
				continue
			}
			successCount += 1
		}
	}
	if len(needCreateOrgUserRelationUserIdsMap) > 0 {
		// 查询 user
		needUserIds := make([]int64, 0, len(needCreateOrgUserRelationUserIdsMap))
		for id := range needCreateOrgUserRelationUserIdsMap {
			needUserIds = append(needUserIds, id)
		}
		//needUpdateNameUserIds = append(needUpdateNameUserIds, needUserIds...)
		userPoList, err := domain.BatchGetUserDetailInfo(needUserIds)
		if err != nil {
			log.Error(err)
			return result, err
		}

		// user 已存在，只需要增加用户和组织的关联
		for _, userPo := range userPoList {
			importUserItem := orgvo.ImportMembersReqVoDataUserItem{}
			if tmpUser, ok := toImportLoginPhonesMap[needCreateOrgUserRelationUserIdsMap[userPo.ID]]; ok {
				importUserItem = tmpUser
			}

			err = domain.AddOrgMember(orgId, userPo.ID, userId, false, false)
			if err != nil {
				errList = append(errList, orgvo.ImportMembersResultErrItem{
					Name:   importUserItem.Name,
					Region: importUserItem.Origin,
					Phone:  importUserItem.Phone,
					Reason: err.Message(),
				})
				phoneNumber := fmt.Sprintf("%s-%s", importUserItem.Origin, importUserItem.Phone)
				errPhones = append(errPhones, phoneNumber)
				continue
			}
			successCount += 1
		}
	}
	if len(needResetCheckStatusUserIdsMap) > 0 {
		// 需要把审核不通过的用户 审核状态重置
		memberIds := make([]int64, 0, len(needResetCheckStatusUserIdsMap))
		for id := range needResetCheckStatusUserIdsMap {
			memberIds = append(memberIds, id)
		}
		//needUpdateNameUserIds = append(needUpdateNameUserIds, memberIds...)
		userPoList, errSys := domain.BatchGetUserDetailInfo(memberIds)
		if errSys != nil {
			log.Error(err)
			return result, err
		}
		errSys = domain.ModifyOrgMemberCheckStatus(orgId, memberIds, consts2.AppCheckStatusSuccess, userId, true)
		if errSys != nil {
			log.Errorf("[CreateOrgUser] ModifyOrgMemberCheckStatus err:%v, orgId:%v, operatorId:%v, memberIds:%v", errSys, orgId, userId, memberIds)
			for _, user := range userPoList {
				importUserItem := orgvo.ImportMembersReqVoDataUserItem{}
				if tmpUser, ok := toImportLoginPhonesMap[needResetCheckStatusUserIdsMap[user.ID]]; ok {
					importUserItem = tmpUser
				}
				errList = append(errList, orgvo.ImportMembersResultErrItem{
					Name:   importUserItem.Name,
					Region: importUserItem.Origin,
					Phone:  importUserItem.Phone,
					Reason: errSys.Message(),
				})
			}
		}
	}

	// 导入的姓名要覆盖 userId对应的其他本地团队
	//needUpdateNameUserIds = slice.SliceUniqueInt64(needUpdateNameUserIds)
	//userOrgIdsMap, errSys := domain.GetLocalOrgUserIdMap(needUpdateNameUserIds)
	//if errSys != nil {
	//	log.Errorf("[ImportMembers] GetLocalOrgUserIdMap err:%v, orgId:%v, phoneNumber:%v", errSys, orgId)
	//	return result, errSys
	//}
	phoneNameMap := map[string]string{}
	for _, data := range importList {
		phoneNumber := fmt.Sprintf("%s-%s", data.Origin, data.Phone)
		phoneNameMap[phoneNumber] = data.Name
	}
	userIdNameMap := map[int64]string{}
	for uId, phone := range needCreateOrgUserRelationUserIdsMap {
		userIdNameMap[uId] = phoneNameMap[phone]
	}
	for uId, phone := range needResetCheckStatusUserIdsMap {
		userIdNameMap[uId] = phoneNameMap[phone]
	}
	//for uId, orgIds := range userOrgIdsMap {
	//	if name, ok := userIdNameMap[uId]; ok {
	//		errSys = domain.UpdateLocalOrgUserNames(orgIds, uId, name)
	//		if errSys != nil {
	//			log.Errorf("[ImportMembers] UpdateLocalOrgUserNames err:%v, orgId:%v, userId:%v", errSys, orgId, userId)
	//			return result, errSys
	//		}
	//	}
	//}

	notSkipPhones = append(notSkipPhones, errPhones...)
	notSkipPhones = slice.SliceUniqueString(notSkipPhones)
	for _, user := range importList {
		phoneNumber := fmt.Sprintf("%s-%s", user.Origin, user.Phone)
		if ok, errSlice := slice.Contain(notSkipPhones, phoneNumber); errSlice == nil && !ok {
			skipPhones = append(skipPhones, phoneNumber)
		}
	}

	// 通过 errList 生成 excel
	relatePath, fileDir, fileName, err := GetExportFileInfo(orgId, "成员导入异常信息.xlsx")
	if err != nil {
		log.Error(err)
		return result, err
	}
	if len(errList) > 0 {
		if err = GenExcelFileForImportMemberErr(fileDir+"/"+fileName, errList); err != nil {
			log.Error(err)
			return result, err
		}
		result.ErrExcelUrl = config.GetOSSConfig().LocalDomain + relatePath + "/" + fileName
	}
	skipPhones = slice.SliceUniqueString(skipPhones)
	result.ErrCount = len(errList)
	result.SuccessCount = len(input.ImportUserList) - len(errList) - len(skipPhones)
	result.SkipCount = len(skipPhones)

	return result, nil
}

func GetInviteInfo(inviteCode string) (*vo.GetInviteInfoResp, errs.SystemErrorInfo) {
	inviteInfo, err := domain.GetUserInviteCodeInfo(inviteCode)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, err
	}

	orgBaseInfo, err := domain.GetBaseOrgInfo(inviteInfo.OrgId)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, err
	}
	userBaseInfo, err := domain.GetBaseUserInfo(inviteInfo.OrgId, inviteInfo.InviterId)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, err
	}

	return &vo.GetInviteInfoResp{
		OrgID:       orgBaseInfo.OrgId,
		OrgName:     orgBaseInfo.OrgName,
		InviterID:   userBaseInfo.UserId,
		InviterName: userBaseInfo.Name,
	}, nil
}

func AddUser(req *orgvo.AddUserReqVo) errs.SystemErrorInfo {
	if !format.VerifyUserNameFormat(req.Input.Name) {
		return errs.UserNameLenError
	}
	phoneInvalid, errNum := phonenumbers.Parse(req.Input.Origin+req.Input.Phone, "")
	if errNum != nil || !phonenumbers.IsValidNumber(phoneInvalid) {
		return errs.MobileInvalidError
	}

	//做登录和自动注册逻辑
	phoneNumber := fmt.Sprintf("%s-%s", req.Input.Origin, req.Input.Phone)
	userBo, err := domain.GetUserInfoByMobile(phoneNumber)
	if err != nil {
		if err == errs.UserNotExist {
			log.Infof("用户%s未注册，开始注册....", phoneNumber)
			userBo, err = domain.UserRegister(bo.UserSMSRegisterInfo{
				PhoneNumber:    phoneNumber,
				SourceChannel:  consts2.AppSourceChannelWeb,
				SourcePlatform: consts2.AppSourceChannelWeb,
				Name:           req.Input.Name,
				MobileRegion:   req.Input.Origin,
			})
			if err != nil {
				log.Errorf("[AddUser] UserRegister input:%v, err:%v", req, err)
				return err
			}
			err = domain.AddOrgMember(req.OrgId, userBo.ID, req.UserId, false, false)
			if err != nil {
				log.Errorf("[AddUser] AddOrgMember orgId:%v, userId:%v, err:%v", req.OrgId, userBo.ID, err)
				return err
			}
			return nil
		} else {
			log.Errorf("[AddUser] GetUserInfoByMobile input:%v, err:%v", req, err)
			return errs.MysqlOperateError
		}
	}

	return errs.MobileSameError
}
