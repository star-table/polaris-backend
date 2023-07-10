package service

import (
	"os"
	"strconv"
	"time"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/date"
	"github.com/star-table/common/core/util/id/snowflake"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/temp"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/format"
	lang2 "github.com/star-table/polaris-backend/common/core/util/lang"
	int642 "github.com/star-table/polaris-backend/common/core/util/slice/int64"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	"github.com/star-table/polaris-backend/facade/msgfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"github.com/tealeg/xlsx/v2"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//更新成员状态
func UpdateOrgMemberStatus(reqVo orgvo.UpdateOrgMemberStatusReq) (*vo.Void, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	userId := reqVo.UserId
	input := reqVo.Input

	//校验当前用户是否具有修改成员状态的权限
	authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgUser, consts.OperationOrgUserModifyStatus)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}

	//如果有权限，修改成员状态
	modifyErr := domain.ModifyOrgMemberStatus(orgId, input.MemberIds, input.Status, userId)
	if modifyErr != nil {
		log.Error(modifyErr)
		return nil, modifyErr
	}
	return &vo.Void{
		ID: orgId,
	}, nil
}

//更新成员检查状态
func UpdateOrgMemberCheckStatus(reqVo orgvo.UpdateOrgMemberCheckStatusReq) (*vo.Void, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	userId := reqVo.UserId
	input := reqVo.Input

	//校验当前用户是否具有修改成员状态的权限
	authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgUser, consts.OperationOrgUserModifyStatus)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}

	//如果有权限，修改成员状态
	modifyErr := domain.ModifyOrgMemberCheckStatus(orgId, input.MemberIds, input.CheckStatus, userId, false)
	if modifyErr != nil {
		log.Error(modifyErr)
		return nil, modifyErr
	}
	return &vo.Void{
		ID: orgId,
	}, nil
}

//func RemoveOrgMember(reqVo orgvo.RemoveOrgMemberReq) (*vo.Void, errs.SystemErrorInfo) {
//	orgId := reqVo.OrgId
//	userId := reqVo.UserId
//	input := reqVo.Input
//
//	//校验当前用户是否具有修改删除成员的权限
//	authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgUser, consts.OperationOrgUserUnbind)
//	if authErr != nil {
//		log.Error(authErr)
//		return nil, authErr
//	}
//	//如果有权限，移除成员
//	modifyErr := domain.RemoveOrgMember(orgId, input.MemberIds, userId)
//	if modifyErr != nil {
//		log.Error(modifyErr)
//		return nil, modifyErr
//	}
//	return &vo.Void{
//		ID: orgId,
//	}, nil
//}

