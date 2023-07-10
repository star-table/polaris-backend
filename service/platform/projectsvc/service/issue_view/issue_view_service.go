package issue_view

import (
	"html"

	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
)

var (
	log = *logger.GetDefaultLogger()
)

// 创建视图
// 项目参与人可以创建私有视图
// 公有视图只能项目负责人才能创建。
func CreateTaskView(orgId, userId int64, input *vo.CreateIssueViewReq) (*vo.GetIssueViewListItem, errs.SystemErrorInfo) {
	if input.IsPrivate == nil {
		defaultVal := true
		input.IsPrivate = &defaultVal
	}
	// 有视图权限/管理员，才能创建/编辑视图
	if input.ProjectID != nil && *input.ProjectID > 0 {
		if err := CheckIssueViewAuthForCreate(orgId, userId, *input.ProjectID, *input.IsPrivate); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	poObj, err := domain.CreateIssueView(orgId, userId, input)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	newObj := vo.GetIssueViewListItem{}
	copyer.Copy(poObj, &newObj)

	return &newObj, nil
}

// 获取我的视图列表。在这里封装一下默认条件。
func GetTaskViewList(orgId, userId int64, input *vo.GetIssueViewListReq) (*vo.GetIssueViewListResp, errs.SystemErrorInfo) {
	boList := make([]*vo.GetIssueViewListItem, 0)
	result := &vo.GetIssueViewListResp{
		Total: 0,
		List:  boList,
	}
	// 非私有表示公共视图+私有视图
	isPrivate := false
	input.IsPrivate = &isPrivate
	total, resp, err := domain.GetTaskViewList(orgId, userId, input)
	if err != nil {
		log.Error(err)
		return result, err
	}
	if err := copyer.Copy(resp, &boList); err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}
	for index, item := range boList {
		(*item).IsPrivate = resp[index].Owner > 0
	}
	result.Total = total
	result.List = boList
	return result, nil
}

// 更新任务视图
// 公共视图：项目负责人、管理员才能更新；私有视图：自己可以更新。
func UpdateTaskView(orgId, userId int64, input *vo.UpdateIssueViewReq) (*vo.Void, errs.SystemErrorInfo) {
	view, err := domain.GetOneTaskView(orgId, input.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	input = CheckInputForUpdateTaskView(input)
	err = CheckIssueViewAuthForUpdate(orgId, userId, view.ProjectID, view, input)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	updateData := assemblyUpdForUpdateTaskView(input, userId)
	_, oriErr := dao.UpdateIssueViewById(orgId, input.ID, updateData)
	if oriErr != nil {
		log.Error(oriErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	return &vo.Void{
		ID: input.ID,
	}, nil
}

func DeleteTaskView(orgId, userId, viewId int64) (*vo.Void, errs.SystemErrorInfo) {
	view, err := domain.GetOneTaskView(orgId, viewId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = CheckIssueViewAuthForUpdate(orgId, userId, view.ProjectID, view, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	affectedNum, oriErr := dao.UpdateIssueViewById(orgId, viewId, mysql.Upd{
		consts.TcDelFlag: consts.AppIsDeleted,
	})
	if oriErr != nil {
		log.Error(oriErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	return &vo.Void{
		ID: affectedNum,
	}, nil
}

// 组装要更新的数据。
func assemblyUpdForUpdateTaskView(input *vo.UpdateIssueViewReq, currentUserId int64) mysql.Upd {
	updateData := mysql.Upd{}
	if input.Config != nil {
		updateData[consts.TcConfig] = *input.Config
	}
	if input.Remark != nil {
		updateData[consts.TcRemark] = *input.Remark
	}
	if input.IsPrivate != nil {
		if *input.IsPrivate {
			updateData[consts.TcOwner] = currentUserId
		} else {
			updateData[consts.TcOwner] = int64(0)
		}
	}
	if input.ViewName != nil {
		updateData[consts.TcIssueViewName] = *input.ViewName
	}
	if input.Type != nil {
		updateData[consts.TcType] = *input.Type
	}
	if input.ProjectObjectTypeID != nil {
		updateData[consts.TcProjectObjectTypeId] = *input.ProjectObjectTypeID
	}
	return updateData
}

func CheckInputForUpdateTaskView(input *vo.UpdateIssueViewReq) *vo.UpdateIssueViewReq {
	// 对用户输入进行转义
	if input.ViewName != nil {
		*input.ViewName = html.EscapeString(*input.ViewName)
	}
	if input.Remark != nil {
		*input.Remark = html.EscapeString(*input.Remark)
	}
	return input
}

func CheckIssueViewAuthForCreate(orgId, userId, projectId int64, isPrivate bool) errs.SystemErrorInfo {
	if isPrivate {
		// 检查是否是项目参与人。
		isOk, err := domain.IsProjectParticipant(orgId, userId, projectId)
		if err != nil {
			log.Error(err)
			return err
		}
		// 不是项目成员，不允许创建视图
		if !isOk {
			return errs.IssueViewOpDeny
		}
	} else {
		// 项目负责人有刪除项目、设置项目，才能管理视图
		err := domain.AuthProject(orgId, userId, projectId, consts.RoleOperationPathOrgProProConfig, consts.RoleOperationDelete)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.IssueViewOpDeny, err)
		}
		// 检查是否有视图管理权限
		err = domain.AuthProject(orgId, userId, projectId, consts.RoleOperationPathOrgProProConfig, consts.RoleOperationModify)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.IssueViewOpDeny, err)
		}
	}

	return nil
}

func CheckIssueViewAuthForUpdate(orgId, userId, projectId int64, oldView *po.PpmPriIssueView, input *vo.UpdateIssueViewReq) errs.SystemErrorInfo {
	if oldView.Owner > 0 {
		// 私有视图
		if userId != oldView.Owner {
			return errs.IssueViewOpDeny
		}
	} else {
		// 项目负责人有刪除项目、设置项目，才能管理视图
		err := domain.AuthProject(orgId, userId, projectId, consts.RoleOperationPathOrgProProConfig, consts.RoleOperationDelete)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.IssueViewOpDeny, err)
		}
		// 检查是否有视图管理权限
		err = domain.AuthProject(orgId, userId, projectId, consts.RoleOperationPathOrgProProConfig, consts.RoleOperationModify)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.IssueViewOpDeny, err)
		}
	}

	return nil
}
