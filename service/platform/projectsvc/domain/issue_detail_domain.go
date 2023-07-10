package domain

import (
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

//func GetIssueDetailBo(orgId, issueId int64) (*bo.IssueDetailBo, errs.SystemErrorInfo) {
//	//获取issue详情
//	issueDetail := &po.PpmPriIssueDetail{}
//	err := mysql.SelectOneByCond(issueDetail.TableName(), db.Cond{
//		consts.TcOrgId:   orgId,
//		consts.TcIssueId: issueId,
//		//consts.TcIsDelete: consts.AppIsNoDelete,
//	}, issueDetail)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDetailNotExist)
//	}
//	issueDetailBo := &bo.IssueDetailBo{}
//	err1 := copyer.Copy(issueDetail, issueDetailBo)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
//	}
//	return issueDetailBo, nil
//}

//func UpdateIssueDetailRemark(issueBo bo.IssueBo, operatorId int64, remark string, remarkDetail *string) errs.SystemErrorInfo {
//	upd := mysql.Upd{
//		consts.TcRemark:     remark,
//		consts.TcUpdator:    operatorId,
//		consts.TcUpdateTime: times.GetBeiJingTime(),
//	}
//	if remarkDetail != nil {
//		upd[consts.TcRemarkDetail] = *remarkDetail
//	}
//	//detail
//	issueDetail := &po.PpmPriIssueDetail{}
//	_, err := mysql.UpdateSmartWithCond(issueDetail.TableName(), db.Cond{
//		consts.TcOrgId:    issueBo.OrgId,
//		consts.TcIssueId:  issueBo.Id,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}, upd)
//	if err != nil {
//		log.Errorf("mysql.UpdateSmartWithCond: %c\n", err)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	return nil
//}

//func GetIssueDetailBoBatch(orgId int64, issueIds []int64) ([]bo.IssueDetailBo, errs.SystemErrorInfo) {
//	//获取issue详情
//	pos := &[]po.PpmPriIssueDetail{}
//	err := mysql.SelectAllByCond(consts.TableIssueDetail, db.Cond{
//		consts.TcOrgId:   orgId,
//		consts.TcIssueId: db.In(issueIds),
//		//consts.TcIsDelete: consts.AppIsNoDelete,
//	}, pos)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.MysqlOperateError
//	}
//	issueDetailBo := &[]bo.IssueDetailBo{}
//	err1 := copyer.Copy(pos, issueDetailBo)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
//	}
//	return *issueDetailBo, nil
//}

// GetDeptListForImportIssue 导入任务时，自定义字段（部门类型）处理时，查询组织的部门列表
func GetDeptListForImportIssue(orgId int64) (map[string]orgvo.SimpleDeptForCustomField, errs.SystemErrorInfo) {
	page := 1
	size := 20000
	resp := orgfacade.Departments(orgvo.DepartmentsReqVo{
		Page:   &page,
		Size:   &size,
		Params: nil,
		OrgId:  orgId,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}

	_, multiDeptIds := GetSameNameDeptByList(resp.DepartmentList.List)
	deptMaps := make(map[string]orgvo.SimpleDeptForCustomField, 0)

	for _, dept := range resp.DepartmentList.List {
		tmpDeptObj := orgvo.SimpleDeptForCustomField{
			ID:       dept.ID,
			Name:     dept.Name,
			Avatar:   "",
			IsDelete: consts.AppIsNoDelete, // 能获取到的都是未被删除的
			Status:   dept.Status,
			Type:     consts.LcCustomFieldDeptType,
		}
		deptMaps[dept.Name] = tmpDeptObj
		// 针对同名的部门，再额外组装一次唯一的名称 key（`dept001#1001#`）
		if exist, _ := slice.Contain(multiDeptIds, dept.ID); exist {
			deptMaps[RenderSameNameKey(dept.ID, dept.Name)] = tmpDeptObj
		}
	}
	return deptMaps, nil
}

// GetSameNameDeptByList 通过部门列表取出有同名情况的部门
func GetSameNameDeptByList(deptList []*vo.Department) ([]*vo.Department, []int64) {
	// 找出同名的组织
	multiMap := make(map[string][]*vo.Department, 0)
	multiDeptMapById := make(map[int64]*vo.Department, 0)
	for _, dept := range deptList {
		multiDeptMapById[dept.ID] = dept

		if _, ok := multiMap[dept.Name]; ok {
			multiMap[dept.Name] = append(multiMap[dept.Name], dept)
		} else {
			multiMap[dept.Name] = []*vo.Department{
				dept,
			}
		}
	}
	multiDeptIds := make([]int64, 0)
	multiDeptList := make([]*vo.Department, 0)
	for _, deptArr := range multiMap {
		if len(deptArr) > 1 {
			for _, dept := range deptArr {
				multiDeptIds = append(multiDeptIds, dept.ID)
			}
		}
	}
	for _, deptId := range multiDeptIds {
		if dept, ok := multiDeptMapById[deptId]; ok {
			multiDeptList = append(multiDeptList, dept)
		}
	}

	return multiDeptList, multiDeptIds
}