func OrgUserList(orgId, userId int64, page, size int, input *vo.OrgUserListReq) (*vo.UserOrganizationList, errs.SystemErrorInfo) {
	//校验当前用户是否具有查看成员的权限
	//authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgUser, consts.RoleOperationView)
	//if authErr != nil {
	//	log.Error(authErr)
	//	return nil, authErr
	//}
	roleUserResp := make([]rolevo.RoleUser, 0)
	//查询成员角色
	//roleUserResp, err := service.GetOrgRoleUser(orgId, 0)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//所有超管用户
	var allUserHaveRoleIds []int64
	for _, v := range roleUserResp {
		if v.RoleLangCode == consts.RoleGroupOrgAdmin {
			allUserHaveRoleIds = append(allUserHaveRoleIds, v.UserId)
		}
	}

	total, info, err := domain.GetOrganizationUserList(orgId, page, size, input, allUserHaveRoleIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//if len(info) == 0 {
	//	return &vo.UserOrganizationList{
	//		Total: int64(total),
	//	}, nil
	//}
	var userIds []int64
	for _, v := range info {
		userIds = append(userIds, v.UserId, v.AuditorId)
	}

	infoVo := &[]*vo.OrganizationUser{}
	copyErr := copyer.Copy(info, infoVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	userIdsInfo, err := domain.BatchGetUserDetailInfo(userIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	userInfoVos := &[]vo.PersonalInfo{}
	copyErr = copyer.Copy(userIdsInfo, userInfoVos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	userInfoMap := maps.NewMap("ID", *userInfoVos)

	roleGroup := make([]*vo.Role, 0)
	//roleGroup, err := service.GetOrgRoleList(orgId)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	lang := lang2.GetLang()
	isOtherLang := lang2.IsEnglish()
	roleNameLanguageMap := make(map[string]string, 0)
	if tmpMap, ok1 := consts.LANG_ROLE_NAME_MAP[lang]; ok1 {
		roleNameLanguageMap = tmpMap
	}
	orgMemberRoleInfo := vo.Role{
		ID:       -1,
		OrgID:    orgId,
		LangCode: "defaultRole",
		Name:     "系统默认成员",
	}
	for _, v := range roleGroup {
		if v.LangCode == consts.RoleGroupSpecialMember {
			// 多语言适配
			if isOtherLang {
				if tmpVal, ok2 := roleNameLanguageMap[v.Name]; ok2 {
					(*v).Name = tmpVal
				}
			}
			orgMemberRoleInfo = *v
		}
	}

	//暂时默认一个用户最多只有一个组织角色)
	//userRoleMap := maps.NewMap("UserId", roleUserResp.NewData)
	// todo 切换到无码的权限系统，先从无码获取所有管理组，再获取所有权限组，做聚合、map，从而得到哪些人，在哪些组中。
	userRoleMap := map[int64]rolevo.RoleUser{}
	for _, datum := range roleUserResp {
		if _, ok := userRoleMap[datum.UserId]; !ok {
			if isOtherLang {
				if tmpVal, ok2 := roleNameLanguageMap[datum.RoleName]; ok2 {
					datum.RoleName = tmpVal
				}
			}
			userRoleMap[datum.UserId] = datum
		}
	}
	for k, v := range *infoVo {
		if _, ok := userRoleMap[v.UserID]; ok {
			role := userRoleMap[v.UserID]
			(*infoVo)[k].UserRole = &vo.UserRoleInfo{
				ID:       role.RoleId,
				Name:     role.RoleName,
				LangCode: role.RoleLangCode,
			}
		} else {
			(*infoVo)[k].UserRole = &vo.UserRoleInfo{
				ID:       orgMemberRoleInfo.ID,
				Name:     orgMemberRoleInfo.Name,
				LangCode: orgMemberRoleInfo.LangCode,
			}
		}
		if _, ok := userInfoMap[v.UserID]; ok {
			user := userInfoMap[v.UserID].(vo.PersonalInfo)
			(*infoVo)[k].UserInfo = &user
			if (*infoVo)[k].UserRole.LangCode == consts.RoleGroupOrgAdmin {
				(*infoVo)[k].UserInfo.IsAdmin = true
			}
			if (*infoVo)[k].UserRole.LangCode == consts.RoleGroupOrgManager {
				(*infoVo)[k].UserInfo.IsManager = true
			}
		}
		if _, ok := userInfoMap[v.AuditorID]; ok {
			user := userInfoMap[v.AuditorID].(vo.PersonalInfo)
			(*infoVo)[k].AuditorInfo = &user
		}
		if v.CheckStatus == consts.AppCheckStatusSuccess && v.AuditTime.String() <= consts.BlankElasticityTime {
			(*infoVo)[k].AuditTime = v.CreateTime
		}
	}
	return &vo.UserOrganizationList{
		Total: int64(total),
		List:  *infoVo,
	}, nil
}

func GetOrgUserInfoListBySourceChannel(reqVo orgvo.GetOrgUserInfoListBySourceChannelReq) (*orgvo.GetOrgUserInfoListBySourceChannelRespData, errs.SystemErrorInfo) {
	userInfoList, total, err := domain.GetOrgUserInfoListBySourceChannel(reqVo.OrgId, reqVo.SourceChannel, reqVo.Page, reqVo.Size)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &orgvo.GetOrgUserInfoListBySourceChannelRespData{
		Total: total,
		List:  userInfoList,
	}, nil
}

func BatchGetUserDetailInfo(userIds []int64) ([]vo.PersonalInfo, errs.SystemErrorInfo) {
	res, err := domain.BatchGetUserDetailInfoWithMobile(userIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	vos := &[]vo.PersonalInfo{}
	copyErr := copyer.Copy(res, vos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return *vos, nil
}

func ExportAddressList(orgId, userId int64, input orgvo.ExportAddressListReq) (string, errs.SystemErrorInfo) {
	if len(input.ExportField) == 0 {
		return "", errs.ExportFieldIsNull
	}

	info, err := UserList(orgId, orgvo.UserListReq{
		SearchCode:   input.SearchCode,
		IsAllocate:   input.IsAllocate,
		Status:       input.Status,
		RoleId:       input.RoleId,
		DepartmentId: input.DepartmentId,
		Page:         0,
		Size:         0,
	})
	if err != nil {
		log.Error(err)
		return "", err
	}

	relatePath := "/user" + "/org_" + strconv.FormatInt(orgId, 10)
	excelDir := config.GetOSSConfig().RootPath + relatePath
	mkdirErr := os.MkdirAll(excelDir, 0777)
	if mkdirErr != nil {
		log.Error(mkdirErr)
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, mkdirErr)
	}
	fileName := "通讯录.xlsx"
	excelPath := excelDir + "/" + fileName
	url := config.GetOSSConfig().LocalDomain + relatePath + "/" + fileName

	var file *xlsx.File
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, err1 := file.AddSheet("Sheet1")
	if err1 != nil {
		log.Error(err1)
		return "", errs.BuildSystemErrorInfo(errs.SystemError, err1)
	}

	row = sheet.AddRow()

	if util.FieldInUpdate(input.ExportField, "name") {
		cell = row.AddCell()
		cell.Value = "成员"
	}
	if util.FieldInUpdate(input.ExportField, "mobile") {
		cell = row.AddCell()
		cell.Value = "手机"
	}
	if util.FieldInUpdate(input.ExportField, "email") {
		cell = row.AddCell()
		cell.Value = "邮箱"
	}
	if util.FieldInUpdate(input.ExportField, "department") {
		cell = row.AddCell()
		cell.Value = "部门"
	}
	if util.FieldInUpdate(input.ExportField, "isLeader") {
		cell = row.AddCell()
		cell.Value = "部门等级"
	}
	if util.FieldInUpdate(input.ExportField, "role") {
		cell = row.AddCell()
		cell.Value = "角色"
	}
	if util.FieldInUpdate(input.ExportField, "statusChangeTime") {
		cell = row.AddCell()
		cell.Value = "禁用时间"
	}
	if util.FieldInUpdate(input.ExportField, "createTime") {
		cell = row.AddCell()
		cell.Value = "创建时间"
	}

	for _, userInfo := range info.List {
		row = sheet.AddRow()

		if util.FieldInUpdate(input.ExportField, "name") {
			cell = row.AddCell()
			cell.Value = userInfo.Name
		}
		if util.FieldInUpdate(input.ExportField, "mobile") {
			cell = row.AddCell()
			cell.Value = userInfo.PhoneNumber
		}
		if util.FieldInUpdate(input.ExportField, "email") {
			cell = row.AddCell()
			cell.Value = userInfo.Email
		}
		if util.FieldInUpdate(input.ExportField, "department") {
			cell = row.AddCell()
			department := ""
			for _, data := range userInfo.DepartmentList {
				department += data.DeparmentName + ","
			}
			if len(department) > 0 {
				department = department[0 : len(department)-1]
			}
			cell.Value = department
		}
		if util.FieldInUpdate(input.ExportField, "isLeader") {
			cell = row.AddCell()
			value := "成员"
			if input.DepartmentId != nil && *input.DepartmentId != 0 {
				for _, data := range userInfo.DepartmentList {
					if data.DepartmentId == *input.DepartmentId && data.IsLeader == 1 {
						value = "部门主管"
					}
				}
			}
			cell.Value = value
		}
		if util.FieldInUpdate(input.ExportField, "role") {
			cell = row.AddCell()
			role := ""
			for _, data := range userInfo.RoleList {
				role += data.RoleName + ","
			}
			if len(role) > 0 {
				role = role[0 : len(role)-1]
			}
			cell.Value = role
		}
		if util.FieldInUpdate(input.ExportField, "statusChangeTime") {
			cell = row.AddCell()
			cell.Value = userInfo.StatusChangeTime.Format(consts.AppTimeFormat)
		}
		if util.FieldInUpdate(input.ExportField, "createTime") {
			cell = row.AddCell()
			cell.Value = userInfo.CreateTime.Format(consts.AppTimeFormat)
		}
	}

	saveErr := file.Save(excelPath)
	if saveErr != nil {
		log.Error(saveErr)
		return "", errs.SystemError
	}

	return url, nil
}

func UserList(orgId int64, req orgvo.UserListReq) (*orgvo.UserListResp, errs.SystemErrorInfo) {
	union := &db.Union{}
	if req.SearchCode != nil && *req.SearchCode != "" {
		union = union.Or(db.Cond{
			//用户
			consts.TcUserId: db.In(db.Raw("select id from ppm_org_user where name like ?", *req.SearchCode+"%")),
		}).Or(db.Cond{
			//部门
			consts.TcUserId: db.In(db.Raw("select ud.user_id from ppm_org_department d, ppm_org_user_department ud where d.id = ud.department_id and d.org_id = ? and ud.org_id = ? and d.name like ? and d.is_delete = 2 and ud.is_delete = 2", orgId, orgId, *req.SearchCode+"%")),
		}).Or(db.Cond{
			//角色
			consts.TcUserId: db.In(db.Raw("select ud.user_id from ppm_rol_role d, ppm_rol_role_user ud where d.id = ud.role_id and d.org_id = ? and ud.org_id = ?  and d.name like ? and d.is_delete = 2 and ud.is_delete = 2", orgId, orgId, *req.SearchCode+"%")),
		})
	}

	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if req.IsAllocate != nil {
		raw := db.Raw("select user_id from ppm_org_user_department where org_id = ? and is_delete = 2", orgId)
		if *req.IsAllocate == 2 {
			cond[consts.TcUserId] = db.NotIn(raw)
		} else {
			cond[consts.TcUserId] = db.In(raw)
		}
	}

	if req.RoleId != nil {
		cond[consts.TcUserId+"  "] = db.In(db.Raw("select user_id from ppm_rol_role_user where org_id = ? and is_delete = 2 and role_id = ?", orgId, *req.RoleId))
	}

	if req.Status != nil {
		cond[consts.TcStatus] = *req.Status
	}

	if req.DepartmentId != nil && *req.DepartmentId != 0 {
		allDepartmentIds, err := domain.GetChildrenDepartmentIds(orgId, []int64{*req.DepartmentId})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		cond[consts.TcUserId+" "] = db.In(db.Raw("select user_id from ppm_org_user_department where org_id = ? and is_delete = 2 and department_id in ?", orgId, allDepartmentIds))
	}
	pos := &[]po.PpmOrgUserOrganization{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOrganization, cond, union, req.Page, req.Size, "", pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	userIds := []int64{}
	for _, organization := range *pos {
		userIds = append(userIds, organization.UserId)
	}

	//查询用户信息
	userInfo, userInfoErr := domain.BatchGetUserDetailInfoWithMobile(userIds)
	if userInfoErr != nil {
		log.Error(userInfoErr)
		return nil, userInfoErr
	}
	userInfoMap := maps.NewMap("ID", userInfo)
	//查询部门信息
	userDepartmentInfo, userDepartmentInfoErr := domain.GetUserDepartmentInfo(orgId, userIds)
	if userDepartmentInfoErr != nil {
		log.Error(userDepartmentInfoErr)
		return nil, userDepartmentInfoErr
	}
	userDepartmentMap := map[int64][]orgvo.UserDepartmentData{}
	for _, info := range userDepartmentInfo {
		userDepartmentMap[info.UserId] = append(userDepartmentMap[info.UserId], orgvo.UserDepartmentData{
			DepartmentId:  info.DepartmentId,
			IsLeader:      info.IsLeader,
			DeparmentName: info.DepartmentName,
		})
	}

	res := &orgvo.UserListResp{
		Total: int64(total),
		List:  []*orgvo.UserInfo{},
	}

	//查询组织创建者
	orgInfo := &po.PpmOrgOrganization{}
	orgErr := mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
		consts.TcId: orgId,
	}, orgInfo)
	if orgErr != nil {
		if orgErr == db.ErrNoMoreRows {
			return nil, errs.OrgNotExist
		} else {
			log.Error(orgErr)
			return nil, errs.MysqlOperateError
		}
	}
	for _, organization := range *pos {
		if _, ok := userInfoMap[organization.UserId]; !ok {
			continue
		}
		userInfo := userInfoMap[organization.UserId].(bo.UserInfoBo)
		temp := &orgvo.UserInfo{
			UserID:           userInfo.ID,
			Name:             userInfo.Name,
			NamePy:           userInfo.NamePinyin,
			Avatar:           userInfo.Avatar,
			Email:            userInfo.Email,
			PhoneNumber:      userInfo.Mobile,
			CreateTime:       time.Time(userInfo.CreateTime),
			StatusChangeTime: organization.StatusChangeTime,
			Status:           organization.Status,
			IsCreator:        false,
		}
		if _, ok := userDepartmentMap[organization.UserId]; ok {
			temp.DepartmentList = userDepartmentMap[organization.UserId]
		} else {
			temp.DepartmentList = []orgvo.UserDepartmentData{}
		}

		if organization.UserId == orgInfo.Creator {
			temp.IsCreator = true
		}

		res.List = append(res.List, temp)
	}

	return res, nil
}

// 查询在某个时间段内新初始化的用户信息
func GetUserListWithCreateTimeRange(input *orgvo.GetUserListWithCreateTimeRangeReq) (*orgvo.GetUserListWithCreateTimeRangeResp, errs.SystemErrorInfo) {
	pos := &[]po.PpmOrgUserOutInfo{}
	cond1 := db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcSourceChannel: input.SourceChannel,
		consts.TcCreateTime:    db.Between(input.CreateTime1, input.CreateTime2),
	}
	_, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOutInfo, cond1, nil, input.Page, input.Size, nil, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	log.Info(json.ToJsonIgnoreError(pos))
	// 通过 orgId 查询 outOrgId
	orgIds := make([]int64, 0)
	for _, item := range *pos {
		orgIds = append(orgIds, item.OrgId)
	}
	// 去重
	orgIds = int642.ArrayUnique(orgIds)
	if len(orgIds) < 1 {
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.MysqlOperateError, "查询到的用户的企业id为空。")
	}
	var outOrgs = make([]*po.PpmOrgOrganizationOutInfo, 0)
	cond1 = db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcSourceChannel: input.SourceChannel,
		consts.TcOrgId:         db.In(orgIds),
	}
	err = mysql.SelectAllByCond(consts.TableOrganizationOutInfo, cond1, &outOrgs)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	outOrgCodeMap := make(map[int64]string, 0)
	for _, outOrg := range outOrgs {
		outOrgCodeMap[outOrg.OrgId] = outOrg.OutOrgId
	}
	// 组装要返回的用户列表
	userList := make([]*orgvo.GetUserListWithCreateTimeRangeItem, 0)
	for _, item := range *pos {
		if outOrgId, ok := outOrgCodeMap[item.OrgId]; ok {
			one := &orgvo.GetUserListWithCreateTimeRangeItem{
				UserId:     item.UserId,
				Name:       item.Name,
				OpenId:     item.OutUserId,
				OrgId:      item.OrgId,
				OutOrgId:   outOrgId,
				CreateTime: item.CreateTime.Format(consts.AppTimeFormat),
			}
			userList = append(userList, one)
		}
	}

	return &orgvo.GetUserListWithCreateTimeRangeResp{
		userList,
	}, nil
}

func EmptyUser(orgId, userId int64, input orgvo.EmptyUserReq) (*vo.Void, errs.SystemErrorInfo) {
	//获取对应的用户
	pos := &[]po.PpmOrgUserOrganization{}
	err := mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcStatus:      input.Status,
		consts.TcOrgId:       orgId,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	if len(*pos) == 0 {
		return &vo.Void{
			ID: orgId,
		}, nil
	}
	memberIds := []int64{}
	for _, organization := range *pos {
		memberIds = append(memberIds, organization.UserId)
	}
	modifyErr := domain.RemoveOrgMember(orgId, memberIds, userId)
	if modifyErr != nil {
		log.Error(modifyErr)
		return nil, modifyErr
	}
	return &vo.Void{
		ID: orgId,
	}, nil
}

//func CreateUser(orgId, userId int64, input orgvo.CreateUserReq) (*vo.Void, errs.SystemErrorInfo) {
//	name := strings.TrimSpace(input.Name)
//	//检测姓名是否合法
//	isNameRight := format.VerifyUserNameFormat(name)
//	if !isNameRight {
//		return nil, errs.UserNameLenError
//	}
//
//	//注册
//	userBo, err := domain.CreateUser(orgId, userId, input)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	return &vo.Void{ID: userBo.ID}, nil
//}

//func UpdateUser(orgId, userId int64, regInfo orgvo.UpdateUserReq) (*vo.Void, errs.SystemErrorInfo) {
//	//查看用户信息
//	info, _, err := domain.GetUserBo(regInfo.UserId)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	upd := mysql.Upd{}
//	if regInfo.Name != nil {
//		name := strings.TrimSpace(*regInfo.Name)
//		//检测姓名是否合法
//		isNameRight := format.VerifyUserNameFormat(name)
//		if !isNameRight {
//			return nil, errs.UserNameLenError
//		}
//
//		upd[consts.TcName] = *regInfo.Name
//	}
//
//	if regInfo.Email != nil && *regInfo.Email != info.Email {
//		if !format.VerifyEmailFormat(*regInfo.Email) {
//			return nil, errs.EmailFormatErr
//		}
//		exist, err := mysql.IsExistByCond(consts.TableUser, db.Cond{
//			consts.TcOrgId:    orgId,
//			consts.TcId:       db.NotEq(regInfo.UserId),
//			consts.TcIsDelete: consts.AppIsNoDelete,
//			consts.TcEmail:    *regInfo.Email,
//		})
//		if err != nil {
//			log.Error(err)
//			return nil, errs.MysqlOperateError
//		}
//
//		if exist {
//			return nil, errs.EmailAlreadyBindByOtherAccountError
//		}
//		upd[consts.TcEmail] = *regInfo.Email
//	}
//
//	//if regInfo.PhoneNumber != nil && *regInfo.PhoneNumber != info.Mobile {
//	//	exist, err := mysql.IsExistByCond(consts.TableUser, db.Cond{
//	//		consts.TcOrgId:    orgId,
//	//		consts.TcId:       db.NotEq(regInfo.UserId),
//	//		consts.TcIsDelete: consts.AppIsNoDelete,
//	//		consts.TcMobile:   *regInfo.PhoneNumber,
//	//	})
//	//	if err != nil {
//	//		log.Error(err)
//	//		return nil, errs.MysqlOperateError
//	//	}
//	//
//	//	if exist {
//	//		return nil, errs.MobileAlreadyBindOtherAccountError
//	//	}
//	//	upd[consts.TcMobile] = *regInfo.PhoneNumber
//	//}
//
//	if len(upd) > 0 {
//		upd[consts.TcUpdator] = userId
//		_, err := mysql.UpdateSmartWithCond(consts.TableUser, db.Cond{
//			consts.TcId: regInfo.UserId,
//		}, upd)
//		if err != nil {
//			log.Error(err)
//			return nil, errs.MysqlOperateError
//		}
//	}
//
//	if regInfo.Status != nil {
//		if ok, _ := slice.Contain([]int{consts.AppStatusEnable, consts.AppStatusDisabled}, *regInfo.Status); ok {
//			_, err := mysql.UpdateSmartWithCond(consts.TableUserOrganization, db.Cond{
//				consts.TcOrgId:  orgId,
//				consts.TcUserId: regInfo.UserId,
//			}, mysql.Upd{
//				consts.TcStatus:  *regInfo.Status,
//				consts.TcUpdator: userId,
//			})
//			if err != nil {
//				log.Error(err)
//				return nil, errs.MysqlOperateError
//			}
//		}
//	}
//
//	if regInfo.DepartmentIds != nil {
//		err := AllocateDepartment(orgvo.AllocateDepartmentReq{
//			UserIds:       []int64{regInfo.UserId},
//			DepartmentIds: *regInfo.DepartmentIds,
//		}, orgId, userId)
//		if err != nil {
//			log.Error(err)
//			return nil, err
//		}
//	}
//
//	// 清除user缓存
//	err = domain.ClearBaseUserInfo(orgId, userId)
//	if err != nil {
//		log.Error(err)
//	}
//
//	return &vo.Void{ID: regInfo.UserId}, nil
//}

func SearchUser(orgId int64, req orgvo.SearchUserReq) (*orgvo.SearchUserResp, errs.SystemErrorInfo) {
	email := req.Email
	if email == "" {
		return nil, errs.EmailFormatErr
	}
	//查看邮箱是否已注册
	info := &po.PpmOrgUser{}
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	selectErr := conn.Select(db.Raw("u.*")).From("ppm_org_user u", "ppm_org_user_organization o").Where(db.Cond{
		"o." + consts.TcOrgId:    orgId,
		"u." + consts.TcEmail:    email,
		"o." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"o." + consts.TcUserId:   db.Raw("u." + consts.TcId),
	}).One(info)
	if selectErr != nil {
		if selectErr == db.ErrNoMoreRows {
			//如果没有的话就去判断是否邀请过
			isInvited, err := mysql.IsExistByCond(consts.TableUserInvite, db.Cond{
				consts.TcIsDelete:       consts.AppIsNoDelete,
				consts.TcEmail:          email,
				consts.TcIsRegister:     2,
				consts.TcOrgId:          orgId,
				consts.TcLastInviteTime: db.Gte(date.Format(time.Now().AddDate(0, 0, -1))),
			})
			if err != nil {
				log.Error(err)
				return nil, errs.MysqlOperateError
			}
			if !isInvited {
				return &orgvo.SearchUserResp{
					Status:   1, //可邀请
					UserInfo: nil,
				}, nil
			} else {
				return &orgvo.SearchUserResp{
					Status:   2, //已邀请
					UserInfo: nil,
				}, nil
			}
		} else {
			log.Error(selectErr)
			return nil, errs.MysqlOperateError
		}
	}

	return &orgvo.SearchUserResp{
		Status: 3, //已注册
		UserInfo: &orgvo.UserInfo{
			UserID: info.Id,
			Name:   info.Name,
			NamePy: info.NamePinyin,
			Avatar: info.Avatar,
		},
	}, nil
}

//已废弃
func InviteUser(orgId, userId int64, param orgvo.InviteUserReq) (*orgvo.InviteUserResp, errs.SystemErrorInfo) {
	result := &orgvo.InviteUserResp{}
	if len(param.Data) == 0 {
		return result, nil
	}

	var validEmail []string
	for _, datum := range param.Data {
		if !format.VerifyEmailFormat(datum.Email) {
			result.InvalidEmail = append(result.InvalidEmail, datum.Email)
		} else {
			validEmail = append(validEmail, datum.Email)
		}
	}
	if len(validEmail) == 0 {
		return result, nil
	}

	validEmail = slice.SliceUniqueString(validEmail)
	emailNameMap := map[string]string{}
	for _, datum := range param.Data {
		emailNameMap[datum.Email] = datum.Name
	}

	//查看是否是用户
	userInfo := &[]po.PpmOrgUser{}
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	selectErr := conn.Select(db.Raw("u.*")).From("ppm_org_user u", "ppm_org_user_organization o").Where(db.Cond{
		"o." + consts.TcOrgId:    orgId,
		"u." + consts.TcEmail:    db.In(validEmail),
		"o." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"o." + consts.TcUserId:   db.Raw("u." + consts.TcId),
	}).All(userInfo)
	if selectErr != nil {
		log.Error(selectErr)
		return nil, errs.MysqlOperateError
	}

	if len(*userInfo) > 0 {
		for _, user := range *userInfo {
			result.IsUserEmail = append(result.IsUserEmail, user.Email)
		}
	}

	var notUserEmail []string
	for _, s := range validEmail {
		if ok, _ := slice.Contain(result.IsUserEmail, s); !ok {
			notUserEmail = append(notUserEmail, s)
		}
	}
	if len(notUserEmail) == 0 {
		return result, nil
	}

	//查看是否已邀请
	inviteInfo := &[]po.PpmOrgUserInvite{}
	inviteErr := mysql.SelectAllByCond(consts.TableUserInvite, db.Cond{
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcEmail:      db.In(notUserEmail),
		consts.TcIsRegister: 2,
		consts.TcOrgId:      orgId,
	}, inviteInfo)
	if inviteErr != nil {
		log.Error(inviteErr)
		return nil, errs.MysqlOperateError
	}

	needInviteAgain := []int64{}
	if len(*inviteInfo) > 0 {
		for _, invite := range *inviteInfo {
			if invite.LastInviteTime.Before(time.Now().AddDate(0, 0, -1)) {
				//已邀请（需要再次邀请,更新数据库）
				needInviteAgain = append(needInviteAgain, invite.Id)
			} else {
				result.InvitedEmail = append(result.InvitedEmail, invite.Email)
			}
		}
	}
	for _, s := range notUserEmail {
		if ok, _ := slice.Contain(result.InvitedEmail, s); !ok {
			result.SuccessEmail = append(result.SuccessEmail, s)
		}
	}

	if len(result.SuccessEmail) == 0 {
		return result, nil
	}
	codeResp, codeErr := GetInviteCode(userId, orgId, "")
	if codeErr != nil {
		log.Error(codeErr)
		return nil, codeErr
	}

	mailErr := msgfacade.SendMailRelaxed(result.SuccessEmail, consts.MailTemplateSubjectInvite, temp.RenderIgnoreError(consts.MailTemplateContentInvite, map[string]string{
		consts.SMSParamsNameInviteHref: config.GetServerConfig().Domain + "/user/entry?inviteCode=" + codeResp.InviteCode,
		consts.SMSParamsNameInviteUrl:  config.GetServerConfig().Domain + "/user/entry?inviteCode=" + codeResp.InviteCode,
	}))
	if mailErr != nil {
		log.Error(mailErr)
		return nil, mailErr
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		now := time.Now()
		if len(needInviteAgain) > 0 {
			_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserInvite, db.Cond{
				consts.TcOrgId: orgId,
				consts.TcId:    db.In(needInviteAgain),
			}, mysql.Upd{
				consts.TcUpdator:        userId,
				consts.TcLastInviteTime: now,
			})
			if err != nil {
				log.Error(err)
				return err
			}
		}

		pos := []interface{}{}
		for _, s := range result.SuccessEmail {
			var name string
			if _, ok := emailNameMap[s]; ok {
				name = emailNameMap[s]
			}
			pos = append(pos, po.PpmOrgUserInvite{
				Id:             snowflake.Id(),
				OrgId:          orgId,
				Name:           name,
				Email:          s,
				InviteUserId:   userId,
				LastInviteTime: now,
				Creator:        userId,
				Updator:        userId,
			})
		}

		err := mysql.TransBatchInsert(tx, &po.PpmOrgUserInvite{}, pos)
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}
	return result, nil
}

func InviteUserList(orgId int64, listReq orgvo.InviteUserListReq) (*orgvo.InviteUserListResp, errs.SystemErrorInfo) {
	pos := &[]po.PpmOrgUserInvite{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserInvite, db.Cond{
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcIsRegister: 2,
		consts.TcOrgId:      orgId,
	}, nil, listReq.Page, listReq.Size, "create_time desc", pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	var resInfo []orgvo.InviteUserInfo
	for _, invite := range *pos {
		resInfo = append(resInfo, orgvo.InviteUserInfo{
			Id:              invite.Id,
			Name:            invite.Name,
			Email:           invite.Email,
			InviteTime:      invite.LastInviteTime,
			IsInvitedRecent: invite.LastInviteTime.Before(time.Now().AddDate(0, 0, -1)),
		})
	}

	return &orgvo.InviteUserListResp{
		Total: int64(total),
		List:  resInfo,
	}, nil
}

func RemoveInviteUser(orgId, userId int64, param orgvo.RemoveInviteUserReq) errs.SystemErrorInfo {
	if len(param.Ids) == 0 && param.IsAll != 1 {
		return nil
	}
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		cond := db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}
		if param.IsAll != 1 {
			cond[consts.TcId] = db.In(param.Ids)
		}
		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserInvite, cond, mysql.Upd{
			consts.TcUpdator:  userId,
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}

	return nil
}

func UserStat(orgId int64) (*orgvo.UserStatResp, errs.SystemErrorInfo) {
	allCount, err := mysql.SelectCountByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcOrgId:       orgId,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		consts.TcStatus:      consts.AppStatusEnable,
	})
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	forbiddenCount, err := mysql.SelectCountByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcOrgId:       orgId,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		consts.TcStatus:      consts.AppStatusDisabled,
	})
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	unreceivedCount, err := mysql.SelectCountByCond(consts.TableUserInvite, db.Cond{
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcIsRegister: 2,
		consts.TcOrgId:      orgId,
	})
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	unallocatedCount, err := mysql.SelectCountByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcOrgId:       orgId,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		consts.TcStatus:      consts.AppStatusEnable,
		consts.TcUserId:      db.NotIn(db.Raw("select user_id from ppm_org_user_department where org_id = ? and is_delete = 2", orgId)),
	})
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	return &orgvo.UserStatResp{
		AllCount:         int64(allCount),
		UnallocatedCount: int64(unallocatedCount),
		UnreceivedCount:  int64(unreceivedCount),
		ForbiddenCount:   int64(forbiddenCount),
	}, nil
}
