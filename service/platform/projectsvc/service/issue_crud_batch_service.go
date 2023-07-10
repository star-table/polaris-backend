package service

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	automationPb "github.com/star-table/interface/golang/automation/v1"

	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/facade/tablefacade"

	"github.com/star-table/polaris-backend/common/core/util/jsonx"

	"github.com/spf13/cast"
	"github.com/star-table/common/core/errors"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/date"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	msgPb "github.com/star-table/interface/golang/msg/v1"
	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/polaris-backend/common/core/businees"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/core/util/batch"
	"github.com/star-table/polaris-backend/common/core/util/convert"
	"github.com/star-table/polaris-backend/common/core/util/format"
	slice2 "github.com/star-table/polaris-backend/common/core/util/slice"
	int64Slice "github.com/star-table/polaris-backend/common/core/util/slice/int64"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/bo/mqbo"
	"github.com/star-table/polaris-backend/common/model/bo/status"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
	"github.com/star-table/polaris-backend/facade/common/report"
	"github.com/star-table/polaris-backend/facade/formfacade"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/permissionfacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	sconsts "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service/work_hour"
	"upper.io/db.v3/lib/sqlbuilder"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//  一些公共方法
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const BATCH_SIZE = 1000

// 一些更新字段参数校验
func checkColumnValue(columnId string, value interface{}) (interface{}, errs.SystemErrorInfo) {
	switch columnId {
	case consts.BasicFieldTitle:
		v, err := cast.ToStringE(value)
		if err != nil {
			return nil, errs.IssueTitleError
		}
		if !format.VerifyIssueNameFormat(v) {
			return nil, errs.IssueTitleError
		}
		return strings.TrimSpace(v), nil
	case consts.BasicFieldRemark:
		v, err := cast.ToStringE(value)
		if err != nil {
			return nil, errs.IssueRemarkLenError
		}
		if !format.VerifyIssueRemarkFormat(v) {
			return nil, errs.IssueRemarkLenError
		}
		return value, nil
	case consts.BasicFieldPlanStartTime, consts.BasicFieldPlanEndTime:
		if cast.ToString(value) == "" {
			return consts.BlankTime, nil
		}
		return value, nil
	default:
		return value, nil
	}
}

func prependInt64(x []int64, y int64) []int64 {
	x = append(x, 0)
	copy(x[1:], x)
	x[0] = y
	return x
}

func prependSliceMapStringInterface(x [][]map[string]interface{}, y []map[string]interface{}) [][]map[string]interface{} {
	x = append(x, nil)
	copy(x[1:], x)
	x[0] = y
	return x
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//  批量创建
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type BatchCreateIssueContext struct {
	TriggerBy        *projectvo.TriggerBy
	AsyncTaskId      string
	Req              *projectvo.BatchCreateIssueReqVo
	OrgBaseInfo      *bo.BaseOrgInfoBo
	Project          *bo.ProjectBo
	TableMeta        *projectvo.TableMetaData
	TableColumns     map[string]*projectvo.TableColumnData         // 表头
	TableColumnsDeep map[string]*lc_table.LcCommonField            // 表头（解析props信息)
	AllIssueStatus   []status.StatusInfoBo                         // IssueStatus表头
	CreateColumnMap  map[string]struct{}                           // 涉及创建的字段
	CreateColumns    []string                                      // 涉及创建的字段
	UserIds          []int64                                       // 涉及的人员Id
	DeptIds          []int64                                       // 涉及的部门Id
	Users            map[int64]*uservo.GetAllUserByIdsRespDataUser // 涉及的用户
	Depts            map[int64]*uservo.GetAllDeptByIdsRespDataUser // 涉及的部门
	UserDepts        map[string]*uservo.MemberDept                 // 组装好的用户、部门
	RelateData       map[string]map[string]interface{}             // 关联/引用的相关数据
	RelatingIssueIds []int64                                       // 关联前后置修改涉及到的任务Id
	RelatingIssues   map[int64]*bo.IssueBo                         // 关联前后置修改涉及到的任务
	AllParentIds     map[int64]struct{}                            // 涉及的父任务id
	AllParentPaths   map[int64]string                              // 涉及的父任务path
	BeforeOrder      *float64                                      // 涉及的beforeId的order
	AfterOrder       *float64                                      // 涉及的afterId的order
	Now              string                                        // 当前时间

	IssueIds             *bo.IdCodes                    // 创建的IssueIds
	IssueCodes           *bo.IdCodes                    // 创建的IssueCodes
	CreateCount          int64                          // 创建任务数
	CreateIndex          int64                          // 创建index
	CreateIssueIds       []int64                        // 需要创建的任务Id
	CreateIssues         []*CreateIssueVo               // 需要创建的任务
	CreateIssuesMap      map[int64]*CreateIssueVo       // 需要创建的任务
	CreateIssueRelations []*po.PpmPriIssueRelation      // 需创建的issue relations
	CreateAttachments    map[int64]*bo.IssueAttachments // 需添加关联的附件
	IssueTrends          []*bo.IssueTrendsBo            // 任务动态结果集
	CreateResults        []map[string]interface{}       // 创建结果
}

type CreateIssueVo struct {
	IssueId              int64
	Data                 map[string]interface{}
	IssueBo              *bo.IssueBo
	UpdateColumns        []string
	RelatingChange       map[string]*bo.RelatingChangeBo
	BaRelatingChange     map[string]*bo.RelatingChangeBo
	SingleRelatingChange map[string]*bo.RelatingChangeBo
}

// SyncBatchCreateIssue 每 BATCH_SIZE 条为一个批次，执行同步批量插入
func SyncBatchCreateIssue(req *projectvo.BatchCreateIssueReqVo, withoutAuth bool, triggerBy *projectvo.TriggerBy) (
	[]map[string]interface{}, map[string]*uservo.MemberDept, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	var count int64
	var d []map[string]interface{}
	var ds [][]map[string]interface{}
	var result []map[string]interface{}
	userDept := make(map[string]*uservo.MemberDept)
	relateData := make(map[string]map[string]interface{})
	var counts []int64
	for _, data := range req.Data {
		d = append(d, data)
		count += 1
		count += cast.ToInt64(data[consts.TempFieldChildrenCount])
		if count >= BATCH_SIZE {
			// 倒序创建，以保证最终生成的任务在无码中的顺序
			counts = prependInt64(counts, count)
			ds = prependSliceMapStringInterface(ds, d)

			// 重新累计
			d = make([]map[string]interface{}, 0)
			count = 0
		}
	}
	// 最后一批率先创建
	if count > 0 {
		counts = prependInt64(counts, count)
		ds = prependSliceMapStringInterface(ds, d)
	}

	for i := 0; i < len(counts); i++ {
		r := &projectvo.BatchCreateIssueReqVo{
			OrgId:         req.OrgId,
			UserId:        req.UserId,
			AppId:         req.AppId,
			ProjectId:     req.ProjectId,
			TableId:       req.TableId,
			BeforeDataId:  req.BeforeDataId,
			AfterDataId:   req.AfterDataId,
			Data:          ds[i],
			IsIdGenerated: req.IsIdGenerated,
		}
		res, ud, rd, err := BatchCreateIssue(r, withoutAuth, triggerBy, "", 0)
		if err != nil {
			return nil, nil, nil, err
		}
		result = append(result, res...)
		if ud != nil {
			for k, v := range ud {
				userDept[k] = v
			}
		}
		if rd != nil {
			for k, v := range rd {
				relateData[k] = v
			}
		}
	}
	return result, userDept, relateData, nil
}

// AsyncBatchCreateIssue 每 BATCH_SIZE 条为一个批次，执行异步批量插入
func AsyncBatchCreateIssue(req *projectvo.BatchCreateIssueReqVo, withoutAuth bool, triggerBy *projectvo.TriggerBy, asyncTaskId string) {
	var count int64
	var d []map[string]interface{}
	var ds [][]map[string]interface{}
	var counts []int64
	for _, data := range req.Data {
		d = append(d, data)
		count += 1
		count += cast.ToInt64(data[consts.TempFieldChildrenCount])
		if count >= BATCH_SIZE {
			// 倒序创建，以保证最终生成的任务在无码中的顺序
			counts = prependInt64(counts, count)
			ds = prependSliceMapStringInterface(ds, d)

			// 重新累计
			count = 0
			d = make([]map[string]interface{}, 0)
		}
	}
	// 最后一批率先创建
	if count > 0 {
		counts = prependInt64(counts, count)
		ds = prependSliceMapStringInterface(ds, d)
	}

	for i := 0; i < len(counts); i++ {
		batch := &mqbo.PushBatchCreateIssue{
			TraceId:     threadlocal.GetTraceId(),
			WithoutAuth: withoutAuth,
			TriggerBy:   triggerBy,
			AsyncTaskId: asyncTaskId,
			Count:       counts[i],
			Req: &projectvo.BatchCreateIssueReqVo{
				OrgId:         req.OrgId,
				UserId:        req.UserId,
				AppId:         req.AppId,
				ProjectId:     req.ProjectId,
				TableId:       req.TableId,
				BeforeDataId:  req.BeforeDataId,
				AfterDataId:   req.AfterDataId,
				Data:          ds[i],
				IsIdGenerated: req.IsIdGenerated,
			},
		}
		domain.PushBatchCreateIssue(batch)
	}
}

func BatchCreateIssue(req *projectvo.BatchCreateIssueReqVo, withoutAuth bool, triggerBy *projectvo.TriggerBy,
	asyncTaskId string, createCount int64) (
	[]map[string]interface{}, map[string]*uservo.MemberDept, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	ctx := &BatchCreateIssueContext{
		TriggerBy:   triggerBy,
		AsyncTaskId: asyncTaskId,
		Req:         req,
	}
	log.Infof("[BatchCreateIssue] req: %v, withoutAuth: %v, triggerBy: %v, asyncTaskId: %v",
		json.ToJsonIgnoreError(req), withoutAuth, triggerBy, asyncTaskId)
	if len(ctx.Req.Data) == 0 {
		return nil, nil, nil, nil
	}

	// 1. 检查参数 组装数据
	errSys := ctx.prepare()
	if errSys != nil {
		log.Error(errSys)
		goto ReturnErr
	}

	// 2. 判断权限
	if !withoutAuth {
		errSys = ctx.checkAuth()
		if errSys != nil {
			log.Error(errSys)
			goto ReturnErr
		}
	}

	// 3. 处理批量创建
	errSys = ctx.process()
	if errSys != nil {
		log.Error(errSys)
		goto ReturnErr
	}

	// 更新异步任务进度
	domain.UpdateAsyncTaskWithSucc(req.OrgId, req.ProjectId, req.TableId, req.UserId, asyncTaskId,
		triggerBy.TriggerBy == consts.TriggerByApplyTemplate, createCount)
	return ctx.CreateResults, ctx.UserDepts, ctx.RelateData, nil

ReturnErr:
	// 更新异步任务进度
	domain.UpdateAsyncTaskWithError(req.OrgId, req.ProjectId, req.TableId, req.UserId, asyncTaskId,
		triggerBy.TriggerBy == consts.TriggerByApplyTemplate, errSys, createCount)
	return nil, nil, nil, errSys
}

func (ctx *BatchCreateIssueContext) prepare() errs.SystemErrorInfo {
	var errSys errs.SystemErrorInfo
	var err error

	ctx.TableColumnsDeep = make(map[string]*lc_table.LcCommonField)
	ctx.CreateColumnMap = make(map[string]struct{})
	ctx.Users = make(map[int64]*uservo.GetAllUserByIdsRespDataUser)
	ctx.Depts = make(map[int64]*uservo.GetAllDeptByIdsRespDataUser)
	ctx.RelatingIssues = make(map[int64]*bo.IssueBo)
	ctx.CreateAttachments = make(map[int64]*bo.IssueAttachments)
	ctx.CreateIssuesMap = make(map[int64]*CreateIssueVo)
	ctx.AllParentIds = make(map[int64]struct{})
	ctx.AllParentPaths = make(map[int64]string)
	ctx.Now = time.Now().Format(consts.AppTimeFormat)

	// 操作人的信息也需要拉取，后续需要用到
	ctx.UserIds = append(ctx.UserIds, ctx.Req.UserId)
	if ctx.TriggerBy.TriggerUserId > 0 {
		ctx.UserIds = append(ctx.UserIds, ctx.TriggerBy.TriggerUserId)
	}

	if ctx.Req.AppId < 0 || ctx.Req.ProjectId < 0 || ctx.Req.TableId < 0 {
		log.Errorf("[BatchCreateIssue] 参数错误 org:%d app:%d proj:%d table:%d user:%d",
			ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId)
		return errs.ParamError
	}

	// 获取org base info
	resp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: ctx.Req.OrgId})
	if resp.Failure() {
		log.Errorf("[BatchCreateIssue] 获取Org失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
			ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
		return resp.Error()
	}
	ctx.OrgBaseInfo = resp.BaseOrgInfo

	// 兼容下有appId但没有projectId的情况
	if ctx.Req.AppId > 0 && ctx.Req.ProjectId == 0 {
		ctx.Project, errSys = domain.GetProjectByAppId(ctx.Req.OrgId, ctx.Req.AppId)
		if errSys != nil {
			log.Errorf("[BatchCreateIssue] GetProjectByAppId org:%d app:%d proj:%d table:%d user:%d, err: %v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
			return errs.BuildSystemErrorInfo(errs.ProjectDomainError, errSys)
		}
		ctx.Req.ProjectId = ctx.Project.Id
	}

	// 获取project
	preCode := consts.NoProjectPreCode
	if ctx.Req.ProjectId > 0 {
		if ctx.Project == nil {
			ctx.Project, errSys = domain.GetProjectSimple(ctx.Req.OrgId, ctx.Req.ProjectId)
			if errSys != nil {
				log.Errorf("[BatchCreateIssue] LoadProjectAuthBo org:%d app:%d proj:%d table:%d user:%d, err: %v",
					ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
				return errs.BuildSystemErrorInfo(errs.ProjectDomainError, errSys)
			}
		}

		// 校验项目是否归档
		if ctx.Project.IsFiling == consts.AppIsFilling {
			return errs.BuildSystemErrorInfo(errs.ProjectIsArchivedWhenModifyIssue)
		}

		if ctx.Project.PreCode != "" {
			preCode = ctx.Project.PreCode
		} else {
			preCode = fmt.Sprintf("$%d", ctx.Req.ProjectId)
		}

		// 兼容下没有传appId的情况
		if ctx.Req.AppId == 0 {
			ctx.Req.AppId = ctx.Project.AppId
		}
	}

	// 这几个id要么都是0 要么都不是0
	if ctx.Req.AppId == 0 || ctx.Req.ProjectId == 0 || ctx.Req.TableId == 0 {
		if ctx.Req.AppId > 0 || ctx.Req.ProjectId > 0 || ctx.Req.TableId > 0 {
			return errs.ParamError
		}
	}

	// 获取表基本信息
	if ctx.Req.TableId > 0 {
		ctx.TableMeta, errSys = domain.GetTableByTableId(ctx.Req.OrgId, ctx.Req.UserId, ctx.Req.TableId)
		if errSys != nil {
			log.Errorf("[BatchCreateIssue] 获取表信息失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
			return errSys
		}
	}

	// 获取表头
	ctx.TableColumns, errSys = domain.GetTableColumnsMap(ctx.Req.OrgId, ctx.Req.TableId, nil, true)
	if errSys != nil {
		log.Errorf("[BatchCreateIssue] 获取表头失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
			ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
		return errSys
	}
	for key, column := range ctx.TableColumns {
		columnDeep := &lc_table.LcCommonField{}
		if err = copyer.Copy(column, columnDeep); err == nil {
			ctx.TableColumnsDeep[key] = columnDeep
		} else {
			log.Errorf("[BatchCreateIssue] 复制表头失败 org:%d app:%d proj:%d table:%d user:%d, column: %v, err: %v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, json.ToJsonIgnoreError(column), errSys)
		}
	}

	// 获取IssueStatus表头，如果获取不到就忽略，这个时候没有issueStatus表头
	if issueStatusColumn, ok := ctx.TableColumns[consts.BasicFieldIssueStatus]; ok {
		ctx.AllIssueStatus = domain.GetStatusListFromStatusColumn(issueStatusColumn)
		log.Infof("[BatchCreateIssue] issueStatus 表头: %v, org:%d app:%d proj:%d table:%d",
			json.ToJsonIgnoreError(ctx.AllIssueStatus), ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId)
	} else {
		log.Infof("[BatchCreateIssue] no need issueStatus, org:%d app:%d proj:%d table:%d", ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId)
	}

	// 参数检查 step 1
	childLevels := make(map[int64]int)
	for _, data := range ctx.Req.Data {
		level, errSys := ctx.prepareIssueStep1(data, 1, 1)
		if errSys != nil {
			log.Infof("[BatchCreateIssue] prepareIssueStep1 failed: %v", errSys)
			return errSys
		}
		parentId := cast.ToInt64(data[consts.BasicFieldParentId])
		if parentId > 0 {
			childLevels[parentId] = level
		}
	}

	// 检查付费任务数限制
	if ctx.CreateCount > 0 && ctx.TriggerBy.TriggerBy != consts.TriggerByApplyTemplate {
		errSys = domain.AuthPayTask(ctx.Req.OrgId, consts.FunctionTaskLimit, int(ctx.CreateCount))
		if errSys != nil {
			log.Errorf("[BatchCreateIssue] 校验付费任务数限制失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
			return errSys
		}
	}

	// 查询顶级子任务的父任务path
	if len(ctx.AllParentIds) > 0 {
		var parentIds []int64
		for id, _ := range ctx.AllParentIds {
			parentIds = append(parentIds, id)
		}
		data, errSys := domain.GetIssueInfosMapLcByIssueIds(ctx.Req.OrgId, ctx.Req.UserId, parentIds,
			lc_helper.ConvertToFilterColumn(consts.BasicFieldPath),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId))
		if errSys != nil {
			log.Errorf("[BatchCreateIssue] GetIssueInfosMapLcByIssueIds org:%d user:%d issueIds:%v, err: %v",
				ctx.Req.OrgId, ctx.Req.UserId, parentIds, errSys)
			return errSys
		}
		for _, d := range data {
			issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
			path := cast.ToString(d[consts.BasicFieldPath])
			tableId := cast.ToInt64(d[consts.BasicFieldTableId])

			// 检查父任务的tableId是否一致
			if tableId != ctx.Req.TableId {
				log.Errorf("[BatchCreateIssue] 父任务tableId不一致 org:%d app:%d proj:%d table:%d user:%d, parentId: %d parentTableId: %d",
					ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, issueId, tableId)
				return errs.ParamError
			}

			if issueId > 0 {
				// 顶级子任务需要继承其父任务的path
				ctx.AllParentPaths[issueId] = fmt.Sprintf("%s%d,", path, issueId)
			}

			// 检查总深度是否超出允许范围
			level := strings.Count(path, ",")
			if childLevel, ok := childLevels[issueId]; ok {
				if level+childLevel > consts.IssueLevel {
					return errs.IssueLevelOutLimit
				}
			}
		}
	}

	// 校验成员/部门是否存在
	ctx.UserIds = slice.SliceUniqueInt64(ctx.UserIds)
	ctx.DeptIds = slice.SliceUniqueInt64(ctx.DeptIds)
	if len(ctx.UserIds) > 0 {
		resp := userfacade.GetAllUserByIds(ctx.Req.OrgId, ctx.UserIds)
		if resp.Failure() {
			log.Errorf("[BatchCreateIssue] GetAllUserByIds err: %v", resp.Error())
			return resp.Error()
		}
		for i := 0; i < len(resp.Data); i++ {
			ctx.Users[resp.Data[i].Id] = &resp.Data[i]
		}
	}
	if len(ctx.DeptIds) > 0 {
		resp := userfacade.GetAllDeptByIds(ctx.Req.OrgId, ctx.DeptIds)
		if resp.Failure() {
			log.Errorf("[BatchCreateIssue] GetAllDeptByIds err: %v", resp.Error())
			return resp.Error()
		}
		for i := 0; i < len(resp.Data); i++ {
			ctx.Depts[resp.Data[i].Id] = &resp.Data[i]
		}
	}

	// 查询before/after的order
	var baDataIds []int64
	if ctx.Req.BeforeDataId > 0 {
		baDataIds = append(baDataIds, ctx.Req.BeforeDataId)
	}
	if ctx.Req.AfterDataId > 0 {
		baDataIds = append(baDataIds, ctx.Req.AfterDataId)
	}
	if len(baDataIds) == 0 {
		order, errSys := domain.GetTableIssueMaxOrder(ctx.Req.OrgId, ctx.Req.UserId)
		if errSys != nil {
			log.Errorf("[BatchCreateIssue] GetTableIssueMaxOrder err: %v", errSys)
			return errSys
		}
		ctx.BeforeOrder = &order
	} else {
		data, errSys := domain.GetIssueInfosMapLcByDataIds(ctx.Req.OrgId, ctx.Req.UserId, baDataIds, lc_helper.ConvertToFilterColumn(consts.BasicFieldOrder))
		if errSys != nil {
			log.Errorf("[BatchCreateIssue] GetIssueInfosMapLcByDataIds err: %v", errSys)
			return errSys
		}
		for _, d := range data {
			dataId := cast.ToInt64(d[consts.BasicFieldId])
			order := cast.ToFloat64(d[consts.BasicFieldOrder])
			if dataId == ctx.Req.BeforeDataId {
				ctx.BeforeOrder = &order
			} else if dataId == ctx.Req.AfterDataId {
				ctx.AfterOrder = &order
			}
		}
		// 做个保护
		if ctx.BeforeOrder != nil &&
			ctx.AfterOrder != nil &&
			*ctx.BeforeOrder > *ctx.AfterOrder {
			temp := ctx.BeforeOrder
			ctx.BeforeOrder = ctx.AfterOrder
			ctx.AfterOrder = temp
		}
	}

	// 生成issue ids
	if !ctx.Req.IsIdGenerated {
		ctx.IssueIds, errSys = idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssue, int(ctx.CreateCount))
		if errSys != nil {
			log.Errorf("[BatchCreateIssue] ApplyMultiplePrimaryIdRelaxed failed, count: %v, err: %v",
				ctx.CreateCount, errSys)
			return errSys
		}

		// 生成issue codes
		ctx.IssueCodes, errSys = idfacade.ApplyMultipleIdRelaxed(ctx.Req.OrgId, preCode, "", int64(ctx.CreateCount))
		if errSys != nil {
			log.Errorf("[BatchCreateIssue] ApplyMultiplePrimaryIdRelaxed failed, count: %v, err: %v",
				ctx.CreateCount, errSys)
			return errSys
		}
	}

	//// 生成data ids
	//sfNode, err := snowflake.NewNode(1)
	//if err != nil {
	//	log.Errorf("[BatchCreateIssue] snowflake.NewNode failed, err: %v", err)
	//	return errs.ServerError
	//}
	//for i := 0; i < ctx.CreateCount; i++ {
	//	ctx.DataIds = append(ctx.DataIds, sfNode.Generate().Int64())
	//}

	// 参数检查 step 2 (赋予issue id/issue code以及其他字段初始值)
	for _, data := range ctx.Req.Data {
		// 顶级任务要继承其父任务的path
		var parentId int64
		path := "0,"
		if v, ok := data[consts.BasicFieldParentId]; ok {
			if v1, ok := ctx.AllParentPaths[cast.ToInt64(v)]; ok {
				parentId = cast.ToInt64(v)
				path = v1
			}
		}
		ctx.prepareIssueStep2(data, parentId, path)
	}

	//// step 3 (处理关联/前后置/单向关联id映射)
	//for _, data := range ctx.Req.Data {
	//	ctx.prepareIssueStep3(data)
	//}

	// 收集列信息
	for column, _ := range ctx.CreateColumnMap {
		ctx.CreateColumns = append(ctx.CreateColumns, column)
	}

	return nil
}

// 单条任务检查参数 step - 1 (收集成员信息, 参数校验)
// 把最大深度返回，结合parent的path深度即可判断深度是否超出
func (ctx *BatchCreateIssueContext) prepareIssueStep1(data map[string]interface{}, level, topLevel int) (int, errs.SystemErrorInfo) {
	ctx.CreateCount += 1

	// 如果不算顶级任务的父任务深度就已经超出限制，那么肯定可以直接返回了
	if level > consts.IssueLevel {
		return 0, errs.IssueLevelOutLimit
	}

	for k, v := range data {
		// 收集顶级任务的父任务信息
		if k == consts.BasicFieldParentId && level == topLevel {
			parentId := cast.ToInt64(v)
			if parentId > 0 {
				ctx.AllParentIds[parentId] = struct{}{}
			}
		}

		if column, ok := ctx.TableColumns[k]; ok {
			// 校验字段值合法性
			newV, errSys := checkColumnValue(k, v)
			if errSys != nil {
				return 0, errSys
			}

			// 收集列信息
			ctx.CreateColumnMap[k] = struct{}{}

			data[k] = newV

			// 收集成员/部门字段信息
			if column.Field.Type == tablePb.ColumnType_member.String() {
				idStrs := cast.ToStringSlice(newV)
				ids := businees.LcMemberToUserIds(idStrs)
				if len(ids) > 0 {
					ctx.UserIds = append(ctx.UserIds, ids...)
				}
			} else if column.Field.Type == tablePb.ColumnType_dept.String() {
				idStrs := cast.ToStringSlice(newV)
				for _, id := range idStrs {
					ctx.DeptIds = append(ctx.DeptIds, cast.ToInt64(id))
				}
			}

		} else if k != consts.BasicFieldRemarkDetail &&
			k != consts.BasicFieldAuditStatus &&
			k != consts.BasicFieldAuditStatusDetail &&
			k != consts.BasicFieldParentId &&
			k != consts.TempFieldChildren &&
			k != consts.TempFieldIssueId &&
			k != consts.TempFieldCode {
			// 删除表头里没有的列，这些数据要么不允许以传入的方式修改，要么是无效数据
			// 注：这里可能存在例外情况，需具体测试确定
			delete(data, k)
		}
	}

	//// 检查起止时间
	//if planStartTimeI, ok := data[consts.BasicFieldPlanStartTime]; ok {
	//	if planEndTimeI, ok := data[consts.BasicFieldPlanEndTime]; ok {
	//		planStartTime := cast.ToTime(planStartTimeI)
	//		planEndTime := cast.ToTime(planEndTimeI)
	//		if planStartTime.After(planEndTime) {
	//			return 0, errs.PlanEndTimeInvalidError
	//		}
	//	}
	//}

	// 递归处理子任务
	if childrenI, ok := data[consts.TempFieldChildren]; ok {
		var children []map[string]interface{}
		err := copyer.Copy(childrenI, &children)
		if err != nil {
			log.Errorf("[BatchCreateIssue] copy children failed, children: %v, err: %v",
				json.ToJsonIgnoreError(childrenI), err)
			return 0, errs.BuildSystemErrorInfo(errs.ParamError, err)
		} else {
			maxChildLevel := level
			for _, c := range children {
				if l, errSys := ctx.prepareIssueStep1(c, level+1, topLevel); errSys != nil {
					return 0, errSys
				} else {
					if l > maxChildLevel {
						maxChildLevel = l
					}
				}
			}
			level = maxChildLevel // 最大深度
		}

		data[consts.TempFieldChildren] = children
	}
	return level, nil
}

// 单条任务组装 step - 2（到这里不应该存在出错情况，可能出错的判断放到step - 1）
//  1. 过滤无效成员
//  2. 组装数据，填入issueId dataId code path parentId
func (ctx *BatchCreateIssueContext) prepareIssueStep2(data map[string]interface{}, parentId int64, path string) {
	var updatedColumns []string
	for k, v := range data {
		if column, ok := ctx.TableColumnsDeep[k]; ok {
			// 过滤无效成员/部门
			if column.Field.Type == tablePb.ColumnType_member.String() &&
				k != consts.BasicFieldCreator &&
				k != consts.BasicFieldUpdator {
				ids := businees.LcMemberToUserIds(cast.ToStringSlice(v))
				idStrs := make([]string, 0)
				for _, id := range ids {
					if _, ok := ctx.Users[id]; ok {
						idStrs = append(idStrs, businees.FormatUserId(id))
					}
				}

				// 不允许多选时，只取第一个人
				isMultiple := column.Field.Props.Member.Multiple || column.Field.Props.Multiple
				if !isMultiple && len(idStrs) > 0 {
					idStrs = idStrs[0:1]
				}

				data[k] = idStrs
			} else if column.Field.Type == tablePb.ColumnType_dept.String() {
				ids := cast.ToStringSlice(v)
				idStrs := make([]string, 0)
				for _, id := range ids {
					i := cast.ToInt64(id)
					if _, ok := ctx.Depts[i]; ok {
						idStrs = append(idStrs, id)
					}
				}
				data[k] = idStrs
			}

			updatedColumns = append(updatedColumns, k)
			if k == consts.BasicFieldIssueStatus {
				updatedColumns = append(updatedColumns, consts.BasicFieldIssueStatusType)
			}
		}
	}

	// 填充issueId dataId等信息
	var issueId int64
	var code string
	if !ctx.Req.IsIdGenerated {
		issueId = ctx.IssueIds.Ids[ctx.CreateIndex].Id
		code = ctx.IssueCodes.Ids[ctx.CreateIndex].Code
		//dataId := ctx.DataIds[ctx.CreateIndex]
	} else {
		issueId = cast.ToInt64(data[consts.TempFieldIssueId])
		code = cast.ToString(data[consts.TempFieldCode])
	}
	data[consts.BasicFieldIssueId] = issueId
	data[consts.BasicFieldCode] = code
	//data[consts.BasicFieldId] = dataId
	//data[consts.BasicFieldDataId] = dataId
	data[consts.BasicFieldParentId] = parentId // 父任务Id
	data[consts.BasicFieldPath] = path
	data[consts.BasicFieldOrgId] = ctx.Req.OrgId
	if ctx.Req.ProjectId > 0 {
		data[consts.BasicFieldAppId] = cast.ToString(ctx.Project.AppId) // 这里用项目对应的appId，避免传过来的appId是menuAppId例如镜像视图appId
	} else {
		data[consts.BasicFieldAppId] = "0"
	}
	data[consts.BasicFieldProjectId] = ctx.Req.ProjectId
	data[consts.BasicFieldTableId] = cast.ToString(ctx.Req.TableId)
	data[consts.BasicFieldRecycleFlag] = 2
	if ctx.BeforeOrder != nil && ctx.AfterOrder != nil {
		step := (*ctx.AfterOrder - *ctx.BeforeOrder) / float64(ctx.CreateCount+1)
		data[consts.BasicFieldOrder] = *ctx.BeforeOrder + float64(ctx.CreateIndex+1)*step
	} else if ctx.BeforeOrder != nil {
		data[consts.BasicFieldOrder] = *ctx.BeforeOrder + float64(ctx.CreateCount-ctx.CreateIndex)*65536
	} else if ctx.AfterOrder != nil {
		data[consts.BasicFieldOrder] = *ctx.AfterOrder - float64(ctx.CreateIndex+1)*65536
	}

	delete(data, consts.TempFieldIssueId)
	delete(data, consts.TempFieldCode)

	ctx.CreateIndex += 1

	// ownerChangeTime
	if ownerIdI, ok := data[consts.BasicFieldOwnerId]; ok {
		ownerIds := cast.ToStringSlice(ownerIdI)
		if len(ownerIds) > 0 {
			data[consts.BasicFieldOwnerChangeTime] = ctx.Now
		}
	}

	// 创建人
	data[consts.BasicFieldCreator] = cast.ToString(ctx.Req.UserId)
	// 创建时间
	data[consts.BasicFieldCreateTime] = ctx.Now
	// 更新人
	data[consts.BasicFieldUpdator] = cast.ToString(ctx.Req.UserId)
	// 更新时间
	data[consts.BasicFieldUpdateTime] = ctx.Now

	// 字段默认值
	if ctx.TriggerBy.TriggerBy == consts.TriggerByNormal ||
		ctx.TriggerBy.TriggerBy == consts.TriggerByAutomation {
		for _, column := range ctx.TableColumnsDeep {
			if column.Field.Props.Default != nil {
				if _, ok := data[column.Name]; !ok {
					data[column.Name] = column.Field.Props.Default
					updatedColumns = append(updatedColumns, column.Name) // 设置了默认值的字段，也算进被设置值的字段
				}
			}
		}
	}
	// 任务状态的默认值强制写进去，无视triggerBy来源
	if column, ok := ctx.TableColumnsDeep[consts.BasicFieldIssueStatus]; ok {
		if column.Field.Props.Default != nil {
			if _, ok := data[column.Name]; !ok {
				data[column.Name] = column.Field.Props.Default
			}
		}
	}

	// 生成issueBo
	createIssueVo := &CreateIssueVo{}
	createIssueVo.IssueId = issueId
	createIssueVo.Data = data
	createIssueVo.UpdateColumns = updatedColumns
	createIssueVo.IssueBo, _ = domain.ConvertIssueDataToIssueBo(data)
	createIssueVo.IssueBo.Version = 1
	createIssueVo.IssueBo.IsDelete = consts.AppIsNoDelete
	createIssueVo.RelatingChange = make(map[string]*bo.RelatingChangeBo)
	createIssueVo.BaRelatingChange = make(map[string]*bo.RelatingChangeBo)
	createIssueVo.SingleRelatingChange = make(map[string]*bo.RelatingChangeBo)
	ctx.CreateIssueIds = append(ctx.CreateIssueIds, issueId)
	ctx.CreateIssues = append(ctx.CreateIssues, createIssueVo)
	ctx.CreateIssuesMap[issueId] = createIssueVo

	// 递归处理子任务
	if childrenI, ok := data[consts.TempFieldChildren]; ok {
		var children []map[string]interface{}
		err := copyer.Copy(childrenI, &children)
		if err != nil {
			// 理论上这里不可能出错，step1已经处理过一遍了
			log.Errorf("[BatchCreateIssue] copy children failed, children: %v, err: %v",
				json.ToJsonIgnoreError(childrenI), err)
		}
		for _, c := range children {
			ctx.prepareIssueStep2(c, issueId, fmt.Sprintf("%s%d,", path, issueId))
		}

		// 删除无码数据结构中的子任务信息
		delete(data, consts.TempFieldChildren)
	}
}

//// 单条任务组装 step - 3（到这里不应该存在出错情况，可能出错的判断放到step - 1）
////  1. 处理关联/前后置/单向关联字段的动态复制的id替换
//func (ctx *BatchCreateIssueContext) prepareIssueStep3(data map[string]interface{}) {
//	issueId := cast.ToInt64(data[consts.BasicFieldIssueId])
//	for k, v := range data {
//		if column, ok := ctx.TableColumns[k]; ok {
//			if column.Field.Type == tablePb.ColumnType_relating.String() ||
//				column.Field.Type == tablePb.ColumnType_baRelating.String() ||
//				column.Field.Type == tablePb.ColumnType_singleRelating.String() {
//				relating := &bo.RelatingIssue{}
//				if err := jsonx.Copy(v, relating); err != nil {
//					log.Errorf("[BatchCreateIssue] parse relating column failed, issueId: %v, columnId: %v, value: %v, err: %v",
//						issueId, k, json.ToJsonIgnoreError(v), err)
//					delete(data, k)
//				} else {
//					for _, uuid := range relating.LinkToUuid {
//						if id, ok := ctx.AllUuids[uuid]; ok {
//							relating.LinkTo = append(relating.LinkTo, cast.ToString(id))
//						}
//					}
//					for _, uuid := range relating.LinkFromUuid {
//						if id, ok := ctx.AllUuids[uuid]; ok {
//							relating.LinkFrom = append(relating.LinkFrom, cast.ToString(id))
//						}
//					}
//					relating.LinkToUuid = nil
//					relating.LinkFromUuid = nil
//					relating.LinkTo = slice.SliceUniqueString(relating.LinkTo)
//					relating.LinkFrom = slice.SliceUniqueString(relating.LinkFrom)
//					data[k] = relating
//				}
//			}
//		}
//	}
//}

func (ctx *BatchCreateIssueContext) checkAuth() errs.SystemErrorInfo {
	var errSys errs.SystemErrorInfo
	//var err error

	// 未归属项目的任务，不校验权限
	if ctx.Req.AppId == 0 {
		return nil
	}

	projectAuth, errSys := domain.LoadProjectAuthBo(ctx.Req.OrgId, ctx.Req.ProjectId)
	if errSys != nil {
		log.Errorf("[BatchCreateIssue] LoadProjectAuthBo org:%d app:%d proj:%d table:%d user:%d, err: %v",
			ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, errSys)
	}

	// 校验字段权限
	errSys = orgfacade.AuthenticateRelaxed(ctx.Req.OrgId, ctx.Req.UserId, projectAuth, nil, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Create, ctx.CreateColumns)
	if errSys != nil {
		log.Errorf("[BatchCreateIssue] 权限校验 org:%d app:%d proj:%d table:%d user:%d, err: %v",
			ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
		if errSys.Code() == errs.NoOperationPermissions.Code() {
			return errs.NoOperationPermissionForProject
		} else {
			return errSys
		}
	}

	return nil
}

func (ctx *BatchCreateIssueContext) process() errs.SystemErrorInfo {
	var errSys errs.SystemErrorInfo

	// 处理更新
	errSys = ctx.processCreate()
	if errSys != nil {
		return errSys
	}

	// 保存数据
	errSys = ctx.saveToDB()
	if errSys != nil {
		return errSys
	}

	// 后续处理（异步）
	asyn.Execute(ctx.processHooks)
	return nil
}

func (ctx *BatchCreateIssueContext) processCreate() errs.SystemErrorInfo {
	for _, ci := range ctx.CreateIssues {
		createIssue := ci

		if errSys := createIssue.processIssueAuditStatus(ctx); errSys != nil {
			return errSys
		}
		if errSys := createIssue.processOwner(ctx); errSys != nil {
			return errSys
		}
		if errSys := createIssue.processFollower(ctx); errSys != nil {
			return errSys
		}
		if errSys := createIssue.processAuditor(ctx); errSys != nil {
			return errSys
		}
		if errSys := createIssue.processRemark(ctx); errSys != nil {
			return errSys
		}
		if errSys := createIssue.processDocument(ctx); errSys != nil {
			return errSys
		}
		if errSys := createIssue.processRelating(ctx, tablePb.ColumnType_relating.String()); errSys != nil {
			return errSys
		}
		if errSys := createIssue.processRelating(ctx, tablePb.ColumnType_baRelating.String()); errSys != nil {
			return errSys
		}
		if errSys := createIssue.processRelating(ctx, tablePb.ColumnType_singleRelating.String()); errSys != nil {
			return errSys
		}

		createIssue.processTrends(ctx)
	}
	return nil
}

func (ctx *BatchCreateIssueContext) saveToDB() errs.SystemErrorInfo {
	var errSys errs.SystemErrorInfo
	// 组装需要创建的issue relations id
	if len(ctx.CreateIssueRelations) > 0 {
		relationIds, errSys := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, len(ctx.CreateIssueRelations))
		if errSys != nil {
			return errs.BuildSystemErrorInfo(errs.ApplyIdError, errSys)
		}
		for i := 0; i < len(ctx.CreateIssueRelations); i++ {
			ctx.CreateIssueRelations[i].Id = relationIds.Ids[i].Id
		}
	}

	appId := ctx.Req.AppId
	if appId == 0 {
		appId, errSys = domain.GetOrgSummaryAppId(ctx.Req.OrgId)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}
	}

	lcReq := formvo.LessCreateIssueReq{
		OrgId:          ctx.Req.OrgId,
		AppId:          appId, // 如果是未归属项目的任务，这里传汇总表appId
		TableId:        ctx.Req.TableId,
		UserId:         ctx.Req.UserId,
		Import:         ctx.TriggerBy.TriggerBy == consts.TriggerByImport,
		CreateTemplate: ctx.TriggerBy.IsCreateTemplate,
		Asc:            true, // form目前是写死true的，这里传不传其实无所谓，写在这里减少理解上的困扰
	}
	if ctx.Req.BeforeDataId > 0 {
		lcReq.BeforeId = ctx.Req.BeforeDataId
	}
	if ctx.Req.AfterDataId > 0 {
		lcReq.AfterId = ctx.Req.AfterDataId
	}
	if ctx.Project != nil {
		lcReq.RedirectIds = []int64{ctx.Project.AppId} // 这里的appId是真正的appId，用于区别镜像视图appId (menuAppId)
	}
	for _, ci := range ctx.CreateIssues {
		createIssue := ci
		lcReq.Form = append(lcReq.Form, createIssue.Data)
	}

	sysErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		fs := make([]batch.TxExecFunc, 0)

		// 创建任务
		//if len(ctx.CreateIssues) > 0 {
		//	fs = append(fs, func(tx sqlbuilder.Tx) errors.SystemErrorInfo {
		//		var issuePos []*po.PpmPriIssue
		//		for _, ci := range ctx.CreateIssues {
		//			createIssue := ci
		//
		//			issuePo := &po.PpmPriIssue{}
		//			err := copyer.Copy(createIssue.IssueBo, issuePo)
		//			if err != nil {
		//				log.Infof("[BatchCreateIssue] ConvertObject err: %v", err)
		//				return errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		//			}
		//			if len(createIssue.IssueBo.OwnerIdI64) > 0 {
		//				issuePo.Owner = createIssue.IssueBo.OwnerIdI64[0]
		//			}
		//			issuePos = append(issuePos, issuePo)
		//		}
		//
		//		err := mysql.TransBatchInsert(tx, &po.PpmPriIssue{}, slice.ToSlice(issuePos))
		//		if err != nil {
		//			log.Errorf("[BatchCreateIssue] batch insert Issues, err: %v", err)
		//			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		//		}
		//		return nil
		//	})
		//}

		// 创建issue relations
		if len(ctx.CreateIssueRelations) > 0 {
			fs = append(fs, func(tx sqlbuilder.Tx) errors.SystemErrorInfo {
				log.Infof("[BatchCreateIssue] batch insert IssueRelations: %v", json.ToJsonIgnoreError(ctx.CreateIssueRelations))
				err := mysql.TransBatchInsert(tx, &po.PpmPriIssueRelation{}, slice.ToSlice(ctx.CreateIssueRelations))
				if err != nil {
					log.Errorf("[BatchCreateIssue] batch insert IssueRelations, err: %v", err)
					return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
				}
				return nil
			})
		}

		// 创建resource relations
		fs = append(fs, func(tx sqlbuilder.Tx) errors.SystemErrorInfo {
			for issueId, v := range ctx.CreateAttachments {
				if len(v.ResourceIds) > 0 {
					errSys := domain.AddResourceRelation(ctx.Req.OrgId, ctx.Req.UserId, ctx.Req.ProjectId, issueId, v.ResourceIds, consts.OssPolicyTypeLesscodeResource, v.ColumnId)
					if errSys != nil {
						log.Errorf("[BatchCreateIssue] AddResourceRelation err:%v", errSys)
						return errSys
					}
				}
			}
			return nil
		})

		txBatchExecutor := &batch.TxBatchExecutor{}
		txBatchExecutor.Init(tx, 20)
		txBatchExecutor.PushJobs(fs)
		errSys := txBatchExecutor.StartAndWaitFinish()
		if errSys != nil {
			log.Errorf("[BatchCreateIssue] txBatchExecutor.StartAndWaitFinish: %v", errSys)
			return errSys
		}

		// 事务的最后创建无码
		log.Infof("[BatchCreateIssue] LessCreateIssue, req: %v", json.ToJsonIgnoreError(lcReq))
		resp := formfacade.LessCreateIssue(lcReq)
		if resp.Failure() {
			log.Errorf("[BatchCreateIssue] LessCreateIssue: %v, req: %v", resp.Error(), json.ToJsonIgnoreError(lcReq))
			return resp.Error()
		}
		ctx.CreateResults = resp.Data

		return nil
	})
	if sysErr != nil {
		log.Errorf("[BatchCreateIssue] saveToDB: %v", sysErr)
		// 如果客户端频繁保存，会出现这个错误
		// Error 1213: Deadlock found when trying to get lock; try restarting transaction
		if strings.Contains(sysErr.Error(), consts.MySQL_DEADLOCK_ERROR) {
			return errs.RequestFrequentError
		}
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, sysErr)
	}

	// 组装成员/部门信息
	ctx.UserDepts = domain.AssembleUserDepts(ctx.Users, ctx.Depts)
	for _, d := range ctx.CreateResults {
		//domain.AssembleLcDataRelated(d, ctx.TableColumns, ctx.Users, ctx.Depts)
		domain.AssembleDataIds(d)
		delete(d, consts.BasicFieldCollaborators)
	}

	// 关联/引用信息
	if len(ctx.RelatingIssueIds) > 0 {
		req := &formvo.LessIssueListReq{
			Condition: vo.LessCondsData{
				Type:   "in",
				Values: ctx.CreateIssueIds,
				Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
			},
			OrgId:         ctx.Req.OrgId,
			AppId:         ctx.Req.AppId,
			TableId:       ctx.Req.TableId,
			UserId:        ctx.Req.UserId,
			NeedRefColumn: true,
		}
		query := json.ToJsonIgnoreError(req)
		listReq := &tablePb.ListRequest{
			Query:   query,
			TableId: ctx.Req.TableId,
		}
		log.Infof("[BatchCreateIssue] query relate data: %v", query)
		reply := tablefacade.List(req.OrgId, req.UserId, listReq)
		if reply.Failure() {
			log.Errorf("[BatchCreateIssue] query relate data: %v, err: %v", query, reply.Error())
		} else if len(reply.Data.RelateData) > 0 {
			ctx.RelateData = make(map[string]map[string]interface{})
			if err := json.Unmarshal(reply.Data.RelateData, &ctx.RelateData); err != nil {
				log.Errorf("[BatchCreateIssue] json decode relateData failed, %q, err: %v",
					reply.Data.RelateData, reply.Error())
			}
		}
	}

	return nil
}

func (ctx *BatchCreateIssueContext) processHooks() {
	ctx.RelatingIssueIds = slice.SliceUniqueInt64(ctx.RelatingIssueIds)
	if len(ctx.RelatingIssueIds) > 0 {
		data, errSys := domain.GetIssueInfosMapLcByDataIds(ctx.Req.OrgId, ctx.Req.UserId, ctx.RelatingIssueIds)
		if errSys != nil {
			log.Errorf("[BatchUpdateIssue] 相关关联前后置任务数据失败 org:%d app:%d proj:%d table:%d user:%d, dept:%v err: %v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, ctx.DeptIds, errSys)
		} else {
			for _, d := range data {
				issueBo, errSys := domain.ConvertIssueDataToIssueBo(d)
				if errSys == nil {
					ctx.RelatingIssues[issueBo.Id] = issueBo
				}
			}
		}
	}

	// 上报事件
	asyn.Execute(func() {
		ctx.reportEvent()
	})

	// 动态相关
	asyn.Execute(func() {
		// 保存动态
		for _, t := range ctx.IssueTrends {
			trend := t
			domain.PushIssueTrends(trend)
		}

		// 个人卡片推送
		asyn.Execute(func() {
			for _, t := range ctx.IssueTrends {
				trend := t
				if ctx.AsyncTaskId == "" {
					domain.PushIssueThirdPlatformNotice(trend)
				}
			}
		})

	})
	// 钉钉酷应用顶部卡片更新
	updateDingTopCard(ctx.Req.OrgId, ctx.Req.ProjectId)
}

func (ctx *BatchCreateIssueContext) reportEvent() {
	mqttFlag := ctx.TriggerBy.TriggerBy != consts.TriggerByImport && ctx.TriggerBy.TriggerBy != consts.TriggerByApplyTemplate

	openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
	openTraceIdStr := cast.ToString(openTraceId)

	// 上报创建任务事件
	for i := len(ctx.CreateResults) - 1; i >= 0; i-- {
		data := ctx.CreateResults[i]
		dataId := cast.ToInt64(data[consts.BasicFieldDataId])
		issueId := cast.ToInt64(data[consts.BasicFieldIssueId])

		e := &commonvo.DataEvent{
			OrgId:     ctx.Req.OrgId,
			AppId:     ctx.Req.AppId,
			ProjectId: ctx.Req.ProjectId,
			TableId:   ctx.Req.TableId,
			DataId:    dataId,
			IssueId:   issueId,
			UserId:    ctx.Req.UserId,
			New:       data,
			UserDepts: ctx.UserDepts,
			TriggerBy: ctx.TriggerBy.TriggerBy,
		}
		if createIssueVo, ok := ctx.CreateIssuesMap[issueId]; ok {
			e.UpdatedColumns = createIssueVo.UpdateColumns
		}
		report.ReportDataEvent(msgPb.EventType_DataCreated, openTraceIdStr, e, mqttFlag)
	}

	// 关联前后置更改，给变更的对端任务也要上报更新事件，此时只上报增量事件
	for _, ci := range ctx.CreateIssues {
		createIssue := ci

		for columnId, relatingChange := range ci.RelatingChange {
			relatingLinkTo := map[string]interface{}{
				columnId: bo.RelatingIssue{
					LinkTo: []string{cast.ToString(createIssue.IssueBo.DataId)},
				},
			}
			relatingLinkFrom := map[string]interface{}{
				columnId: bo.RelatingIssue{
					LinkFrom: []string{cast.ToString(createIssue.IssueBo.DataId)},
				},
			}

			for _, dataId := range relatingChange.LinkToAdd {
				if issueBo, ok := ctx.RelatingIssues[dataId]; ok {
					e := &commonvo.DataEvent{
						OrgId:       issueBo.OrgId,
						AppId:       issueBo.AppId,
						ProjectId:   issueBo.ProjectId,
						TableId:     issueBo.TableId,
						DataId:      issueBo.DataId,
						IssueId:     issueBo.Id,
						UserId:      ctx.Req.UserId,
						Incremental: relatingLinkFrom, // 增量-添加
						TriggerBy:   ctx.TriggerBy.TriggerBy,
					}
					report.ReportDataEvent(msgPb.EventType_DataUpdated, openTraceIdStr, e, mqttFlag)
				}
			}
			for _, dataId := range relatingChange.LinkFromAdd {
				if issueBo, ok := ctx.RelatingIssues[dataId]; ok {
					e := &commonvo.DataEvent{
						OrgId:       issueBo.OrgId,
						AppId:       issueBo.AppId,
						ProjectId:   issueBo.ProjectId,
						TableId:     issueBo.TableId,
						DataId:      issueBo.DataId,
						IssueId:     issueBo.Id,
						UserId:      ctx.Req.UserId,
						Incremental: relatingLinkTo, // 增量-添加
						TriggerBy:   ctx.TriggerBy.TriggerBy,
					}
					report.ReportDataEvent(msgPb.EventType_DataUpdated, openTraceIdStr, e, mqttFlag)
				}
			}
		}
	}
}

func (i *CreateIssueVo) processIssueAuditStatus(ctx *BatchCreateIssueContext) errs.SystemErrorInfo {
	// 没有issueStatus表头的话，各种状态都不需要了
	if ctx.AllIssueStatus == nil {
		return nil
	}

	// 组装任务状态
	var issueStatusType int
	var auditStatus int
	if issueStatusI, ok := i.Data[consts.BasicFieldIssueStatus]; ok {
		issueStatus := cast.ToInt64(issueStatusI)
		for _, is := range ctx.AllIssueStatus {
			if issueStatus == is.ID {
				issueStatusType = is.Type
				break
			}
		}
		if issueStatusType > 0 {
			i.Data[consts.BasicFieldIssueStatusType] = issueStatusType
			i.IssueBo.IssueStatusType = int32(issueStatusType)
		}
	}

	// 组装待确认状态
	if _, ok := i.Data[consts.BasicFieldAuditStatus]; !ok {
		if ctx.Req.ProjectId > 0 && ctx.Project.ProjectTypeId == consts.ProjectTypeNormalId {
			// 审批确认项目，设置为待确认状态
			auditStatus = consts.AuditStatusNotView
			if len(i.IssueBo.AuditorIds) == 0 && issueStatusType == consts.StatusTypeComplete {
				// 如果是已完成状态且没有确认人，则设置为审批通过
				auditStatus = consts.AuditStatusPass
			}
		} else {
			// 非审批确认项目
			auditStatus = consts.AuditStatusNoNeed
		}
		i.Data[consts.BasicFieldAuditStatus] = auditStatus
		i.Data[consts.BasicFieldAuditStatusDetail] = make(map[string]int)
		i.IssueBo.AuditStatus = auditStatus
		i.IssueBo.AuditStatusDetail = make(map[string]int)
	}
	return nil
}

func (i *CreateIssueVo) processOwner(ctx *BatchCreateIssueContext) errs.SystemErrorInfo {
	for _, userId := range i.IssueBo.OwnerIdI64 {
		ctx.CreateIssueRelations = append(ctx.CreateIssueRelations, &po.PpmPriIssueRelation{
			//Id:           ids.Ids[idx].Id, // ID到最后统一生成
			OrgId:        i.IssueBo.OrgId,
			ProjectId:    i.IssueBo.ProjectId,
			IssueId:      i.IssueId,
			RelationId:   userId,
			RelationType: consts.IssueRelationTypeOwner,
			Creator:      ctx.Req.UserId,
			CreateTime:   time.Now(),
			Updator:      ctx.Req.UserId,
			UpdateTime:   time.Now(),
			IsDelete:     consts.AppIsNoDelete,
		})
	}
	return nil
}

func (i *CreateIssueVo) processFollower(ctx *BatchCreateIssueContext) errs.SystemErrorInfo {
	for _, userId := range i.IssueBo.FollowerIdsI64 {
		ctx.CreateIssueRelations = append(ctx.CreateIssueRelations, &po.PpmPriIssueRelation{
			//Id:           ids.Ids[idx].Id, // ID到最后统一生成
			OrgId:        i.IssueBo.OrgId,
			ProjectId:    i.IssueBo.ProjectId,
			IssueId:      i.IssueId,
			RelationId:   userId,
			RelationType: consts.IssueRelationTypeFollower,
			Creator:      ctx.Req.UserId,
			CreateTime:   time.Now(),
			Updator:      ctx.Req.UserId,
			UpdateTime:   time.Now(),
			IsDelete:     consts.AppIsNoDelete,
		})
	}
	return nil
}

func (i *CreateIssueVo) processAuditor(ctx *BatchCreateIssueContext) errs.SystemErrorInfo {
	for _, userId := range i.IssueBo.AuditorIdsI64 {
		ctx.CreateIssueRelations = append(ctx.CreateIssueRelations, &po.PpmPriIssueRelation{
			//Id:           ids.Ids[idx].Id, // ID到最后统一生成
			OrgId:        i.IssueBo.OrgId,
			ProjectId:    i.IssueBo.ProjectId,
			IssueId:      i.IssueId,
			RelationId:   userId,
			RelationType: consts.IssueRelationTypeAuditor,
			Creator:      ctx.Req.UserId,
			CreateTime:   time.Now(),
			Updator:      ctx.Req.UserId,
			UpdateTime:   time.Now(),
			IsDelete:     consts.AppIsNoDelete,
		})
	}
	return nil
}

func (i *CreateIssueVo) processRemark(ctx *BatchCreateIssueContext) errs.SystemErrorInfo {
	return nil
}

func (i *CreateIssueVo) processDocument(ctx *BatchCreateIssueContext) errs.SystemErrorInfo {
	for _, columnId := range ctx.CreateColumns {
		if tableColumn, ok := ctx.TableColumns[columnId]; ok {
			if tableColumn.Field.Type == consts.LcColumnFieldTypeDocument {
				if document, ok := i.Data[columnId]; ok {
					attachments := map[string]*bo.Attachments{}
					copyer.Copy(document, &attachments)
					if len(attachments) == 0 {
						delete(i.Data, columnId) // 把无效数据清理掉
					}

					var resourceIds []int64
					for resourceId, _ := range attachments {
						resourceIds = append(resourceIds, cast.ToInt64(resourceId))
					}
					if len(resourceIds) > 0 {
						ctx.CreateAttachments[i.IssueId] = &bo.IssueAttachments{
							ColumnId:    columnId,
							ResourceIds: resourceIds,
						}
					}
				}
			}
		}
	}
	return nil
}

func (i *CreateIssueVo) processRelating(ctx *BatchCreateIssueContext, relatingColumnType string) errs.SystemErrorInfo {
	for columnId, column := range ctx.TableColumns {
		if column.Field.Type == relatingColumnType {
			if _, ok := i.Data[columnId]; ok {
				relatingI := i.Data[columnId]
				relating := &bo.RelatingIssue{}
				jsonx.Copy(relatingI, relating)

				// 收集变化的对端任务ID
				relatingChange := &bo.RelatingChangeBo{
					LinkToAdd:   slice2.StringToInt64Slice(relating.LinkTo),
					LinkFromAdd: slice2.StringToInt64Slice(relating.LinkFrom),
				}
				switch relatingColumnType {
				case tablePb.ColumnType_relating.String():
					i.RelatingChange[columnId] = relatingChange
				case tablePb.ColumnType_baRelating.String():
					i.BaRelatingChange[columnId] = relatingChange
				case tablePb.ColumnType_singleRelating.String():
					i.SingleRelatingChange[columnId] = relatingChange
				}

				ctx.RelatingIssueIds = append(ctx.RelatingIssueIds, relatingChange.LinkToAdd...)
				ctx.RelatingIssueIds = append(ctx.RelatingIssueIds, relatingChange.LinkFromAdd...)
			}
		}
	}
	return nil
}

func (i *CreateIssueVo) processTrends(ctx *BatchCreateIssueContext) {
	ext := bo.TrendExtensionBo{}
	ext.IssueType = "T"
	ext.ObjName = i.IssueBo.Title
	if ctx.TriggerBy.TriggerBy == consts.TriggerByAutomation {
		ext.AutomationInfo = &bo.TrendAutomationInfo{
			WorkflowId:   ctx.TriggerBy.WorkflowId,
			WorkflowName: ctx.TriggerBy.WorkflowName,
			ExecutionId:  ctx.TriggerBy.ExecutionId,
		}
		if ctx.TriggerBy.TriggerUserId > 0 {
			if user, ok := ctx.Users[ctx.TriggerBy.TriggerUserId]; ok {
				ext.AutomationInfo.TriggerUser = &bo.SimpleUserInfoBo{
					Id:     user.Id,
					Name:   user.Name,
					Avatar: user.Avatar,
					Status: user.Status,
				}
			}
		}
	}

	issueTrendsBo := &bo.IssueTrendsBo{
		PushType:              consts.PushTypeCreateIssue,
		OrgId:                 ctx.Req.OrgId,
		OperatorId:            ctx.Req.UserId,
		DataId:                i.IssueBo.DataId,
		IssueId:               i.IssueId,
		ProjectId:             ctx.Req.ProjectId,
		TableId:               ctx.Req.TableId,
		ParentIssueId:         i.IssueBo.ParentId,
		PriorityId:            i.IssueBo.PriorityId,
		ParentId:              i.IssueBo.ParentId,
		IssueTitle:            i.IssueBo.Title,
		IssueStatusId:         i.IssueBo.Status,
		BeforeOwner:           i.IssueBo.OwnerIdI64,
		AfterOwner:            i.IssueBo.OwnerIdI64,
		BeforeChangeFollowers: i.IssueBo.FollowerIdsI64,
		AfterChangeFollowers:  i.IssueBo.FollowerIdsI64,
		BeforeChangeAuditors:  i.IssueBo.AuditorIdsI64,
		AfterChangeAuditors:   i.IssueBo.AuditorIdsI64,
		IssuePlanStartTime:    &i.IssueBo.PlanStartTime,
		IssuePlanEndTime:      &i.IssueBo.PlanEndTime,
		NewValue:              json.ToJsonIgnoreError(i.Data),
		Ext:                   ext,
	}
	log.Infof("[BatchCreateIssue] processTrends, issueId: %d, trend: %v", i.IssueId, json.ToJsonIgnoreError(issueTrendsBo))
	ctx.IssueTrends = append(ctx.IssueTrends, issueTrendsBo)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//  批量更新
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type BatchUpdateIssueContext struct {
	TriggerBy            *projectvo.TriggerBy
	Req                  *projectvo.BatchUpdateIssueReqInnerVo
	OrgBaseInfo          *bo.BaseOrgInfoBo
	Project              *bo.ProjectBo
	TableMeta            *projectvo.TableMetaData
	TableColumns         map[string]*projectvo.TableColumnData // 表头
	AllIssueStatus       []status.StatusInfoBo                 // IssueStatus表头
	Todo                 *automationPb.Todo                    // 待办
	UpdateColumnMap      map[string]struct{}                   // 涉及更新的字段
	UpdateColumns        []string                              // 涉及更新的字段
	UpdateIssueIds       []int64                               // 涉及更新的任务IssueID
	UpdateDataIds        []int64                               // 涉及更新的任务DataID
	UpdateIssues         map[int64]*UpdateIssueVo              // 需更新的任务相关信息
	SqlIssueRelations    []*dao.SqlInfo                        // 需更新的issue relations sql语句
	CreateIssueRelations []*po.PpmPriIssueRelation             // 需创建的issue relations
	CreateRecycleBins    []*po.PpmPrsRecycleBin                // 需要放入回收站的数据
	CreateAttachments    map[int64]*bo.IssueAttachments        // 需要添加的附件
	RecycleAttachments   map[int64]*bo.IssueAttachments        // 需要删除的附件
	DeleteAuditors       map[int64][]int64                     // 需要删除的确认人 issueId, userId
	BeforeOrder          *float64                              // 涉及的beforeId的order
	AfterOrder           *float64                              // 涉及的afterId的order
	Now                  string                                // 当前时间

	// 动态相关
	HandledTrendColumnIds map[string]struct{}                           // 需特别处理动态的字段，其余都走通用流程
	IssueTrends           []*bo.IssueTrendsBo                           // 任务动态结果集
	IterationIds          []int64                                       // 涉及到的迭代Id
	UserIds               []int64                                       // 涉及到的UserIds
	DeptIds               []int64                                       // 涉及到的DeptIds
	Iterations            map[int64]*bo.IterationBo                     // 涉及的迭代
	Users                 map[int64]*uservo.GetAllUserByIdsRespDataUser // 涉及的用户
	Depts                 map[int64]*uservo.GetAllDeptByIdsRespDataUser // 涉及的部门
	RelatingIssueIds      []int64                                       // 关联/前后置/单向关联涉及到的`变化的`任务ID
	AllRelatingIssueIds   []int64                                       // 关联/前后置/单向关联涉及到的`所有的`任务ID
	RelatingIssues        map[int64]*bo.IssueBo                         // 关联前后置修改涉及到的任务

	// 自动排期相关
	AutoScheduleSourceIssues []*UpdateIssueVo

	// 待办相关
	UpdateTodo *commonvo.UpdateTodoReq
	TodoResult int
}

type UpdateIssueVo struct {
	IssueId                        int64
	PData                          map[string]interface{} // 极星issue表需要更新的字段及数据
	LcData                         map[string]interface{} // 无码pg表需要更新的字段及数据
	OldLcData                      map[string]interface{} // 老的无码pg表数据(所有字段)
	OldIssueBo                     *bo.IssueBo
	NewIssueBo                     *bo.IssueBo
	TrendChangeList                []bo.TrendChangeListBo
	UpdateColumns                  []string
	UpdateTitle                    bool
	UpdateOwner                    bool
	UpdateFollower                 bool
	UpdateAuditor                  bool
	NoticeAuditorChange            bool
	NoticeAuditorIssueStatusChange bool
	DeltaPlanEndTime               time.Duration
	RelatingChange                 map[string]*bo.RelatingChangeBo
	BaRelatingChange               map[string]*bo.RelatingChangeBo
	SingleRelatingChange           map[string]*bo.RelatingChangeBo
}

// 自动排期任务信息
type AutoScheduleIssueVo struct {
	Id                 string
	IssueId            int64
	PlanStartTime      types.Time
	PlanEndTime        types.Time
	PlanStartTimeDelta time.Duration
	PlanEndTimeDelta   time.Duration
	BaRelating         map[string][]string
	LinkTo             []*AutoScheduleIssueVo
	LinkFrom           []*AutoScheduleIssueVo
}

// SyncBatchUpdateIssue 每BATCH_SIZE条为一个批次，执行同步批量更新
func SyncBatchUpdateIssue(req *projectvo.BatchUpdateIssueReqInnerVo, withoutAuth bool, triggerBy *projectvo.TriggerBy) errs.SystemErrorInfo {
	var count int64
	var d []map[string]interface{}
	var ds [][]map[string]interface{}
	for _, data := range req.Data {
		d = append(d, data)
		count += 1
		if count >= BATCH_SIZE {
			ds = append(ds, d)

			// 重新累计
			d = make([]map[string]interface{}, 0)
			count = 0
		}
	}
	if count > 0 {
		ds = append(ds, d)
	}

	for i := 0; i < len(ds); i++ {
		r := &projectvo.BatchUpdateIssueReqInnerVo{
			OrgId:          req.OrgId,
			UserId:         req.UserId,
			AppId:          req.AppId,
			ProjectId:      req.ProjectId,
			TableId:        req.TableId,
			Data:           ds[i],
			BeforeDataId:   req.BeforeDataId,
			AfterDataId:    req.AfterDataId,
			TodoId:         req.TodoId,
			TodoOp:         req.TodoOp,
			TrendPushType:  req.TrendPushType,
			TrendExtension: req.TrendExtension,
		}
		err := BatchUpdateIssue(r, withoutAuth, triggerBy)
		if err != nil {
			return err
		}
	}
	return nil
}

func BatchUpdateIssue(reqVo *projectvo.BatchUpdateIssueReqInnerVo, withoutAuth bool, triggerBy *projectvo.TriggerBy) errs.SystemErrorInfo {
	ctx := &BatchUpdateIssueContext{
		TriggerBy: triggerBy,
		Req:       reqVo,
	}
	log.Infof("[BatchUpdateIssue] req: %v", json.ToJsonIgnoreError(reqVo))
	if len(ctx.Req.Data) == 0 && ctx.Req.TodoId == 0 {
		return nil
	}

	// 1. 检查参数 组装数据
	errSys := ctx.prepare()
	if errSys != nil {
		return errSys
	}

	// 没有需要更新的数据就直接返回
	if len(ctx.UpdateIssues) == 0 && ctx.Req.TodoId == 0 {
		return nil
	}

	// 2. 判断权限
	if !withoutAuth {
		errSys = ctx.checkAuth()
		if errSys != nil {
			return errSys
		}
	}

	// 3. 处理批量更新
	errSys = ctx.process()
	if errSys != nil {
		return errSys
	}

	return nil
}

func (ctx *BatchUpdateIssueContext) prepare() errs.SystemErrorInfo {
	var errSys errs.SystemErrorInfo
	var err error
	ctx.UpdateColumnMap = make(map[string]struct{})
	ctx.UpdateIssues = make(map[int64]*UpdateIssueVo)
	ctx.HandledTrendColumnIds = make(map[string]struct{})
	ctx.RecycleAttachments = make(map[int64]*bo.IssueAttachments)
	ctx.CreateAttachments = make(map[int64]*bo.IssueAttachments)
	ctx.DeleteAuditors = make(map[int64][]int64)
	ctx.Users = make(map[int64]*uservo.GetAllUserByIdsRespDataUser)
	ctx.Depts = make(map[int64]*uservo.GetAllDeptByIdsRespDataUser)
	ctx.Iterations = make(map[int64]*bo.IterationBo)
	ctx.RelatingIssues = make(map[int64]*bo.IssueBo)
	ctx.Now = time.Now().Format(consts.AppTimeFormat)

	// 操作人的信息也需要拉取，后续需要用到
	ctx.UserIds = append(ctx.UserIds, ctx.Req.UserId)
	if ctx.TriggerBy.TriggerUserId > 0 {
		ctx.UserIds = append(ctx.UserIds, ctx.TriggerBy.TriggerUserId)
	}

	// 批量更新的任务数上限
	if len(ctx.Req.Data) > BATCH_SIZE {
		return errs.BatchOperateTooManyRows
	}

	if ctx.Req.TodoId > 0 && len(ctx.Req.Data) > 1 {
		return errs.BuildSystemErrorInfoWithMessage(errs.BatchOperateTooManyRows, "Update more than one row in Todo.")
	}

	// 解析数据，并校验
	for _, data := range ctx.Req.Data {
		d := data
		// id必传检查
		if id, ok := d[consts.BasicFieldId]; !ok {
			log.Errorf("[BatchUpdateIssue] 参数校验失败 Id必传 org:%d app:%d proj:%d table:%d user:%d",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId)
			return errs.ReqParamsValidateError
		} else {
			if _, err = cast.ToInt64E(id); err != nil {
				log.Errorf("[BatchUpdateIssue] 参数校验失败 Id值非法 org:%d app:%d proj:%d table:%d user:%d",
					ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId)
				return errs.ReqParamsValidateError
			}
		}
		updateIssueVo := &UpdateIssueVo{}
		updateIssueVo.PData = make(map[string]interface{}, 0)
		updateIssueVo.LcData = make(map[string]interface{}, 0)
		updateIssueVo.RelatingChange = make(map[string]*bo.RelatingChangeBo)
		updateIssueVo.BaRelatingChange = make(map[string]*bo.RelatingChangeBo)
		updateIssueVo.SingleRelatingChange = make(map[string]*bo.RelatingChangeBo)
		for k, v := range d {
			switch k {
			case consts.BasicFieldId: // 这里前端传的id实际是issueId
				updateIssueVo.IssueId = cast.ToInt64(v)
				updateIssueVo.LcData[consts.BasicFieldIssueId] = v
				ctx.UpdateIssueIds = append(ctx.UpdateIssueIds, updateIssueVo.IssueId)
			default:
				newV, errSys := checkColumnValue(k, v)
				if errSys != nil {
					return errSys
				}
				updateIssueVo.LcData[k] = newV
				updateIssueVo.UpdateColumns = append(updateIssueVo.UpdateColumns, k)
				ctx.UpdateColumnMap[k] = struct{}{}
			}
		}
		// 起止时间校验
		if planStartTimeI, ok := updateIssueVo.LcData[consts.BasicFieldPlanStartTime]; ok {
			if planEndTimeI, ok := updateIssueVo.LcData[consts.BasicFieldPlanEndTime]; ok {
				planStartTime, err := cast.ToTimeE(planStartTimeI)
				if err != nil {
					return errs.ReqParamsValidateError
				}
				planEndTime, err := cast.ToTimeE(planEndTimeI)
				if err != nil {
					return errs.ReqParamsValidateError
				}
				if planStartTime.After(consts.BlankTimeObject) && planEndTime.After(consts.BlankTimeObject) {
					if planEndTime.Before(planStartTime) || planEndTime.Equal(planStartTime) {
						return errs.PlanEndTimeInvalidError
					}
				}
			}
		}
		ctx.UpdateIssues[updateIssueVo.IssueId] = updateIssueVo
	}

	if len(ctx.UpdateIssueIds) == 0 {
		return nil
	}

	for column, _ := range ctx.UpdateColumnMap {
		// 判断涉及字段是否允许进行批量编辑
		if errSys = ctx.checkColumnAllowBatchUpdate(column); errSys != nil {
			log.Errorf("[BatchUpdateIssue] 参数校验失败 字段不允许批量编辑 org:%d app:%d proj:%d table:%d user:%d column:%v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, column)
			return errSys
		}
		ctx.UpdateColumns = append(ctx.UpdateColumns, column)
	}

	// 获取修改前的老数据
	errSys = ctx.fetchOldDatas()
	if errSys != nil {
		log.Errorf("[BatchUpdateIssue] 获取老数据失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
			ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
		return errSys
	}

	// 判断是否有已删除的数据
	for _, updateIssue := range ctx.UpdateIssues {
		if cast.ToInt64(updateIssue.OldLcData[consts.BasicFieldRecycleFlag]) != consts.AppIsNoDelete {
			return errs.IssueAlreadyBeDeleted
		}
	}

	// 兼容没拿到AppId/ProjectId/TableId的情况
	if ctx.Req.AppId == -1 || ctx.Req.ProjectId == -1 || ctx.Req.TableId == -1 {
		if len(ctx.UpdateIssues) == 0 {
			return errs.BuildSystemErrorInfoWithMessage(errs.ParamError, fmt.Sprintf("appId: %d, projectId: %d, tableId: %d", ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId))
		}
		ctx.prepareAppProjectTableId()
	}

	// 获取org base info
	resp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: ctx.Req.OrgId})
	if resp.Failure() {
		log.Errorf("[BatchUpdateIssue] 获取Org失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
			ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
		return resp.Error()
	}
	ctx.OrgBaseInfo = resp.BaseOrgInfo

	// 获取project
	if ctx.Req.ProjectId > 0 {
		ctx.Project, errSys = domain.GetProjectSimple(ctx.Req.OrgId, ctx.Req.ProjectId)
		if errSys != nil {
			log.Errorf("[BatchUpdateIssue] 获取Project失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
			return errSys
		}
	}

	// 获取表基本信息
	if ctx.Req.TableId > 0 {
		ctx.TableMeta, errSys = domain.GetTableByTableId(ctx.Req.OrgId, ctx.Req.UserId, ctx.Req.TableId)
		if errSys != nil {
			log.Errorf("[BatchUpdateIssue] 获取表信息失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
			return errSys
		}
	}

	// 获取表头
	ctx.TableColumns, errSys = domain.GetTableColumnsMap(ctx.Req.OrgId, ctx.Req.TableId, nil, true)
	if errSys != nil {
		log.Errorf("[BatchUpdateIssue] 获取表头失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
			ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, errSys)
		return errSys
	}

	// 获取IssueStatus表头
	// 获取IssueStatus表头，如果获取不到就忽略，这个时候没有issueStatus表头
	if issueStatusColumn, ok := ctx.TableColumns[consts.BasicFieldIssueStatus]; ok {
		ctx.AllIssueStatus = domain.GetStatusListFromStatusColumn(issueStatusColumn)
		log.Infof("[BatchUpdateIssue] issueStatus 表头: %v, org:%d app:%d proj:%d table:%d",
			json.ToJsonIgnoreError(ctx.AllIssueStatus), ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId)
	} else {
		log.Infof("[BatchUpdateIssue] no need issueStatus, org:%d app:%d proj:%d table:%d", ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId)
	}

	// 查询before/after的order
	var baDataIds []int64
	if ctx.Req.BeforeDataId > 0 {
		baDataIds = append(baDataIds, ctx.Req.BeforeDataId)
	}
	if ctx.Req.AfterDataId > 0 {
		baDataIds = append(baDataIds, ctx.Req.AfterDataId)
	}
	if len(baDataIds) > 0 {
		data, errSys := domain.GetIssueInfosMapLcByDataIds(ctx.Req.OrgId, ctx.Req.UserId, baDataIds, lc_helper.ConvertToFilterColumn(consts.BasicFieldOrder))
		if errSys != nil {
			return errSys
		}
		for _, d := range data {
			dataId := cast.ToInt64(d[consts.BasicFieldId])
			order := cast.ToFloat64(d[consts.BasicFieldOrder])
			if dataId == ctx.Req.BeforeDataId {
				ctx.BeforeOrder = &order
			} else if dataId == ctx.Req.AfterDataId {
				ctx.AfterOrder = &order
			}
		}
		// 做个保护
		if ctx.BeforeOrder != nil &&
			ctx.AfterOrder != nil &&
			*ctx.BeforeOrder > *ctx.AfterOrder {
			temp := ctx.BeforeOrder
			ctx.BeforeOrder = ctx.AfterOrder
			ctx.AfterOrder = temp
		}
	}

	// 收集相关数据
	ctx.prepareData()

	// 判断 如果记录没有变化就不需要更新和推送
	ctx.issueChange()

	return nil
}

// 对比 老数据和 传参要修改的数据 如果不产生变化就不需要更新
func (ctx *BatchUpdateIssueContext) issueChange() {
	if ctx.BeforeOrder != nil || ctx.AfterOrder != nil { // 移动order的时候，其他字段一般是没有变化的
		return
	}
	for index, updateIssue := range ctx.UpdateIssues {
		lcData := updateIssue.LcData
		oldLcData := updateIssue.OldLcData
		for col, value := range lcData {
			if col == consts.BasicFieldIssueId || col == consts.BasicFieldId {
				continue
			}
			if oldData, ok := oldLcData[col]; ok {
				oldDataStr := json.ToJsonIgnoreError(oldData)
				newDataStr := json.ToJsonIgnoreError(value)
				if oldDataStr == newDataStr {
					delete(updateIssue.LcData, col)
				}
			}
		}
		// 如果传参只剩下issueId或者id，那么就说明没有需要更新的内容
		needUpdateColumn := []string{}
		for k, _ := range lcData {
			if k != consts.BasicFieldIssueId && k != consts.BasicFieldId {
				needUpdateColumn = append(needUpdateColumn, k)
			}
		}
		if len(needUpdateColumn) == 0 {
			delete(ctx.UpdateIssues, index)
		}
	}
}

// 填装AppId ProjectId TableId(外层没传入的情况下)
func (ctx *BatchUpdateIssueContext) prepareAppProjectTableId() {
	// 获取一条数据的ProjectId TableId
	for _, updateIssue := range ctx.UpdateIssues {
		if updateIssue.OldIssueBo == nil {
			log.Errorf("[prepareAppProjectTableId] invalid issue data: %v", json.ToJsonIgnoreError(updateIssue.LcData))
			break
		}
		if ctx.Req.AppId <= 0 {
			ctx.Req.AppId = updateIssue.OldIssueBo.AppId
		}
		if ctx.Req.ProjectId <= 0 {
			ctx.Req.ProjectId = updateIssue.OldIssueBo.ProjectId
		}
		if ctx.Req.TableId <= 0 {
			ctx.Req.TableId = updateIssue.OldIssueBo.TableId
		}
		break
	}
}

func (ctx *BatchUpdateIssueContext) prepareData() {
	for _, updateIssue := range ctx.UpdateIssues {
		for k, v := range updateIssue.LcData {
			// 成员 部门
			if header, ok := ctx.TableColumns[k]; ok {
				if header.Field.Type == consts.LcColumnFieldTypeMember {
					newMember := cast.ToStringSlice(v)
					newMemberIds := businees.LcMemberToUserIds(newMember)
					ctx.UserIds = append(ctx.UserIds, newMemberIds...)

					updateIssue.LcData[k] = newMember // 重新赋值一次，可以把null转成[]
				} else if header.Field.Type == consts.LcColumnFieldTypeDept {
					newDept := cast.ToStringSlice(v)
					var newDeptIds []int64
					for _, id := range newDept {
						newDeptIds = append(newDeptIds, cast.ToInt64(id))
					}
					ctx.DeptIds = append(ctx.DeptIds, newDeptIds...)

					updateIssue.LcData[k] = newDept // 重新赋值一次，可以把null转成[]
				}
			}
		}

		for k, v := range updateIssue.OldLcData {
			// 成员 部门
			if header, ok := ctx.TableColumns[k]; ok {
				if header.Field.Type == consts.LcColumnFieldTypeMember {
					oldMember := cast.ToStringSlice(v)
					oldMemberIds := businees.LcMemberToUserIds(oldMember)
					ctx.UserIds = append(ctx.UserIds, oldMemberIds...)

				} else if header.Field.Type == consts.LcColumnFieldTypeDept {
					oldDept := cast.ToStringSlice(v)
					var oldDeptIds []int64
					for _, id := range oldDept {
						oldDeptIds = append(oldDeptIds, cast.ToInt64(id))
					}
					ctx.DeptIds = append(ctx.DeptIds, oldDeptIds...)
				}
			}
		}
	}
}

func (ctx *BatchUpdateIssueContext) checkColumnAllowBatchUpdate(columnId string) errs.SystemErrorInfo {
	forbidenColumnIds := map[string]struct{}{
		consts.BasicFieldCode:       struct{}{},
		consts.BasicFieldCreator:    struct{}{},
		consts.BasicFieldCreateTime: struct{}{},
		consts.BasicFieldUpdator:    struct{}{},
		consts.BasicFieldUpdateTime: struct{}{},
	}
	forbidenColumnTypes := map[string]struct{}{
		//consts.LcColumnFieldTypeRichText:   struct{}{},
		//consts.LcColumnFieldTypeImage:      struct{}{},
		//consts.LcColumnFieldTypeDocument:   struct{}{},
		//consts.LcColumnFieldTypeRelating:   struct{}{},
		//consts.LcColumnFieldTypeBaRelating: struct{}{},
		//consts.LcColumnFieldTypeWorkHour:   struct{}{},
	}

	if _, ok := forbidenColumnIds[columnId]; ok {
		return errs.BatchUpdateForbidenColumn
	}
	if column, ok := ctx.TableColumns[columnId]; ok {
		if _, ok := forbidenColumnTypes[column.Field.Type]; ok {
			return errs.BatchUpdateForbidenColumn
		}
	}
	return nil
}

func (ctx *BatchUpdateIssueContext) checkAuth() errs.SystemErrorInfo {
	operation := consts.OperationProIssue4Modify
	orgId := ctx.Req.OrgId
	appId := ctx.Req.AppId
	projectId := ctx.Req.ProjectId
	tableId := ctx.Req.TableId
	userId := ctx.Req.UserId
	checkAuthFields := ctx.UpdateColumns
	checkAuthDataIds := ctx.UpdateDataIds

	// 待办填写任务，根据待办的权限进行校验
	if ctx.Todo != nil {
		return ctx.checkAuthTodo()
	}

	if len(checkAuthFields) == 0 || len(checkAuthDataIds) == 0 {
		return nil
	}

	log.Infof("[BatchUpdateIssue] 权限校验，org:%d app:%d proj:%d table:%d user:%d op:%s fields:%v dataIds:%v",
		orgId, appId, projectId, tableId, userId, operation, checkAuthFields, checkAuthDataIds)

	// 放开不属于任何项目的任务的权限判断
	if projectId == 0 || appId == 0 {
		return nil
	}

	// 校验项目是否归档
	if ctx.Project.IsFiling == consts.AppIsFilling {
		return errs.BuildSystemErrorInfo(errs.ProjectIsArchivedWhenModifyIssue)
	}

	appAuthResp := permissionfacade.GetAppAuthWithoutCollaborator(orgId, appId, userId)
	if appAuthResp.Failure() {
		log.Errorf("[BatchAuthIssue] GetAppAuthWithoutCollaborator, org:%d app:%d user:%d, err:%v", orgId, appId, userId, appAuthResp.Message)
		return appAuthResp.Error()
	}
	log.Infof("[BatchAuthIssue] GetAppAuthWithoutCollaborator，org:%d app:%d user:%d, result:%v", orgId, appId, userId, appAuthResp.Data)

	// 鉴权时，检查是否是系统管理员
	if appAuthResp.Data.OrgOwner || appAuthResp.Data.AppOwner ||
		appAuthResp.Data.SysAdmin || appAuthResp.Data.HasAppRootPermission {
		return nil
	}

	orgConfigResp := orgfacade.GetOrgConfig(orgvo.GetOrgConfigReq{OrgId: orgId})
	if orgConfigResp.Failure() {
		log.Errorf("[BatchAuthIssue] orgId:%d, GetOrgConfig err:%v", orgId, orgConfigResp.Message)
		return orgConfigResp.Error()
	}

	access := domain.CheckAppOptFieldAuth(orgConfigResp.Data.PayLevel, tableId, operation, checkAuthFields, appAuthResp.Data)
	if !access {
		// 通过app的权限校验失败，再通过数据的协作人权限校验一次
		dataAuthResp := permissionfacade.GetDataAuthBatch(permissionvo.GetDataAuthBatchReq{
			OrgId:   orgId,
			AppId:   appId,
			UserId:  userId,
			DataIds: checkAuthDataIds,
		})
		if dataAuthResp.Failure() {
			log.Errorf("[BatchAuthIssue] org:%d app:%d user:%d dataIds:%v, GetDataAuthBatch err: %v", orgId, appId, userId, checkAuthDataIds, dataAuthResp.Message)
			return dataAuthResp.Error()
		}

		log.Infof("[BatchAuthIssue] GetDataAuthBatch, org:%d app:%d user:%d, result:%v",
			orgId, appId, userId, json.ToJsonIgnoreError(dataAuthResp.Data))

		// 每条数据都需要单独校验权限
		for dataId, appAuth := range dataAuthResp.Data {
			// 有一条数据权限不足就返回失败
			if !domain.CheckAppOptFieldAuth(orgConfigResp.Data.PayLevel, tableId, operation, checkAuthFields, appAuth) {
				log.Infof("[BatchAuthIssue] 数据权限校验失败 org:%d app:%d user:%d dataId:%d", orgId, appId, userId, dataId)
				return errs.BuildSystemErrorInfo(errs.NoOperationPermissionForIssueUpdate)
			}
		}
	}

	return nil
}

func (ctx *BatchUpdateIssueContext) checkAuthTodo() errs.SystemErrorInfo {
	formSettings := make(map[string]*automationPb.FormSetting)

	switch v := ctx.Todo.Parameters.(type) {
	case *automationPb.Todo_TodoAudit:
		parameters := v.TodoAudit
		for _, s := range parameters.FormSettings {
			formSettings[s.ColumnId] = s
		}
	case *automationPb.Todo_TodoFillIn:
		parameters := v.TodoFillIn
		for _, s := range parameters.FormSettings {
			formSettings[s.ColumnId] = s
		}
	}
	if len(formSettings) == 0 {
		return nil
	}

	for _, updateIssue := range ctx.UpdateIssues {
		// check required
		for columnId, setting := range formSettings {
			if setting.Required {
				if _, ok := updateIssue.LcData[columnId]; !ok {
					return errs.TodoFillInMissRequiredColumn
				}
			}
		}

		// check read only
		for columnId, _ := range updateIssue.LcData {
			if setting, ok := formSettings[columnId]; ok {
				if !setting.CanWrite {
					return errs.TodoFillInChangeReadOnlyColumn
				}
			} else {
				// 未配置的列: 暂时放行
			}
		}
	}
	return nil
}

func (ctx *BatchUpdateIssueContext) process() errs.SystemErrorInfo {
	var errSys errs.SystemErrorInfo

	// 处理更新
	errSys = ctx.processUpdate()
	if errSys != nil {
		return errSys
	}

	// 处理待办
	if ctx.Todo != nil {
		ctx.processTodo()
	}

	// 保存数据
	errSys = ctx.saveToDB()
	if errSys != nil {
		return errSys
	}

	// 后续处理（异步）
	asyn.Execute(ctx.processHooks)

	return nil
}

func (ctx *BatchUpdateIssueContext) fetchOldDatas() errs.SystemErrorInfo {
	oldDatas, errSys := domain.GetIssueInfosMapLcByIssueIds(ctx.Req.OrgId, ctx.Req.UserId, ctx.UpdateIssueIds)
	if errSys != nil {
		return errSys
	}
	for _, d := range oldDatas {
		data := d
		id := cast.ToInt64(data[consts.BasicFieldIssueId])
		if updateIssue, ok := ctx.UpdateIssues[id]; ok {
			// 得到老的issueBo
			oldIssueBo, errSys := domain.ConvertIssueDataToIssueBo(data)
			if errSys != nil {
				return errSys
			}
			updateIssue.OldLcData = data
			updateIssue.OldIssueBo = oldIssueBo

			// 复制一份老的data，合并新的data，生成新的issueBo
			newData := make(map[string]interface{})
			copyer.Copy(data, &newData)
			for k, v := range updateIssue.LcData {
				newData[k] = v
			}
			newIssueBo, errSys := domain.ConvertIssueDataToIssueBo(newData)
			if errSys != nil {
				return errSys
			}
			updateIssue.NewIssueBo = newIssueBo

			ctx.UpdateDataIds = append(ctx.UpdateDataIds, oldIssueBo.DataId)
			//ctx.IterationIds = append(ctx.IterationIds, oldIssueBo.IterationId)
		} else {
			return errs.IssueNotExist
		}
	}
	return nil
}

func (ctx *BatchUpdateIssueContext) processUpdate() errs.SystemErrorInfo {
	var i int
	for _, ui := range ctx.UpdateIssues {
		updateIssue := ui
		if errSys := updateIssue.processIssueStatus(ctx); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processTitle(ctx); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processPlanTime(ctx); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processIteration(ctx); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processRemark(ctx); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processOwner(ctx); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processFollower(ctx); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processAuditor(ctx); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processRelating(ctx, tablePb.ColumnType_relating.String()); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processRelating(ctx, tablePb.ColumnType_baRelating.String()); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processRelating(ctx, tablePb.ColumnType_singleRelating.String()); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processDocument(ctx); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processImage(ctx); errSys != nil {
			return errSys
		}
		if errSys := updateIssue.processOrder(ctx, i); errSys != nil {
			return errSys
		}

		i += 1
		updateIssue.PData[consts.TcUpdator] = ctx.Req.UserId
		updateIssue.LcData[consts.BasicFieldUpdator] = ctx.Req.UserId
		updateIssue.LcData[consts.BasicFieldUpdateTime] = ctx.Now
		updateIssue.LcData[consts.BasicFieldId] = cast.ToString(updateIssue.OldIssueBo.DataId)
	}
	return nil
}

func (ctx *BatchUpdateIssueContext) processTodo() errs.SystemErrorInfo {
	result, ok := ctx.Todo.Operators[ctx.Req.UserId]
	if !ok {
		log.Errorf("[BatchUpdateIssue] invalid todo operator: %d, %v.", ctx.Req.UserId, json.ToJsonIgnoreError(result))
		return errs.TodoInvalidOperator
	}
	ctx.UpdateTodo = &commonvo.UpdateTodoReq{
		OrgId:  ctx.Req.OrgId,
		UserId: ctx.Req.UserId,
		Input: &automationPb.UpdateTodoReq{
			Id:  cast.ToString(ctx.Req.TodoId),
			Op:  automationPb.TodoOp(ctx.Req.TodoOp),
			Msg: ctx.Req.TodoMsg,
		},
	}

	if ctx.Req.TodoOp == int(automationPb.TodoOp_OpWithdraw) {
		// 撤回
		ctx.TodoResult = int(automationPb.TodoOp_OpWithdraw)
		ctx.UpdateTodo.Input.Status = automationPb.TodoStatus_SWithdrew
	} else {
		// 通过、拒绝
		switch v := ctx.Todo.Parameters.(type) {
		case *automationPb.Todo_TodoAudit:
			// 审批
			parameters := v.TodoAudit
			if !ok {
				log.Error("[BatchUpdateIssue] invalid audit todo parameters.")
				return errs.TodoInvalidParameter
			}

			switch parameters.SignMode {
			case automationPb.SignMode_Or:
				// 或签
				ctx.TodoResult = ctx.Req.TodoOp
				if ctx.Req.TodoOp == int(automationPb.TodoOp_OpPass) {
					ctx.UpdateTodo.Input.Status = automationPb.TodoStatus_SFinishedPassed
				} else {
					ctx.UpdateTodo.Input.Status = automationPb.TodoStatus_SFinishedRejected
				}

			case automationPb.SignMode_AndOnePass:
				// 会签: 一人通过+全部驳回
				if ctx.Req.TodoOp == int(automationPb.TodoOp_OpPass) {
					ctx.TodoResult = int(automationPb.TodoOp_OpPass)
					ctx.UpdateTodo.Input.Status = automationPb.TodoStatus_SFinishedPassed
				} else {
					isAllReject := true
					for _, result := range ctx.Todo.Operators {
						if result.Op != automationPb.TodoOp_OpReject {
							isAllReject = false
							break
						}
					}
					if isAllReject {
						ctx.TodoResult = int(automationPb.TodoOp_OpReject)
						ctx.UpdateTodo.Input.Status = automationPb.TodoStatus_SFinishedRejected
					}
				}

			case automationPb.SignMode_AndAllPass:
				// 会签: 全部通过+一人驳回
				if ctx.Req.TodoOp == int(automationPb.TodoOp_OpPass) {
					isAllPass := true
					for _, result := range ctx.Todo.Operators {
						if result.Op != automationPb.TodoOp_OpPass {
							isAllPass = false
							break
						}
					}
					if isAllPass {
						ctx.TodoResult = int(automationPb.TodoOp_OpPass)
						ctx.UpdateTodo.Input.Status = automationPb.TodoStatus_SFinishedPassed
					}
				} else {
					ctx.TodoResult = int(automationPb.TodoOp_OpReject)
					ctx.UpdateTodo.Input.Status = automationPb.TodoStatus_SFinishedRejected
				}

			default:
				log.Errorf("[BatchUpdateIssue] invalid sign mode: %v", parameters.SignMode)
				return errs.TodoInvalidParameter
			}

		case *automationPb.Todo_TodoFillIn:
			// 填写
			ctx.TodoResult = ctx.Req.TodoOp
			ctx.UpdateTodo.Input.Status = automationPb.TodoStatus_SFinishedPassed

		default:
			log.Errorf("[BatchUpdateIssue] invalid todo type: %v", ctx.Todo.Type)
			return errs.TodoInvalidParameter
		}
	}

	return nil
}

func (ctx *BatchUpdateIssueContext) saveToDB() errs.SystemErrorInfo {
	// 组装需要创建的issue relations id
	if len(ctx.CreateIssueRelations) > 0 {
		relationIds, errSys := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, len(ctx.CreateIssueRelations))
		if errSys != nil {
			return errs.BuildSystemErrorInfo(errs.ApplyIdError, errSys)
		}
		for i := 0; i < len(ctx.CreateIssueRelations); i++ {
			ctx.CreateIssueRelations[i].Id = relationIds.Ids[i].Id
		}
	}

	// 组装需要放入回收站的附件资源
	if len(ctx.CreateRecycleBins) > 0 {
		relationIds, errSys := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableRecycleBin, len(ctx.CreateRecycleBins))
		if errSys != nil {
			return errs.BuildSystemErrorInfo(errs.ApplyIdError, errSys)
		}
		versionId, errSys := idfacade.ApplyPrimaryIdRelaxed(consts.RecycleVersion)
		if errSys != nil {
			return errs.BuildSystemErrorInfo(errs.ApplyIdError, errSys)
		}
		for i := 0; i < len(ctx.CreateRecycleBins); i++ {
			ctx.CreateRecycleBins[i].Id = relationIds.Ids[i].Id
			ctx.CreateRecycleBins[i].Version = int(versionId)
		}
	}

	var lcReqRelating *formvo.LessUpdateIssueReq
	var relatingIssueIds []int64
	ctx.AllRelatingIssueIds = slice.SliceUniqueInt64(ctx.AllRelatingIssueIds)
	for _, id := range ctx.AllRelatingIssueIds {
		if _, ok := ctx.UpdateIssues[id]; !ok {
			relatingIssueIds = append(relatingIssueIds, id)
		}
	}
	if len(relatingIssueIds) > 0 {
		lcReqRelating = &formvo.LessUpdateIssueReq{
			OrgId: ctx.Req.OrgId,
			AppId: ctx.Req.AppId,
			//TableId: ctx.Req.TableId,
			UserId: ctx.Req.UserId,
		}
		for _, id := range relatingIssueIds {
			lcReqRelating.Form = append(lcReqRelating.Form, map[string]interface{}{
				consts.BasicFieldIssueId:    id,
				consts.BasicFieldUpdateTime: ctx.Now,
			})
		}
	}

	lcReq := formvo.LessUpdateIssueReq{
		OrgId:   ctx.Req.OrgId,
		AppId:   ctx.Req.AppId,
		TableId: ctx.Req.TableId,
		UserId:  ctx.Req.UserId,
	}
	for _, ui := range ctx.UpdateIssues {
		updateIssue := ui
		lcReq.Form = append(lcReq.Form, updateIssue.LcData)
	}

	sysErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		fs := make([]batch.TxExecFunc, 0)

		// 更新issue
		//fs = append(fs, func(tx sqlbuilder.Tx) errors.SystemErrorInfo {
		//	for _, ui := range ctx.UpdateIssues {
		//		updateIssue := ui
		//
		//		log.Infof("[BatchUpdateIssue] update issues: %v", json.ToJsonIgnoreError(updateIssue.PData))
		//		err := mysql.TransUpdateSmart(tx, consts.TableIssue, updateIssue.IssueId, updateIssue.PData)
		//		if err != nil {
		//			log.Errorf("[BatchUpdateIssue] mysql.TransUpdateSmart: %v", err)
		//			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		//		}
		//	}
		//	return nil
		//})

		// 更新issue relations
		fs = append(fs, func(tx sqlbuilder.Tx) errors.SystemErrorInfo {
			for _, s := range ctx.SqlIssueRelations {
				sql := s
				log.Infof("[BatchUpdateIssue] tx.Exec, sql: %s, args: %v", sql.Query, json.ToJsonIgnoreError(sql.Args))
				_, err := tx.Exec(sql.Query, sql.Args...)
				if err != nil {
					log.Errorf("[BatchUpdateIssue] tx.Exec: %v, sql: %s, args: %v", err, sql.Query, sql.Args)
					return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
				}
			}

			// 创建新issue relations
			if len(ctx.CreateIssueRelations) > 0 {
				log.Infof("[BatchUpdateIssue] insert issue relations: %v", json.ToJsonIgnoreError(ctx.CreateIssueRelations))
				err := mysql.TransBatchInsert(tx, &po.PpmPriIssueRelation{}, slice.ToSlice(ctx.CreateIssueRelations))
				if err != nil {
					log.Errorf("[BatchUpdateIssue] mysql.TransBatchInsert: %v", err)
					return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
				}
			}
			return nil
		})

		fs = append(fs, func(tx sqlbuilder.Tx) errors.SystemErrorInfo {
			// 插入resource relations
			for issueId, v := range ctx.CreateAttachments {
				if len(v.ResourceIds) > 0 {
					errSys := domain.AddResourceRelation(ctx.Req.OrgId, ctx.Req.UserId, ctx.Req.ProjectId, issueId, v.ResourceIds, consts.OssPolicyTypeLesscodeResource, v.ColumnId)
					if errSys != nil {
						log.Errorf("[BatchUpdateIssue] AddResourceRelation err:%v", errSys)
						return errSys
					}
				}
			}
			return nil
		})

		fs = append(fs, func(tx sqlbuilder.Tx) errors.SystemErrorInfo {
			// 删除附件到回收站的数据
			if len(ctx.CreateRecycleBins) > 0 {
				log.Infof("[BatchUpdateIssue] insert attachment recycle: %v", json.ToJsonIgnoreError(ctx.CreateRecycleBins))
				err := mysql.TransBatchInsert(tx, &po.PpmPrsRecycleBin{}, slice.ToSlice(ctx.CreateRecycleBins))
				if err != nil {
					log.Errorf("[BatchUpdateIssue] mysql.TransBatchInsert: %v", err)
					return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
				}
			}
			for issueId, v := range ctx.RecycleAttachments {
				if len(v.ResourceIds) > 0 {
					errSys := domain.RecycleIssueResources(ctx.Req.OrgId, ctx.Req.UserId, ctx.Req.ProjectId, issueId, v.ResourceIds,
						int64(ctx.CreateRecycleBins[0].Version), v.ColumnId)
					if errSys != nil {
						log.Errorf("[BatchUpdateIssue] RecycleIssueResources error:%v", errSys)
						return errSys
					}
				}
			}
			return nil
		})

		fs = append(fs, func(tx sqlbuilder.Tx) errors.SystemErrorInfo {
			// 删除被移出的确认人的确认状态
			for issueId, deleteAuditors := range ctx.DeleteAuditors {
				domain.DeleteLcIssueAuditStatusDetailByUsers(ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.UserId, issueId, deleteAuditors...)
			}
			return nil
		})

		txBatchExecutor := &batch.TxBatchExecutor{}
		txBatchExecutor.Init(tx, 20)
		txBatchExecutor.PushJobs(fs)
		errSys := txBatchExecutor.StartAndWaitFinish()
		if errSys != nil {
			log.Errorf("[BatchUpdateIssue] txBatchExecutor.StartAndWaitFinish: %v", errSys)
			return errSys
		}

		// 刷新关联的任务updateTimes
		if lcReqRelating != nil {
			log.Infof("[BatchUpdateIssue] LessUpdateIssue, relating update time, req: %v", json.ToJsonIgnoreError(lcReqRelating))
			resp := formfacade.LessUpdateIssue(*lcReqRelating)
			if resp.Failure() {
				log.Errorf("[BatchUpdateIssue] LessUpdateIssue, relating update time failed: %v, req: %v", resp.Error(), json.ToJsonIgnoreError(lcReqRelating))
				return resp.Error()
			}
		}

		// 事务的最后更新无码
		if len(lcReq.Form) > 0 {
			log.Infof("[BatchUpdateIssue] LessUpdateIssue, req: %v", json.ToJsonIgnoreError(lcReq))
			resp := formfacade.LessUpdateIssue(lcReq)
			if resp.Failure() {
				log.Errorf("[BatchUpdateIssue] LessUpdateIssue: %v, req: %v", resp.Error(), json.ToJsonIgnoreError(lcReq))
				return resp.Error()
			}
		}

		return nil
	})
	if sysErr != nil {
		log.Errorf("[BatchUpdateIssue] saveToDB: %v", sysErr)
		// 如果客户端频繁保存，会出现这个错误
		// Error 1213: Deadlock found when trying to get lock; try restarting transaction
		if strings.Contains(sysErr.Error(), consts.MySQL_DEADLOCK_ERROR) {
			return errs.RequestFrequentError
		}
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, sysErr)
	}
	return nil
}

func (ctx *BatchUpdateIssueContext) processHooks() {
	ctx.IterationIds = slice.SliceUniqueInt64(ctx.IterationIds)
	ctx.UserIds = slice.SliceUniqueInt64(ctx.UserIds)
	ctx.DeptIds = slice.SliceUniqueInt64(ctx.DeptIds)
	ctx.RelatingIssueIds = slice.SliceUniqueInt64(ctx.RelatingIssueIds)

	if len(ctx.IterationIds) > 0 {
		iterations, errSys := domain.GetIterationsInfo(ctx.Req.OrgId, ctx.IterationIds)
		if errSys != nil {
			log.Errorf("[BatchUpdateIssue] 获取迭代失败 org:%d app:%d proj:%d table:%d user:%d, iterations:%v err: %v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, ctx.IterationIds, errSys)
		} else {
			for _, it := range iterations {
				i := it
				ctx.Iterations[i.Id] = &i
			}
		}
	}
	if len(ctx.UserIds) > 0 {
		resp := userfacade.GetAllUserByIds(ctx.Req.OrgId, ctx.UserIds)
		if resp.Failure() {
			log.Errorf("[BatchUpdateIssue] 获取用户失败 org:%d app:%d proj:%d table:%d user:%d, user:%v err: %v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, ctx.UserIds, resp.Error())
		} else {
			for _, user := range resp.Data {
				u := user
				ctx.Users[u.Id] = &u
			}
		}
	}
	if len(ctx.DeptIds) > 0 {
		resp := userfacade.GetAllDeptByIds(ctx.Req.OrgId, ctx.DeptIds)
		if resp.Failure() {
			log.Errorf("[BatchUpdateIssue] 获取部门失败 org:%d app:%d proj:%d table:%d user:%d, dept:%v err: %v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, ctx.DeptIds, resp.Error())
		} else {
			for _, dept := range resp.Data {
				d := dept
				ctx.Depts[d.Id] = &d
			}
		}
	}
	if len(ctx.RelatingIssueIds) > 0 {
		data, errSys := domain.GetIssueInfosMapLcByIssueIds(ctx.Req.OrgId, ctx.Req.UserId, ctx.RelatingIssueIds)
		if errSys != nil {
			log.Errorf("[BatchUpdateIssue] 相关关联前后置任务数据失败 org:%d app:%d proj:%d table:%d user:%d, ids:%v err: %v",
				ctx.Req.OrgId, ctx.Req.AppId, ctx.Req.ProjectId, ctx.Req.TableId, ctx.Req.UserId, ctx.RelatingIssueIds, errSys)
		} else {
			for _, d := range data {
				issueBo, errSys := domain.ConvertIssueDataToIssueBo(d)
				if errSys == nil {
					ctx.RelatingIssues[issueBo.Id] = issueBo
				}
			}
		}
	}

	for _, ui := range ctx.UpdateIssues {
		updateIssue := ui
		// 生成剩余的动态
		updateIssue.processTrends(ctx)
	}

	asyn.Execute(func() {
		for _, ui := range ctx.UpdateIssues {
			updateIssue := ui

			if ctx.TriggerBy.TriggerBy != consts.TriggerByAutoSchedule {
				// 飞书卡片通知确认人
				if updateIssue.NoticeAuditorIssueStatusChange {
					domain.NoticeIssueAudit(ctx.Req.OrgId, ctx.Req.UserId, ctx.OrgBaseInfo, updateIssue.NewIssueBo, nil)
				} else if updateIssue.NoticeAuditorChange {
					domain.NoticeIssueAudit(ctx.Req.OrgId, ctx.Req.UserId, ctx.OrgBaseInfo, updateIssue.NewIssueBo, updateIssue.OldIssueBo.AuditorIdsI64)
				}
			}

			// 工时相关处理
			if updateIssue.UpdateOwner {
				var oldOwnerId, newOwnerId int64
				if updateIssue.OldIssueBo.OwnerIdI64 != nil && len(updateIssue.OldIssueBo.OwnerIdI64) > 0 {
					oldOwnerId = updateIssue.OldIssueBo.OwnerIdI64[0]
				}
				if updateIssue.NewIssueBo.OwnerIdI64 != nil && len(updateIssue.NewIssueBo.OwnerIdI64) > 0 {
					newOwnerId = updateIssue.NewIssueBo.OwnerIdI64[0]
				}
				err := work_hour.ChangeWorkerIdWhenChangedIssueOwner(ctx.Req.OrgId, updateIssue.IssueId, oldOwnerId, newOwnerId)
				if err != nil {
					log.Errorf("[BatchUpdateIssue] ChangeWorkerIdWhenChangedIssueOwner。oldOwnerId: %v, newOwnerId: %v, err: %v", oldOwnerId, newOwnerId, err)
				}
			}
			if _, ok := ctx.UpdateColumnMap[consts.BasicFieldPlanStartTime]; ok {
				work_hour.TriggerWorkHourWhenChangedPlanStartOrEndTime(ctx.Req.OrgId, updateIssue.IssueId, consts.TcStartTime, updateIssue.OldIssueBo.PlanStartTime, updateIssue.NewIssueBo.PlanStartTime)
			}
			if _, ok := ctx.UpdateColumnMap[consts.BasicFieldPlanEndTime]; ok {
				work_hour.TriggerWorkHourWhenChangedPlanStartOrEndTime(ctx.Req.OrgId, updateIssue.IssueId, consts.TcEndTime, updateIssue.OldIssueBo.PlanEndTime, updateIssue.NewIssueBo.PlanEndTime)
			}
		}
	})

	// 动态相关
	asyn.Execute(func() {
		for _, t := range ctx.IssueTrends {
			trend := t

			// 保存动态
			domain.PushIssueTrends(trend)

			if ctx.TriggerBy.TriggerBy != consts.TriggerByAutoSchedule {
				// 推送动态，如：飞书卡片（个人机器人）、钉钉卡片等
				domain.PushIssueThirdPlatformNotice(trend)
			}
		}
	})

	// 自动排期相关
	asyn.Execute(func() {
		if len(ctx.AutoScheduleSourceIssues) > 0 {
			ctx.processAutoSchedule()
		}
	})

	// 上报事件
	asyn.Execute(func() {
		userDepts := domain.AssembleUserDepts(ctx.Users, ctx.Depts)
		for _, ui := range ctx.UpdateIssues {
			ui.reportEvent(ctx, userDepts)
		}
	})

	// 钉钉酷应用顶部卡片更新
	updateDingTopCard(ctx.Req.OrgId, ctx.Req.ProjectId)
}

func (ctx *BatchUpdateIssueContext) processAutoSchedule() {
	// 拿出本表格内所有的前后置关联任务
	condition := &tablePb.Condition{Type: tablePb.ConditionType_and, Conditions: domain.GetNoRecycleCondition(
		domain.GetRowsCondition(consts.BasicFieldTableId, tablePb.ConditionType_equal, cast.ToString(ctx.Req.TableId), nil),
		domain.GetRowsCondition(consts.BasicFieldBaRelating, tablePb.ConditionType_not_null, nil, nil),
		&tablePb.Condition{
			Type: tablePb.ConditionType_or,
			Conditions: []*tablePb.Condition{
				domain.GetRowsCondition(fmt.Sprintf("jsonb_array_length(data->'%s'->'%s')", consts.BasicFieldBaRelating, consts.BasicFieldRelatingLinkTo),
					tablePb.ConditionType_gt, 0, nil),
				domain.GetRowsCondition(fmt.Sprintf("jsonb_array_length(data->'%s'->'%s')", consts.BasicFieldBaRelating, consts.BasicFieldRelatingLinkFrom),
					tablePb.ConditionType_gt, 0, nil),
			},
		},
	)}

	columns := []string{
		lc_helper.ConvertToFilterColumn(consts.BasicFieldPlanStartTime),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldPlanEndTime),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldBaRelating),
	}
	datas, errSys := domain.GetIssueInfosMapLc(ctx.Req.OrgId, ctx.Req.UserId, condition, columns, -1, -1)
	if errSys != nil {
		log.Errorf("[AutoSchedule] GetIssueInfosMapLc err: %s", errSys)
		return
	}
	log.Infof("[AutoSchedule] all baRelatPlanWorkHouring data: %s", json.ToJsonIgnoreError(datas))

	allIssueMap := make(map[int64]*AutoScheduleIssueVo)
	for _, data := range datas {
		issue := &AutoScheduleIssueVo{}
		issue.IssueId = cast.ToInt64(data[consts.BasicFieldIssueId])
		issue.Id = cast.ToString(data[consts.BasicFieldId])

		var planStartTime, planEndTime types.Time
		planStartTime = types.Time(cast.ToTime(data[consts.BasicFieldPlanStartTime]))
		if planStartTime.IsNull() {
			planStartTime = types.Time(consts.BlankTimeObject)
		}
		planEndTime = types.Time(cast.ToTime(data[consts.BasicFieldPlanEndTime]))
		if planEndTime.IsNull() {
			planEndTime = types.Time(consts.BlankTimeObject)
		}
		issue.PlanStartTime = planStartTime
		issue.PlanEndTime = planEndTime

		if err := copyer.Copy(data[consts.BasicFieldBaRelating], &issue.BaRelating); err != nil {
			log.Errorf("[AutoSchedule] issue: %d, decode baRelating err: %v", issue.IssueId, err)
			continue
		}

		allIssueMap[issue.IssueId] = issue
	}

	// 构造链表
	for _, issue := range allIssueMap {
		for _, id := range issue.BaRelating[consts.BasicFieldRelatingLinkTo] {
			// 如果没找到，可能是跨表的前后置，也可能是脏数据，总之忽略
			if i, ok := allIssueMap[cast.ToInt64(id)]; ok {
				issue.LinkTo = append(issue.LinkTo, i)
			}
		}
		for _, id := range issue.BaRelating[consts.BasicFieldRelatingLinkFrom] {
			// 如果没找到，可能是跨表的前后置，也可能是脏数据，总之忽略
			if i, ok := allIssueMap[cast.ToInt64(id)]; ok {
				issue.LinkFrom = append(issue.LinkFrom, i)
			}
		}
	}

	// 手动更新的任务都不再参与自动排期，链路执行到这些任务将自动中断
	manualSetIdMap := make(map[int64]struct{})
	for _, i := range ctx.AutoScheduleSourceIssues {
		manualSetIdMap[i.NewIssueBo.Id] = struct{}{}
	}

	// 执行批量自动排期
	for _, i := range ctx.AutoScheduleSourceIssues {
		if issue, ok := allIssueMap[i.NewIssueBo.Id]; ok {
			processedIdMap := make(map[int64]struct{})

			// 对后置任务执行自动排期
			for _, is := range issue.LinkFrom {
				ctx.processAutoScheduleIssue(is, i.DeltaPlanEndTime, manualSetIdMap, processedIdMap)
			}
		}
	}

	form := make([]map[string]interface{}, 0)
	for _, is := range allIssueMap {
		data := make(map[string]interface{})
		if is.PlanStartTime.IsNotNull() && is.PlanStartTimeDelta != 0 {
			data[consts.BasicFieldPlanStartTime] = date.Format(time.Time(is.PlanStartTime).Add(is.PlanStartTimeDelta))
		}
		if is.PlanEndTime.IsNotNull() && is.PlanEndTimeDelta != 0 {
			data[consts.BasicFieldPlanEndTime] = date.Format(time.Time(is.PlanEndTime).Add(is.PlanEndTimeDelta))
		}
		if len(data) > 0 {
			data[consts.BasicFieldId] = is.IssueId
			form = append(form, data)
		}
	}
	if len(form) > 0 {
		log.Infof("[AutoSchedule] BatchUpdateIssue req form: %s", json.ToJsonIgnoreError(form))
		err := BatchUpdateIssue(&projectvo.BatchUpdateIssueReqInnerVo{
			OrgId:     ctx.Req.OrgId,
			UserId:    ctx.Req.UserId,
			AppId:     ctx.Req.AppId,
			ProjectId: ctx.Req.ProjectId,
			TableId:   ctx.Req.TableId,
			Data:      form,
		}, true, &projectvo.TriggerBy{
			TriggerBy: consts.TriggerByAutoSchedule,
		})
		if err != nil {
			log.Errorf("[AutoSchedule] BatchUpdateIssue failed: %s", err)
		}
	}
}

func (ctx *BatchUpdateIssueContext) processAutoScheduleIssue(issue *AutoScheduleIssueVo, delta time.Duration, manualSetIdMap, processedIdMap map[int64]struct{}) {
	// 打断重复循环处理
	if _, ok := processedIdMap[issue.IssueId]; ok {
		return
	}
	processedIdMap[issue.IssueId] = struct{}{}

	// 以手动编辑的结果为准，非手动编辑的任务才进行自动排期
	if _, ok := manualSetIdMap[issue.IssueId]; !ok {
		if issue.PlanStartTime.IsNotNull() {
			issue.PlanStartTimeDelta += delta
		}
		if issue.PlanEndTime.IsNotNull() {
			issue.PlanEndTimeDelta += delta
		}
	}

	// 对后置任务执行自动排期
	for _, is := range issue.LinkFrom {
		ctx.processAutoScheduleIssue(is, delta, manualSetIdMap, processedIdMap)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (i *UpdateIssueVo) assignUpdateData(pData, lcData map[string]interface{}) {
	for k, v := range pData {
		i.PData[k] = v

		switch k {
		case consts.TcStatus:
			i.NewIssueBo.Status = cast.ToInt64(v)
		case consts.TcAuditStatus:
			i.NewIssueBo.AuditStatus = cast.ToInt(v)
		case consts.TcUpdator:
			i.NewIssueBo.Updator = cast.ToInt64(v)
		case consts.TcStartTime:
			i.NewIssueBo.StartTime = types.Time(cast.ToTime(v))
		case consts.TcEndTime:
			i.NewIssueBo.EndTime = types.Time(cast.ToTime(v))
		}
	}
	for k, v := range lcData {
		i.LcData[k] = v
		i.NewIssueBo.LessData[k] = v
	}
}

func (i *UpdateIssueVo) processIssueStatus(ctx *BatchUpdateIssueContext) errs.SystemErrorInfo {
	var oldIssueStatus *status.StatusInfoBo
	var newIssueStatus *status.StatusInfoBo
	var isChangeAuditors bool
	var finalAuditorsCount int
	var err error

	isNormalProject := ctx.Project != nil && ctx.Project.ProjectTypeId == consts.ProjectTypeNormalId

	finalAuditorsCount = len(i.NewIssueBo.AuditorIdsI64)
	if _, ok := i.LcData[consts.BasicFieldAuditorIds]; ok {
		isChangeAuditors = true
	}

	// 匹配旧的status
	for _, is := range ctx.AllIssueStatus {
		if is.ID == i.OldIssueBo.Status {
			oldIssueStatus = &is
			break
		}
	}
	if oldIssueStatus == nil {
		log.Errorf("[BatchUpdateIssue] processIssueStatus issueId: %d old issue status not exists: %v", i.IssueId, i.OldIssueBo.Status)
		//return errs.IssueStatusNotExist
		// 这里做下兼容，允许修复异常数据
		oldIssueStatus = &status.StatusInfoBo{
			ID:   i.OldIssueBo.Status,
			Type: 0,
		}
	}

	var newIssueStatusId int64
	if value, ok := i.LcData[consts.BasicFieldIssueStatus]; ok {
		newIssueStatusId, err = cast.ToInt64E(value)
		if err != nil {
			return errs.ReqParamsValidateError
		}
	} else if value, ok = i.LcData[consts.BasicFieldIssueStatusType]; ok {
		// 支持通过修改issueStatusType来修改issueStatus
		newIssueStatusTypeId, err := cast.ToIntE(value)
		if err != nil {
			return errs.ReqParamsValidateError
		}
		for _, is := range ctx.AllIssueStatus {
			if is.Type == newIssueStatusTypeId {
				newIssueStatusId = is.ID
				break
			}
		}
	}
	var newAuditorStatusId int64
	if value, ok := i.LcData[consts.BasicFieldAuditStatus]; ok {
		newAuditorStatusId, err = cast.ToInt64E(value)
		if err != nil {
			return errs.ReqParamsValidateError
		}
	}

	// 匹配新的status
	if newIssueStatusId != 0 {
		for _, is := range ctx.AllIssueStatus {
			if is.ID == newIssueStatusId {
				newIssueStatus = &is
				break
			}
		}
		if newIssueStatus == nil {
			log.Errorf("[BatchUpdateIssue] processIssueStatus issueId: %d new issue status not exists: %v", i.IssueId, newIssueStatusId)
			return errs.IssueStatusNotExist
		}
		log.Infof("[BatchUpdateIssue] issueId %d, new issueStatus: %v %v", i.IssueId, newIssueStatusId, json.ToJsonIgnoreError(newIssueStatus))
	}

	// 老通用项目处理待确认逻辑: 如果是已完成状态下变化了确认人，也可能引起状态需要自动变化，这里简化下处理，覆盖掉用户手动的修改（如果也传了issueStatus, auditStatus过来则会被这里的处理覆盖）
	if isNormalProject && isChangeAuditors {
		// 任务修改前后都是已完成状态
		if oldIssueStatus.Type == consts.StatusTypeComplete && (newIssueStatus == nil || newIssueStatus.Type == consts.StatusTypeComplete) {
			if finalAuditorsCount == 0 {
				// 如果确认人为空，则更新审核状态为通过
				newAuditorStatusId = consts.AuditStatusPass
			} else if i.OldIssueBo.AuditStatus != consts.AuditStatusPass && i.OldIssueBo.AuditStatus != consts.AuditStatusReject {
				// 如果是待确认状态，且确认人只剩下驳回的，则更新状态为未完成(进行中)
				isAllReject := true
				if i.OldIssueBo.AuditStatusDetail == nil {
					isAllReject = false
				} else {
					for _, auditorId := range i.NewIssueBo.AuditorIdsI64 {
						auditStatus, ok := i.OldIssueBo.AuditStatusDetail[cast.ToString(auditorId)]
						if !ok || auditStatus != consts.AuditStatusReject {
							isAllReject = false
							break
						}
					}
				}

				if isAllReject {
					newAuditorStatusId = consts.AuditStatusReject
					for _, is := range ctx.AllIssueStatus {
						if is.Type == consts.StatusTypeRunning {
							newIssueStatusId = is.ID
							newIssueStatus = &is
							break
						}
					}
				}
			}
		}
	}

	// 没有需要更新的状态
	if newIssueStatusId == 0 && newAuditorStatusId == 0 {
		return nil
	}

	// 再匹配一次新status(因为newIssueStatusId可能发生了变化，或者从无到有)
	if newIssueStatusId != 0 {
		for _, is := range ctx.AllIssueStatus {
			if is.ID == newIssueStatusId {
				newIssueStatus = &is
				break
			}
		}
		if newIssueStatus == nil {
			log.Errorf("[BatchUpdateIssue] processIssueStatus issueId: %d new issue status not exists: %v", i.IssueId, newIssueStatusId)
			return errs.IssueStatusNotExist
		}
		log.Infof("[BatchUpdateIssue] issueId %d, new issueStatus: %v %v", i.IssueId, newIssueStatusId, json.ToJsonIgnoreError(newIssueStatus))
	}

	// 需更新的数据
	pData := make(map[string]interface{}, 0)
	if newIssueStatusId != 0 {
		pData[consts.TcStatus] = newIssueStatusId
		if newIssueStatus.Type == consts.StatusTypeComplete {
			pData[consts.TcEndTime] = time.Now().Format(consts.AppTimeFormat)
		} else if newIssueStatus.Type == consts.StatusTypeNotStart {
			pData[consts.TcStartTime] = consts.BlankTime
			pData[consts.TcEndTime] = consts.BlankTime
		} else if newIssueStatus.Type == consts.StatusTypeRunning {
			pData[consts.TcStartTime] = time.Now().Format(consts.AppTimeFormat)
			pData[consts.TcEndTime] = consts.BlankTime
		}
	}
	if newAuditorStatusId != 0 {
		pData[consts.TcAuditStatus] = newAuditorStatusId
	}

	// 老通用项目处理待确认逻辑
	if isNormalProject && newIssueStatus != nil {
		// 未完成->已完成 如果没有确认人，就更新任务待确认状态为确认通过
		if oldIssueStatus.Type != consts.StatusTypeComplete && newIssueStatus.Type == consts.StatusTypeComplete {
			if finalAuditorsCount == 0 {
				pData[consts.TcAuditStatus] = consts.AuditStatusPass
			}
		} else if oldIssueStatus.Type == consts.StatusTypeComplete && newIssueStatus.Type != consts.StatusTypeComplete {
			// 已完成->未完成 重置下待确认状态
			pData[consts.TcAuditStatus] = consts.AuditStatusNotView
		}
	}

	// pData 生成 lcData
	lcData := slice2.CaseCamelCopy(pData)
	delete(lcData, consts.BasicFieldStatus)

	// 老通用项目处理待确认逻辑
	if isNormalProject && newIssueStatus != nil {
		// 如果是从已完成变为未完成/未完成到已完成 都 需要将任务的确认人都置为初始状态
		if oldIssueStatus.Type != consts.StatusTypeComplete && newIssueStatus.Type == consts.StatusTypeComplete {
			lcData[consts.BasicFieldAuditStatusDetail] = map[string]int{}

			if finalAuditorsCount > 0 {
				i.NoticeAuditorIssueStatusChange = true
				// 保存一下审批发起人
				lcData[consts.BasicFieldAuditStarter] = ctx.Req.UserId
			}
		} else if oldIssueStatus.Type == consts.StatusTypeComplete && newIssueStatus.Type != consts.StatusTypeComplete {
			lcData[consts.BasicFieldAuditStatusDetail] = map[string]int{}
		}
	}

	if newIssueStatus != nil {
		lcData[consts.BasicFieldIssueStatus] = newIssueStatus.ID
		lcData[consts.BasicFieldIssueStatusType] = newIssueStatus.Type

		if oldIssueStatus.Type != newIssueStatus.Type {
			i.UpdateColumns = append(i.UpdateColumns, consts.BasicFieldIssueStatusType)
		}
	}
	i.assignUpdateData(pData, lcData)

	return nil
}

func (i *UpdateIssueVo) processTitle(ctx *BatchUpdateIssueContext) errs.SystemErrorInfo {
	if value, ok := i.LcData[consts.BasicFieldTitle]; ok {
		i.PData[consts.TcTitle] = value
		i.TrendChangeList = append(i.TrendChangeList, bo.TrendChangeListBo{
			Field:     consts.BasicFieldTitle,
			FieldName: consts.Title,
			FieldType: consts.LcColumnFieldTypeInput,
			OldValue:  i.OldIssueBo.Title,
			NewValue:  i.NewIssueBo.Title,
		})
		i.UpdateTitle = true
		ctx.HandledTrendColumnIds[consts.BasicFieldTitle] = struct{}{}
	}
	return nil
}

func (i *UpdateIssueVo) processPlanTime(ctx *BatchUpdateIssueContext) errs.SystemErrorInfo {
	if _, ok := i.LcData[consts.BasicFieldPlanEndTime]; ok {
		i.PData[consts.TcPlanEndTime] = date.Format(time.Time(i.NewIssueBo.PlanEndTime))
		i.TrendChangeList = append(i.TrendChangeList, bo.TrendChangeListBo{
			Field:     consts.BasicFieldPlanEndTime,
			FieldName: consts.PlanEndTime,
			FieldType: consts.LcColumnFieldTypeDatepicker,
			OldValue:  businees.GetTimeString(i.OldIssueBo.PlanEndTime),
			NewValue:  businees.GetTimeString(i.NewIssueBo.PlanEndTime),
		})
		ctx.HandledTrendColumnIds[consts.BasicFieldPlanEndTime] = struct{}{}

		// 自动排期相关 (不归属于项目的任务不触发，理论上一次批量编辑的任务必须都在同一张表里，或者都不归属于项目)
		// 自动排期不再次触发自动排期
		if ctx.TriggerBy.TriggerBy != consts.TriggerByAutoSchedule &&
			ctx.Req.TableId > 0 && ctx.TableMeta.AutoScheduleFlag == consts.SwitchOn {
			if i.OldIssueBo.PlanEndTime.IsNotNull() && i.NewIssueBo.PlanEndTime.IsNotNull() {
				newTime := time.Time(i.NewIssueBo.PlanEndTime)
				oldTime := time.Time(i.OldIssueBo.PlanEndTime)
				i.DeltaPlanEndTime = newTime.Sub(oldTime)
				ctx.AutoScheduleSourceIssues = append(ctx.AutoScheduleSourceIssues, i)
			}
		}
	}
	if _, ok := i.LcData[consts.BasicFieldPlanStartTime]; ok {
		i.PData[consts.TcPlanStartTime] = date.Format(time.Time(i.NewIssueBo.PlanStartTime))
		i.TrendChangeList = append(i.TrendChangeList, bo.TrendChangeListBo{
			Field:     consts.BasicFieldPlanStartTime,
			FieldName: consts.PlanStartTime,
			FieldType: consts.LcColumnFieldTypeDatepicker,
			OldValue:  businees.GetTimeString(i.OldIssueBo.PlanStartTime),
			NewValue:  businees.GetTimeString(i.NewIssueBo.PlanStartTime),
		})
		ctx.HandledTrendColumnIds[consts.BasicFieldPlanStartTime] = struct{}{}
	}

	// 更新时间校验
	planStartTime := i.NewIssueBo.PlanStartTime
	planEndTime := i.NewIssueBo.PlanEndTime
	if planStartTime.IsNotNull() && planEndTime.IsNotNull() {
		if time.Time(planEndTime).Before(time.Time(planStartTime)) {
			return errs.PlanEndTimeInvalidError
		}
	}
	return nil
}

func (i *UpdateIssueVo) processIteration(ctx *BatchUpdateIssueContext) errs.SystemErrorInfo {
	if _, ok := i.LcData[consts.BasicFieldIterationId]; ok {
		if i.NewIssueBo.IterationId == i.OldIssueBo.IterationId {
			return nil
		}

		i.PData[consts.TcIterationId] = i.NewIssueBo.IterationId
		ctx.IterationIds = append(ctx.IterationIds, i.OldIssueBo.IterationId, i.NewIssueBo.IterationId)
	}
	return nil
}

func (i *UpdateIssueVo) processRemark(ctx *BatchUpdateIssueContext) errs.SystemErrorInfo {
	if _, ok := i.LcData[consts.BasicFieldRemark]; ok {
		i.TrendChangeList = append(i.TrendChangeList, bo.TrendChangeListBo{
			Field:     consts.BasicFieldRemark,
			FieldName: consts.Description,
			FieldType: consts.LcColumnFieldTypeRichText,
			OldValue:  i.OldIssueBo.Remark,
			NewValue:  i.NewIssueBo.Remark,
		})
		ctx.HandledTrendColumnIds[consts.BasicFieldRemark] = struct{}{}
	}
	return nil
}

func (i *UpdateIssueVo) processOwner(ctx *BatchUpdateIssueContext) errs.SystemErrorInfo {
	if _, ok := i.LcData[consts.BasicFieldOwnerId]; ok {
		//如果所有者变动， 则更新
		same, add, del := int64Slice.CompareSliceAddDelInt64(i.NewIssueBo.OwnerIdI64, i.OldIssueBo.OwnerIdI64)
		if !same {
			i.UpdateOwner = true

			ownerChangeTime := types.NowTime()
			i.PData[consts.TcOwnerChangeTime] = date.FormatTime(ownerChangeTime)
			i.LcData[consts.BasicFieldOwnerChangeTime] = date.FormatTime(ownerChangeTime)
			i.NewIssueBo.OwnerChangeTime = ownerChangeTime

			// 删除的issue relation
			relationType := consts.IssueRelationTypeOwner
			if len(del) > 0 {
				ctx.SqlIssueRelations = append(ctx.SqlIssueRelations, dao.SQLDeleteIssueRelationByIdsType(ctx.Req.OrgId, i.IssueId, del, relationType))
			}
			// 新增的issue relation
			if len(add) > 0 {
				for _, userId := range add {
					ctx.CreateIssueRelations = append(ctx.CreateIssueRelations, &po.PpmPriIssueRelation{
						//Id:           ids.Ids[idx].Id, // ID到最后统一生成
						OrgId:        i.OldIssueBo.OrgId,
						ProjectId:    i.OldIssueBo.ProjectId,
						IssueId:      i.IssueId,
						RelationId:   userId,
						RelationType: relationType,
						Creator:      ctx.Req.UserId,
						CreateTime:   time.Now(),
						Updator:      ctx.Req.UserId,
						UpdateTime:   time.Now(),
						IsDelete:     consts.AppIsNoDelete,
					})
				}
			}
		}
	}
	return nil
}

func (i *UpdateIssueVo) processFollower(ctx *BatchUpdateIssueContext) errs.SystemErrorInfo {
	if _, ok := i.LcData[consts.BasicFieldFollowerIds]; ok {
		//如果所有者变动， 则更新
		same, add, del := int64Slice.CompareSliceAddDelInt64(i.NewIssueBo.FollowerIdsI64, i.OldIssueBo.FollowerIdsI64)
		if !same {
			i.UpdateFollower = true

			// 删除的issue relation
			relationType := consts.IssueRelationTypeFollower
			if len(del) > 0 {
				ctx.SqlIssueRelations = append(ctx.SqlIssueRelations, dao.SQLDeleteIssueRelationByIdsType(ctx.Req.OrgId, i.IssueId, del, relationType))
			}
			// 新增的issue relation
			if len(add) > 0 {
				for _, userId := range add {
					ctx.CreateIssueRelations = append(ctx.CreateIssueRelations, &po.PpmPriIssueRelation{
						//Id:           ids.Ids[idx].Id, // ID到最后统一生成
						OrgId:        i.OldIssueBo.OrgId,
						ProjectId:    i.OldIssueBo.ProjectId,
						IssueId:      i.IssueId,
						RelationId:   userId,
						RelationType: relationType,
						Creator:      ctx.Req.UserId,
						CreateTime:   time.Now(),
						Updator:      ctx.Req.UserId,
						UpdateTime:   time.Now(),
						IsDelete:     consts.AppIsNoDelete,
					})
				}
			}
		}
	}
	return nil
}

// 注意：修改确认人可能需要引起任务状态、任务确认状态的变化，这些处理已集中放在IssueStatus一起处理了，此处仅仅关注auditor本身相关的更新
func (i *UpdateIssueVo) processAuditor(ctx *BatchUpdateIssueContext) errs.SystemErrorInfo {
	if _, ok := i.LcData[consts.BasicFieldAuditorIds]; ok {
		//如果所有者变动， 则更新
		same, add, del := int64Slice.CompareSliceAddDelInt64(i.NewIssueBo.AuditorIdsI64, i.OldIssueBo.AuditorIdsI64)
		if !same {
			i.UpdateAuditor = true

			// 删除的issue relation
			relationType := consts.IssueRelationTypeAuditor
			if len(del) > 0 {
				ctx.SqlIssueRelations = append(ctx.SqlIssueRelations, dao.SQLDeleteIssueRelationByIdsType(ctx.Req.OrgId, i.IssueId, del, relationType))
				ctx.DeleteAuditors[i.IssueId] = del
			}
			// 新增的issue relation
			if len(add) > 0 {
				for _, userId := range add {
					ctx.CreateIssueRelations = append(ctx.CreateIssueRelations, &po.PpmPriIssueRelation{
						//Id:           ids.Ids[idx].Id, // ID到最后统一生成
						OrgId:        i.OldIssueBo.OrgId,
						ProjectId:    i.OldIssueBo.ProjectId,
						IssueId:      i.IssueId,
						RelationId:   userId,
						RelationType: relationType,
						Creator:      ctx.Req.UserId,
						CreateTime:   time.Now(),
						Updator:      ctx.Req.UserId,
						UpdateTime:   time.Now(),
						IsDelete:     consts.AppIsNoDelete,
					})
				}
			}

			// 变更确认人时，给新确认人发通知。如果任务审核状态为通过，那么不需要通知。
			if i.NewIssueBo.AuditStatus != consts.AuditStatusPass {
				// 如果任务还未完成也不需要通知
				var issueStatus *status.StatusInfoBo
				for _, is := range ctx.AllIssueStatus {
					if is.ID == i.NewIssueBo.Status {
						issueStatus = &is
						break
					}
				}
				if issueStatus != nil && issueStatus.Type == consts.StatusTypeComplete {
					i.NoticeAuditorChange = true
				}
			}
		}
	}
	return nil
}

func (i *UpdateIssueVo) processRelating(ctx *BatchUpdateIssueContext, relatingColumnType string) errs.SystemErrorInfo {
	for columnId, column := range ctx.TableColumns {
		if column.Field.Type == relatingColumnType {
			if _, ok := i.LcData[columnId]; ok {
				oldRelatingI := i.OldLcData[columnId]
				newRelatingI := i.LcData[columnId]
				oldRelating := &bo.RelatingIssue{}
				newRelating := &bo.RelatingIssue{}
				jsonx.Copy(oldRelatingI, oldRelating)
				jsonx.Copy(newRelatingI, newRelating)

				// 收集变化的对端任务ID
				_, toAdd, toDel := int64Slice.CompareSliceAddDelString(newRelating.LinkTo, oldRelating.LinkTo)
				_, fromAdd, fromDel := int64Slice.CompareSliceAddDelString(newRelating.LinkFrom, oldRelating.LinkFrom)

				relatingChange := &bo.RelatingChangeBo{
					LinkToAdd:   slice2.StringToInt64Slice(toAdd),
					LinkToDel:   slice2.StringToInt64Slice(toDel),
					LinkFromAdd: slice2.StringToInt64Slice(fromAdd),
					LinkFromDel: slice2.StringToInt64Slice(fromDel),
				}
				switch relatingColumnType {
				case tablePb.ColumnType_relating.String():
					i.RelatingChange[columnId] = relatingChange
				case tablePb.ColumnType_baRelating.String():
					i.BaRelatingChange[columnId] = relatingChange
				case tablePb.ColumnType_singleRelating.String():
					i.SingleRelatingChange[columnId] = relatingChange
				}

				ctx.RelatingIssueIds = append(ctx.RelatingIssueIds, relatingChange.LinkToAdd...)
				ctx.RelatingIssueIds = append(ctx.RelatingIssueIds, relatingChange.LinkToDel...)
				ctx.RelatingIssueIds = append(ctx.RelatingIssueIds, relatingChange.LinkFromAdd...)
				ctx.RelatingIssueIds = append(ctx.RelatingIssueIds, relatingChange.LinkFromDel...)
				ctx.AllRelatingIssueIds = append(ctx.AllRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkTo)...)
				ctx.AllRelatingIssueIds = append(ctx.AllRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkFrom)...)
				ctx.AllRelatingIssueIds = append(ctx.AllRelatingIssueIds, slice2.StringToInt64Slice(newRelating.LinkTo)...)
				ctx.AllRelatingIssueIds = append(ctx.AllRelatingIssueIds, slice2.StringToInt64Slice(newRelating.LinkFrom)...)
				ctx.HandledTrendColumnIds[columnId] = struct{}{}
			} else {
				if v, ok := i.OldLcData[columnId]; ok {
					oldRelating := &bo.RelatingIssue{}
					jsonx.Copy(v, oldRelating)
					ctx.AllRelatingIssueIds = append(ctx.AllRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkTo)...)
					ctx.AllRelatingIssueIds = append(ctx.AllRelatingIssueIds, slice2.StringToInt64Slice(oldRelating.LinkFrom)...)
				}
			}
		}
	}
	return nil
}

// 附件特殊处理，有recycleFlag的逻辑
func (i *UpdateIssueVo) processDocument(ctx *BatchUpdateIssueContext) errs.SystemErrorInfo {
	for columnId, _ := range ctx.UpdateColumnMap {
		if tableColumn, ok := ctx.TableColumns[columnId]; ok {
			if tableColumn.Field.Type == consts.LcColumnFieldTypeDocument {
				if newDocument, ok := i.LcData[columnId]; ok {

					if i.NewIssueBo.Documents == nil {
						i.NewIssueBo.Documents = make(map[string]interface{})
					}
					if i.OldIssueBo.Documents == nil {
						i.OldIssueBo.Documents = make(map[string]interface{})
					}

					newAttachmentsMap := map[string]bo.Attachments{}
					copyer.Copy(newDocument, &newAttachmentsMap)

					oldAttachmentsMap := map[string]bo.Attachments{}
					if oldDocument, ok := i.OldLcData[columnId]; ok {
						copyer.Copy(oldDocument, &oldAttachmentsMap)
					}
					recycleAttachmentIds := []int64{}
					addAttachmentIds := []int64{}
					if len(newAttachmentsMap) > len(oldAttachmentsMap) {
						log.Infof("添加附件, newAttachmentsMap:%v", json.ToJsonIgnoreError(newAttachmentsMap))
						// 创建resource relation
						for resourceId, newAttach := range newAttachmentsMap {
							if _, ok := oldAttachmentsMap[resourceId]; !ok {
								addAttachmentIds = append(addAttachmentIds, newAttach.Id)
							}
						}

					} else if len(newAttachmentsMap) == len(oldAttachmentsMap) {
						for resourceIdStr, old := range oldAttachmentsMap {
							if newAttachmentsMap[resourceIdStr].RecycleFlag == consts.AppIsDeleted && old.RecycleFlag == consts.AppIsNoDelete {
								// 说明是删除附件
								delete(newAttachmentsMap, resourceIdStr)
								// 删除的附件进入回收站
								ctx.CreateRecycleBins = append(ctx.CreateRecycleBins, &po.PpmPrsRecycleBin{
									//Id:           0,  // 最后统一生成
									OrgId:        i.OldIssueBo.OrgId,
									ProjectId:    i.OldIssueBo.ProjectId,
									RelationId:   old.Id,
									RelationType: consts.RecycleTypeAttachment,
									Creator:      ctx.Req.UserId,
									CreateTime:   time.Now(),
									Updator:      ctx.Req.UserId,
									UpdateTime:   time.Now(),
									//Version:      0,  // 最后统一生成
									IsDelete: consts.AppIsNoDelete,
								})
								recycleAttachmentIds = append(recycleAttachmentIds, old.Id)
							} else if newAttachmentsMap[resourceIdStr].RecycleFlag == consts.AppIsNoDelete && old.RecycleFlag == consts.AppIsDeleted {
								// 说明是从回收站恢复附件
								delete(oldAttachmentsMap, resourceIdStr)
							}
						}

					} else if len(newAttachmentsMap) < len(oldAttachmentsMap) {
						log.Errorf("[processDocument] err:%s, newAttachmentsMap:%v, oldAttachmentsMap:%v",
							"附件处理错误", json.ToJsonIgnoreError(newAttachmentsMap), json.ToJsonIgnoreError(oldAttachmentsMap))
					}

					i.NewIssueBo.Documents[columnId] = newAttachmentsMap
					i.OldIssueBo.Documents[columnId] = oldAttachmentsMap

					//if len(fsResourceIds) > 0 {
					//	// 如果附件引用的是飞书云文档等网络文件，需要手工插入资源关联表
					//	ctx.InsertResourceRelations = append(ctx.InsertResourceRelations, bo.ResourceRelationBo{
					//		OrgId:       i.OldIssueBo.OrgId,
					//		ProjectId:   i.OldIssueBo.ProjectId,
					//		IssueId:     i.IssueId,
					//		ResourceIds: fsResourceIds,
					//		SourceType:  consts.OssPolicyTypeLesscodeResource,
					//	})
					//}
					ctx.RecycleAttachments[i.IssueId] = &bo.IssueAttachments{
						ColumnId:    columnId,
						ResourceIds: recycleAttachmentIds,
					}
					ctx.CreateAttachments[i.IssueId] = &bo.IssueAttachments{
						ColumnId:    columnId,
						ResourceIds: addAttachmentIds,
					}
					ctx.HandledTrendColumnIds[columnId] = struct{}{}
				}
			}
		}
	}
	return nil
}

func (i *UpdateIssueVo) processImage(ctx *BatchUpdateIssueContext) errs.SystemErrorInfo {
	for columnId, _ := range ctx.UpdateColumnMap {
		if tableColumn, ok := ctx.TableColumns[columnId]; ok {
			if tableColumn.Field.Type == consts.LcColumnFieldTypeImage {
				if newImage, ok := i.LcData[columnId]; ok {
					if i.NewIssueBo.Images == nil {
						i.NewIssueBo.Images = make(map[string]interface{})
					}
					if i.OldIssueBo.Images == nil {
						i.OldIssueBo.Images = make(map[string]interface{})
					}
					i.NewIssueBo.Images[columnId] = newImage
					if oldImage, ok := i.OldLcData[columnId]; ok {
						i.OldIssueBo.Images[columnId] = oldImage
					}

					ctx.HandledTrendColumnIds[columnId] = struct{}{}
				}
			}
		}

	}
	return nil
}

func (i *UpdateIssueVo) processOrder(ctx *BatchUpdateIssueContext, index int) errs.SystemErrorInfo {
	if ctx.BeforeOrder != nil && ctx.AfterOrder != nil {
		step := (*ctx.AfterOrder - *ctx.BeforeOrder) / float64(len(ctx.UpdateIssues)+1)
		log.Infof("[BatchUpdateIssue] processOrder: %v %v %v %v", *ctx.BeforeOrder, *ctx.AfterOrder, index, step)
		i.LcData[consts.BasicFieldOrder] = *ctx.BeforeOrder + float64(index+1)*step
	} else if ctx.BeforeOrder != nil {
		log.Infof("[BatchUpdateIssue] processOrder: %v %v", *ctx.BeforeOrder, index)
		i.LcData[consts.BasicFieldOrder] = *ctx.BeforeOrder + float64(len(ctx.UpdateIssues)-index)*65536
	} else if ctx.AfterOrder != nil {
		log.Infof("[BatchUpdateIssue] processOrder: %v %v", *ctx.AfterOrder, index)
		i.LcData[consts.BasicFieldOrder] = *ctx.AfterOrder - float64(index+1)*65536
	}

	return nil
}

func (i *UpdateIssueVo) processTrends(ctx *BatchUpdateIssueContext) {
	pushType := consts.PushTypeUpdateIssue
	if ctx.Req.TrendPushType != nil {
		pushType = *ctx.Req.TrendPushType // 用传入的push type覆盖
	}
	customChangeList := i.handleCustomFieldChangeList(ctx, i.LcData, i.OldLcData, ctx.Iterations, ctx.Users, ctx.Depts)
	i.TrendChangeList = append(i.TrendChangeList, customChangeList...)

	log.Infof("[BatchUpdateIssue] processTrends, issueId: %d, changeList: %v, handledColumns: %v", i.IssueId, i.TrendChangeList, ctx.HandledTrendColumnIds)

	// 成员、关联、前后置、附件、图片的ChangeList没有实际生成，是通过HandledTrendColumnIds来判断的，传给Trendsvc做处理
	if pushType == consts.PushTypeUpdateIssue {
		if len(i.TrendChangeList) == 0 && len(ctx.HandledTrendColumnIds) == 0 {
			return
		}
	}

	ext := bo.TrendExtensionBo{}

	// 合并传入的TrendExtensionBo
	if ctx.Req.TrendExtension != nil {
		json.FromJson(json.ToJsonIgnoreError(ctx.Req.TrendExtension), &ext)
	}

	ext.IssueType = "T"
	ext.ObjName = i.OldIssueBo.Title
	ext.ChangeList = i.TrendChangeList
	if ctx.TriggerBy.TriggerBy == consts.TriggerByAutomation {
		ext.AutomationInfo = &bo.TrendAutomationInfo{
			WorkflowId:   ctx.TriggerBy.WorkflowId,
			WorkflowName: ctx.TriggerBy.WorkflowName,
			ExecutionId:  ctx.TriggerBy.ExecutionId,
		}
		if ctx.TriggerBy.TriggerUserId > 0 {
			if user, ok := ctx.Users[ctx.TriggerBy.TriggerUserId]; ok {
				ext.AutomationInfo.TriggerUser = &bo.SimpleUserInfoBo{
					Id:     user.Id,
					Name:   user.Name,
					Avatar: user.Avatar,
					Status: user.Status,
				}
			}
		}
	}

	//最新的计划时间
	var planStartTime *types.Time = nil
	var planEndTime *types.Time = nil
	if i.NewIssueBo.PlanStartTime.IsNotNull() {
		planStartTime = &i.NewIssueBo.PlanStartTime
	}
	if i.NewIssueBo.PlanEndTime.IsNotNull() {
		planEndTime = &i.NewIssueBo.PlanEndTime
	}

	oldValue := convert.ObjectToMap(i.OldIssueBo)
	newValue := convert.ObjectToMap(i.NewIssueBo)
	for k, v := range i.PData {
		if _, ok := newValue[k]; !ok {
			newValue[k] = v
		}
	}
	for k, v := range i.OldLcData {
		if _, ok := oldValue[k]; !ok {
			oldValue[k] = v
		}
	}

	issueTrendsBo := bo.IssueTrendsBo{
		PushType:              pushType,
		OrgId:                 ctx.Req.OrgId,
		OperatorId:            ctx.Req.UserId,
		DataId:                i.OldIssueBo.DataId,
		IssueId:               i.IssueId,
		ParentIssueId:         i.OldIssueBo.ParentId,
		ProjectId:             i.OldIssueBo.ProjectId,
		TableId:               i.OldIssueBo.TableId,
		PriorityId:            i.OldIssueBo.PriorityId,
		ParentId:              i.OldIssueBo.ParentId,
		IssueTitle:            i.NewIssueBo.Title,
		IssueStatusId:         i.NewIssueBo.Status,
		BeforeOwner:           i.OldIssueBo.OwnerIdI64,
		AfterOwner:            i.NewIssueBo.OwnerIdI64,
		BeforeChangeFollowers: i.OldIssueBo.FollowerIdsI64,
		AfterChangeFollowers:  i.NewIssueBo.FollowerIdsI64,
		BeforeChangeAuditors:  i.OldIssueBo.AuditorIdsI64,
		AfterChangeAuditors:   i.NewIssueBo.AuditorIdsI64,
		RelatingChange:        i.RelatingChange,
		BaRelatingChange:      i.BaRelatingChange,
		SingleRelatingChange:  i.SingleRelatingChange,
		BeforeChangeDocuments: i.OldIssueBo.Documents,
		AfterChangeDocuments:  i.NewIssueBo.Documents,
		BeforeChangeImages:    i.OldIssueBo.Images,
		AfterChangeImages:     i.NewIssueBo.Images,
		UpdateOwner:           i.UpdateOwner,
		UpdateFollower:        i.UpdateFollower,
		UpdateAuditor:         i.UpdateAuditor,
		IssuePlanStartTime:    planStartTime,
		IssuePlanEndTime:      planEndTime,
		SourceChannel:         ctx.OrgBaseInfo.SourceChannel,
		NewValue:              json.ToJsonIgnoreError(newValue),
		OldValue:              json.ToJsonIgnoreError(oldValue),
		Ext:                   ext,
	}
	log.Infof("[BatchUpdateIssue] processTrends, issueId: %d, trend: %v", i.IssueId, json.ToJsonIgnoreError(issueTrendsBo))
	ctx.IssueTrends = append(ctx.IssueTrends, &issueTrendsBo)
}

func (i *UpdateIssueVo) handleCustomFieldChangeList(ctx *BatchUpdateIssueContext, newData, oldData map[string]interface{},
	iterationMap map[int64]*bo.IterationBo, userMap map[int64]*uservo.GetAllUserByIdsRespDataUser, deptMap map[int64]*uservo.GetAllDeptByIdsRespDataUser) []bo.TrendChangeListBo {
	customChangeList := make([]bo.TrendChangeListBo, 0)

	headers := make(map[string]lc_table.LcCommonField, 0)
	copyer.Copy(ctx.TableColumns, &headers)

	for k, v := range newData {
		if header, ok := headers[k]; ok {
			oldValue := oldData[k]
			if _, ok := ctx.HandledTrendColumnIds[k]; !ok &&
				k != consts.BasicFieldId &&
				k != consts.BasicFieldIssueId &&
				k != consts.BasicFieldUpdator &&
				k != consts.BasicFieldUpdateTime {
				// 迭代字段特殊处理
				if k == consts.BasicFieldIterationId {
					newName := ""
					oldName := ""
					if i.NewIssueBo.IterationId != 0 {
						if iteration, ok := iterationMap[i.NewIssueBo.IterationId]; ok {
							newName = iteration.Name
						}
					}
					if i.OldIssueBo.IterationId != 0 {
						if iteration, ok := iterationMap[i.OldIssueBo.IterationId]; ok {
							oldName = iteration.Name
						}
					}
					customChangeList = append(customChangeList, bo.TrendChangeListBo{
						Field:     consts.BasicFieldIterationId,
						FieldName: consts.Iteration,
						FieldType: consts.LcColumnFieldTypeSelect,
						OldValue:  oldName,
						NewValue:  newName,
					})
				} else {
					changeNewValue, changeNewValue2 := i.handleCustomChangeValue(ctx.Req.OrgId, v, header, userMap, deptMap)
					changeOldValue, changeOldValue2 := i.handleCustomChangeValue(ctx.Req.OrgId, oldValue, header, userMap, deptMap)
					customChangeList = append(customChangeList, bo.TrendChangeListBo{
						Field:                    header.Name,
						FieldName:                header.Label,
						FieldType:                header.Field.Type,
						AliasTitle:               header.AliasTitle,
						NewValue:                 changeNewValue,
						OldValue:                 changeOldValue,
						NewUserIdsOrDeptIdsValue: changeNewValue2,
						OldUserIdsOrDeptIdsValue: changeOldValue2,
					})
				}
			}
		}
	}
	return customChangeList
}

func (i *UpdateIssueVo) handleCustomChangeValue(orgId int64, v interface{}, header lc_table.LcCommonField, userMap map[int64]*uservo.GetAllUserByIdsRespDataUser, deptMap map[int64]*uservo.GetAllDeptByIdsRespDataUser) (string, []string) {
	newV := ""
	newVList := make([]string, 0)
	if v != nil {
		if header.Field.Type == consts.LcColumnFieldTypeSelect {
			options := map[interface{}]string{}
			for _, option := range header.Field.Props.Select.Options {
				options[str.ToString(option.Id)] = option.Value
			}
			newV = options[str.ToString(v)]
		} else if header.Field.Type == consts.LcColumnFieldTypeGroupSelect && header.Field.Props.GroupSelect != nil {
			options := map[interface{}]string{}
			for _, option := range header.Field.Props.GroupSelect.Options {
				options[str.ToString(option.Id)] = option.Value
			}
			newV = options[str.ToString(v)]
		} else if header.Field.Type == consts.LcColumnFieldTypeMultiSelect {
			options := map[interface{}]string{}
			for _, option := range header.Field.Props.MultiSelect.Options {
				options[str.ToString(option.Id)] = option.Value
			}
			strBuf := bytes.Buffer{}
			vs := cast.ToStringSlice(v)
			for _, id := range vs {
				if value, ok := options[id]; ok {
					strBuf.WriteString(value)
					strBuf.WriteString(",")
				}
			}
			if strBuf.Len() > 0 {
				newV = strBuf.String()[:strBuf.Len()-1]
			}
			//} else if header.Field.Type == "image" {
			//	docs := make([]lc_table.LcDocumentValue, 0)
			//	//docs := make(map[string]lc_table.LcDocumentValue, 0)
			//	_ = json.FromJson(str.ToString(v), &docs)
			//	for _, v := range docs {
			//		newV += v.Name + " "
			//	}
		} else if header.Field.Type == consts.LcColumnFieldTypeMember {
			userIds, _ := businees.LcMemberToUserIdsWithError(v, true, true)
			for _, userId := range userIds {
				if user, ok := userMap[userId]; ok {
					newV += user.Name + " "
					newVList = append(newVList, businees.FormatUserId(userId))
				}
			}
		} else if header.Field.Type == consts.LcColumnFieldTypeDept {
			deptIds := cast.ToStringSlice(v)
			for _, deptId := range deptIds {
				id := cast.ToInt64(deptId)
				if dept, ok := deptMap[id]; ok {
					newV += dept.Name + " "
					newVList = append(newVList, businees.FormatDeptId(id))
				}
			}
		} else {
			newV = str.ToString(v)
		}
	}
	return newV, newVList
}

func (i *UpdateIssueVo) reportEvent(ctx *BatchUpdateIssueContext, userDepts map[string]*uservo.MemberDept) {
	// 上报数据更新事件
	keys := make(map[string]struct{})
	for key, _ := range i.LcData {
		keys[key] = struct{}{}
	}
	for key, value := range i.OldLcData {
		if _, ok := i.LcData[key]; !ok {
			i.LcData[key] = value
		}
	}
	//domain.AssembleLcDataRelated(i.LcData, ctx.TableColumns, ctx.Users, ctx.Depts)
	//domain.AssembleLcDataRelated(i.OldLcData, ctx.TableColumns, ctx.Users, ctx.Depts)
	domain.AssembleDataIds(i.LcData)
	domain.AssembleDataIds(i.OldLcData)
	newData := i.LcData
	oldData := i.OldLcData
	e := &commonvo.DataEvent{
		OrgId:          i.NewIssueBo.OrgId,
		AppId:          i.NewIssueBo.AppId,
		ProjectId:      i.NewIssueBo.ProjectId,
		TableId:        i.NewIssueBo.TableId,
		DataId:         i.NewIssueBo.DataId,
		IssueId:        i.NewIssueBo.Id,
		UserId:         ctx.Req.UserId,
		New:            newData,
		Old:            oldData,
		UpdatedColumns: i.UpdateColumns,
		UserDepts:      userDepts,
		TriggerBy:      ctx.TriggerBy.TriggerBy,
	}
	eventType := msgPb.EventType_DataUpdated
	if _, ok := keys[consts.BasicFieldOrder]; ok {
		if len(i.UpdateColumns) == 0 {
			// 只修改了order，转为EventType_DataMoved事件
			eventType = msgPb.EventType_DataMoved
		}
		e.UpdatedColumns = append(e.UpdatedColumns, consts.BasicFieldOrder)
	}
	openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
	openTraceIdStr := cast.ToString(openTraceId)
	report.ReportDataEvent(eventType, openTraceIdStr, e)

	// 关联前后置更改，给变更的对端任务也要上报更新事件，此时只上报增量事件
	for columnId, relatingChange := range i.RelatingChange {
		relatingLinkTo := map[string]interface{}{
			columnId: bo.RelatingIssue{
				LinkTo: []string{cast.ToString(i.NewIssueBo.Id)},
			},
		}
		relatingLinkFrom := map[string]interface{}{
			columnId: bo.RelatingIssue{
				LinkFrom: []string{cast.ToString(i.NewIssueBo.Id)},
			},
		}
		for _, issueId := range relatingChange.LinkToAdd {
			if issueBo, ok := ctx.RelatingIssues[issueId]; ok {
				e = &commonvo.DataEvent{
					OrgId:       issueBo.OrgId,
					AppId:       issueBo.AppId,
					ProjectId:   issueBo.ProjectId,
					TableId:     issueBo.TableId,
					DataId:      issueBo.DataId,
					IssueId:     issueBo.Id,
					UserId:      ctx.Req.UserId,
					Incremental: relatingLinkFrom, // 增量-添加
					TriggerBy:   ctx.TriggerBy.TriggerBy,
				}
				report.ReportDataEvent(msgPb.EventType_DataUpdated, openTraceIdStr, e)
			}
		}
		for _, issueId := range relatingChange.LinkToDel {
			if issueBo, ok := ctx.RelatingIssues[issueId]; ok {
				e = &commonvo.DataEvent{
					OrgId:       issueBo.OrgId,
					AppId:       issueBo.AppId,
					ProjectId:   issueBo.ProjectId,
					TableId:     issueBo.TableId,
					DataId:      issueBo.DataId,
					IssueId:     issueBo.Id,
					UserId:      ctx.Req.UserId,
					Decremental: relatingLinkFrom, // 增量-删除
					TriggerBy:   ctx.TriggerBy.TriggerBy,
				}
				report.ReportDataEvent(msgPb.EventType_DataUpdated, openTraceIdStr, e)
			}
		}
		for _, issueId := range relatingChange.LinkFromAdd {
			if issueBo, ok := ctx.RelatingIssues[issueId]; ok {
				e = &commonvo.DataEvent{
					OrgId:       issueBo.OrgId,
					AppId:       issueBo.AppId,
					ProjectId:   issueBo.ProjectId,
					TableId:     issueBo.TableId,
					DataId:      issueBo.DataId,
					IssueId:     issueBo.Id,
					UserId:      ctx.Req.UserId,
					Incremental: relatingLinkTo, // 增量-添加
					TriggerBy:   ctx.TriggerBy.TriggerBy,
				}
				report.ReportDataEvent(msgPb.EventType_DataUpdated, openTraceIdStr, e)
			}
		}
		for _, issueId := range relatingChange.LinkFromDel {
			if issueBo, ok := ctx.RelatingIssues[issueId]; ok {
				e = &commonvo.DataEvent{
					OrgId:       issueBo.OrgId,
					AppId:       issueBo.AppId,
					ProjectId:   issueBo.ProjectId,
					TableId:     issueBo.TableId,
					DataId:      issueBo.DataId,
					IssueId:     issueBo.Id,
					UserId:      ctx.Req.UserId,
					Decremental: relatingLinkTo, // 增量-删除
					TriggerBy:   ctx.TriggerBy.TriggerBy,
				}
				report.ReportDataEvent(msgPb.EventType_DataUpdated, openTraceIdStr, e)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//  批量审批/驳回
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type BatchAuditIssueContext struct {
	Req                *projectvo.BatchAuditIssueReqVo
	AuditIssues        map[int64]*AuditIssueVo
	AuditIssueIds      []int64 // 真正需要审批/驳回的任务
	RunningIssueStatus int64   // 进行中的第一个任务状态 用于审批驳回时使用
	AppId              int64
	ProjectId          int64
	TableId            int64
	Data               []map[string]interface{}
}

type AuditIssueVo struct {
	IssueId    int64
	OldIssueBo *bo.IssueBo
	OldData    map[string]interface{} // 老的无码pg表数据(所有字段)
	Data       map[string]interface{}
}

func BatchAuditIssue(reqVo *projectvo.BatchAuditIssueReqVo) ([]int64, errs.SystemErrorInfo) {
	ctx := &BatchAuditIssueContext{
		Req: reqVo,
	}
	log.Infof("[BatchAuditIssue] req: %v", json.ToJsonIgnoreError(reqVo))
	if len(ctx.Req.Input.IssueIds) == 0 {
		return nil, nil
	}

	// 1. 检查参数 组装数据 判断合法性
	errSys := ctx.prepare()
	if errSys != nil {
		return nil, errSys
	}

	log.Infof("[BatchAuditIssue] orgId: %v, userId: %v, auditIssueIds: %v", ctx.Req.OrgId, ctx.Req.UserId, json.ToJsonIgnoreError(ctx.AuditIssueIds))
	// 没有需要处理的任务
	if len(ctx.AuditIssueIds) == 0 {
		return nil, nil
	}

	// 2. 处理
	errSys = ctx.process()
	if errSys != nil {
		return nil, errSys
	}

	return ctx.AuditIssueIds, nil
}

func (ctx *BatchAuditIssueContext) prepare() errs.SystemErrorInfo {
	if ctx.Req.Input.AuditStatus != consts.AuditStatusPass &&
		ctx.Req.Input.AuditStatus != consts.AuditStatusReject {
		return errs.ParamError
	}

	// 批量任务数上限
	if len(ctx.Req.Input.IssueIds) > BATCH_SIZE {
		return errs.BatchOperateTooManyRows
	}

	ctx.AuditIssues = make(map[int64]*AuditIssueVo)

	// 获取修改前的老数据
	errSys := ctx.fetchOldDatas()
	if errSys != nil {
		log.Errorf("[BatchAuditIssue] 获取老数据失败 org:%d user:%d, err: %v", ctx.Req.OrgId, ctx.Req.UserId, errSys)
		return errSys
	}

	for _, issueBo := range ctx.AuditIssues {
		ctx.AppId = issueBo.OldIssueBo.AppId
		ctx.ProjectId = issueBo.OldIssueBo.ProjectId
		ctx.TableId = issueBo.OldIssueBo.TableId
		break
	}
	// 未归属项目的任务，暂不允许这个操作
	if ctx.AppId <= 0 || ctx.ProjectId <= 0 || ctx.TableId <= 0 {
		return errs.ParamError

	}

	// 校验项目类型
	project, errSys := domain.GetProjectSimple(ctx.Req.OrgId, ctx.ProjectId)
	if errSys != nil {
		log.Errorf("[BatchAuditIssue] 获取Project失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
			ctx.Req.OrgId, ctx.AppId, ctx.ProjectId, ctx.TableId, ctx.Req.UserId, errSys)
		return errSys
	}
	if project.ProjectTypeId != consts.ProjectTypeNormalId {
		return errs.ParamError
	}

	// 获取表头 (驳回可能要回滚任务状态，通过也只需要保持已完成状态，用不到表头)
	if ctx.Req.Input.AuditStatus == consts.AuditStatusReject {
		tableColumns, errSys := domain.GetTableColumnsMap(ctx.Req.OrgId, ctx.TableId, []string{consts.BasicFieldIssueStatus})
		if errSys != nil {
			log.Errorf("[BatchAuditIssue] 获取表头失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
				ctx.Req.OrgId, ctx.AppId, ctx.ProjectId, ctx.TableId, ctx.Req.UserId, errSys)
			return errSys
		}
		if column, ok := tableColumns[consts.BasicFieldIssueStatus]; ok {
			allIssueStatus := domain.GetStatusListFromStatusColumn(column)
			for _, is := range allIssueStatus {
				if is.Type == consts.StatusTypeRunning {
					ctx.RunningIssueStatus = is.ID
					break
				}
			}
		}
		if ctx.RunningIssueStatus == 0 {
			return errs.IssueStatusNotExist
		}
	}

	for _, auditIssue := range ctx.AuditIssues {
		// 跳过：已删除的数据
		if cast.ToInt64(auditIssue.OldData[consts.BasicFieldRecycleFlag]) != consts.AppIsNoDelete {
			continue
		}

		// 跳过：已经审批过的任务
		if auditIssue.OldIssueBo.AuditStatus == consts.AuditStatusPass ||
			auditIssue.OldIssueBo.AuditStatus == consts.AuditStatusReject {
			continue
		}

		// 跳过：任务状态不是完成状态的任务
		if auditIssue.OldIssueBo.IssueStatusType != consts.StatusTypeComplete {
			continue
		}

		// 跳过：当前操作人不是任务确认人的任务
		if has, err := slice.Contain(auditIssue.OldIssueBo.AuditorIdsI64, ctx.Req.UserId); !has || err != nil {
			continue
		}

		// 跳过：当前操作人已经审批/驳回过任务
		if auditIssue.OldIssueBo.AuditStatusDetail != nil {
			if auditStatus, ok := auditIssue.OldIssueBo.AuditStatusDetail[cast.ToString(ctx.Req.UserId)]; ok {
				if auditStatus == consts.AuditStatusPass || auditStatus == consts.AuditStatusReject {
					continue
				}
			}
		}
		ctx.AuditIssueIds = append(ctx.AuditIssueIds, auditIssue.IssueId)
	}

	return nil
}

func (ctx *BatchAuditIssueContext) fetchOldDatas() errs.SystemErrorInfo {
	oldDatas, errSys := domain.GetIssueInfosMapLcByIssueIds(ctx.Req.OrgId, ctx.Req.UserId, ctx.Req.Input.IssueIds)
	if errSys != nil {
		return errSys
	}
	for _, d := range oldDatas {
		oldIssueBo, errSys := domain.ConvertIssueDataToIssueBo(d)
		if errSys != nil {
			return errSys
		}

		auditIssue := &AuditIssueVo{}
		auditIssue.IssueId = oldIssueBo.Id
		auditIssue.OldData = d
		auditIssue.OldIssueBo = oldIssueBo
		ctx.AuditIssues[auditIssue.IssueId] = auditIssue
	}
	return nil
}

func (ctx *BatchAuditIssueContext) process() errs.SystemErrorInfo {
	var errSys errs.SystemErrorInfo

	// 处理更新
	errSys = ctx.processAudit()
	if errSys != nil {
		return errSys
	}

	// 后续处理（异步）
	asyn.Execute(ctx.processHooks)

	return nil
}

func (ctx *BatchAuditIssueContext) processAudit() errs.SystemErrorInfo {
	for _, issueId := range ctx.AuditIssueIds {
		ai, _ := ctx.AuditIssues[issueId]
		ai.OldIssueBo.AuditStatusDetail[cast.ToString(ctx.Req.UserId)] = ctx.Req.Input.AuditStatus
		ai.Data = make(map[string]interface{})

		ai.Data[consts.BasicFieldId] = issueId
		if ctx.Req.Input.AuditStatus == consts.AuditStatusPass {
			// 只要一个人审批通过，即可修改审批状态为通过
			ai.Data[consts.BasicFieldAuditStatus] = ctx.Req.Input.AuditStatus
			ai.Data[consts.BasicFieldAuditStatusDetail] = ai.OldIssueBo.AuditStatusDetail
		} else {
			// 检查是否所有人都审批驳回
			isAllReject := true
			for _, auditorId := range ai.OldIssueBo.AuditorIdsI64 {
				if auditStatus, ok := ai.OldIssueBo.AuditStatusDetail[cast.ToString(auditorId)]; !ok {
					isAllReject = false // 有人尚未审核
					break
				} else if auditStatus != consts.AuditStatusReject {
					isAllReject = false
					break
				}
			}
			// 如果所有人都审批驳回，则回到进行中状态
			if isAllReject {
				ai.Data[consts.BasicFieldAuditStatus] = consts.AuditStatusNotView
				ai.Data[consts.BasicFieldIssueStatus] = ctx.RunningIssueStatus
			}
			ai.Data[consts.BasicFieldAuditStatusDetail] = ai.OldIssueBo.AuditStatusDetail
		}
		ctx.Data = append(ctx.Data, ai.Data)
	}

	var attachments []bo.ResourceInfoBo
	now := time.Now()
	for _, attachment := range ctx.Req.Input.Attachments {
		attachments = append(attachments, bo.ResourceInfoBo{
			Url:        attachment.URL,
			Name:       attachment.Name,
			Suffix:     attachment.Suffix,
			Size:       attachment.Size,
			UploadTime: now,
		})
	}

	trendPushType := consts.PushTypeAuditIssue
	trendExt := bo.TrendExtensionBo{AuditInfo: &bo.TrendAuditInfo{
		Status:      ctx.Req.Input.AuditStatus,
		Remark:      ctx.Req.Input.Message,
		Attachments: attachments,
	}}
	batchUpdateReq := &projectvo.BatchUpdateIssueReqInnerVo{
		OrgId:          ctx.Req.OrgId,
		UserId:         ctx.Req.UserId,
		AppId:          ctx.AppId,
		ProjectId:      ctx.ProjectId,
		TableId:        ctx.TableId,
		Data:           ctx.Data,
		TrendPushType:  &trendPushType,
		TrendExtension: &trendExt,
	}
	return BatchUpdateIssue(batchUpdateReq, true, &projectvo.TriggerBy{
		TriggerBy: consts.TriggerByNormal,
	})
}

func (ctx *BatchAuditIssueContext) processHooks() {
	// 发送审批结果给发起人
	for _, issueId := range ctx.AuditIssueIds {
		if ai, ok := ctx.AuditIssues[issueId]; ok {
			// 没有审批结果的不需要推卡片
			if _, ok := ai.Data[consts.BasicFieldAuditStatus]; ok {
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//  批量催办
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type BatchUrgeIssueContext struct {
	Req                     *projectvo.BatchUrgeIssueReqVo
	UrgeIssues              map[int64]*UrgeIssueVo
	OrgBaseInfo             *bo.BaseOrgInfoBo
	Project                 *bo.ProjectBo
	TableMeta               *projectvo.TableMetaData
	UserBaseInfos           map[int64]*bo.BaseUserInfoBo
	AppId                   int64
	ProjectId               int64
	TableId                 int64
	UrgeOwnerIssueIds       []int64
	UrgeAuditorIssueIds     []int64
	UrgeOwnerSuccIssueIds   []int64
	UrgeAuditorSuccIssueIds []int64
	AllParentTitles         map[int64]string
}

type UrgeIssueVo struct {
	IssueId    int64
	OldLcData  map[string]interface{} // 老的无码pg表数据(所有字段)
	OldIssueBo *bo.IssueBo
}

func BatchUrgeIssue(reqVo *projectvo.BatchUrgeIssueReqVo) ([]int64, errs.SystemErrorInfo) {
	ctx := &BatchUrgeIssueContext{
		Req: reqVo,
	}
	log.Infof("[BatchUrgeIssue] req: %v", json.ToJsonIgnoreError(reqVo))
	if len(ctx.Req.Input.IssueIds) == 0 {
		return nil, nil
	}

	// 1. 检查参数 组装数据 判断合法性
	errSys := ctx.prepare()
	if errSys != nil {
		return nil, errSys
	}

	// 没有需要处理的任务
	log.Infof("[BatchUrgeIssue] orgId: %v, userId: %v, urgeOwnerIssueIds: %v, urgeAuditorIssueIds: %v",
		ctx.Req.OrgId, ctx.Req.UserId, ctx.UrgeOwnerIssueIds, ctx.UrgeAuditorIssueIds)
	if len(ctx.UrgeOwnerIssueIds) == 0 && len(ctx.UrgeAuditorIssueIds) == 0 {
		return nil, nil
	}

	// 2. 处理
	errSys = ctx.process()
	if errSys != nil {
		return nil, errSys
	}

	allIssueIds := append(ctx.UrgeOwnerSuccIssueIds, ctx.UrgeAuditorSuccIssueIds...)
	return allIssueIds, nil
}

func (ctx *BatchUrgeIssueContext) prepare() errs.SystemErrorInfo {
	ctx.UserBaseInfos = make(map[int64]*bo.BaseUserInfoBo)

	// 批量任务数上限
	if len(ctx.Req.Input.IssueIds) > BATCH_SIZE {
		return errs.BatchOperateTooManyRows
	}

	ctx.UrgeIssues = make(map[int64]*UrgeIssueVo)
	ctx.AllParentTitles = make(map[int64]string)

	// 获取修改前的老数据
	errSys := ctx.fetchOldDatas()
	if errSys != nil {
		log.Errorf("[BatchUrgeIssue] 获取老数据失败 org:%d user:%d, err: %v", ctx.Req.OrgId, ctx.Req.UserId, errSys)
		return errSys
	}

	for _, issueBo := range ctx.UrgeIssues {
		ctx.AppId = issueBo.OldIssueBo.AppId
		ctx.ProjectId = issueBo.OldIssueBo.ProjectId
		ctx.TableId = issueBo.OldIssueBo.TableId
		break
	}
	// 未归属项目的任务，暂不允许这个操作
	if ctx.AppId <= 0 || ctx.ProjectId <= 0 || ctx.TableId <= 0 {
		return errs.ParamError
	}

	// 获取org base info
	orgResp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: ctx.Req.OrgId})
	if orgResp.Failure() {
		log.Errorf("[BatchUrgeIssue] 获取Org失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
			ctx.Req.OrgId, ctx.AppId, ctx.ProjectId, ctx.TableId, ctx.Req.UserId, orgResp.Error())
		return orgResp.Error()
	}
	ctx.OrgBaseInfo = orgResp.BaseOrgInfo

	// 获取项目信息
	ctx.Project, errSys = domain.GetProject(ctx.Req.OrgId, ctx.ProjectId)
	if errSys != nil {
		return errSys
	}

	// 获取表格信息
	ctx.TableMeta, errSys = domain.GetTableByTableId(ctx.Req.OrgId, ctx.Req.UserId, ctx.TableId)
	if errSys != nil {
		log.Errorf("[BatchUrgeIssue] 获取表信息失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
			ctx.Req.OrgId, ctx.AppId, ctx.ProjectId, ctx.TableId, ctx.Req.UserId, errSys)
		return errSys
	}

	// 得到项目管理员
	orgAdmins, errSys := domain.GetAdminUserIdsOfOrg(ctx.Req.OrgId, ctx.Project.AppId)
	if errSys != nil {
		log.Errorf("[BatchAuditIssue] GetAdminUserIdsOfOrg err: %v", errSys)
		return errSys
	}
	appAdmins := make(map[int64]struct{})
	for _, id := range orgAdmins {
		appAdmins[id] = struct{}{}
	}
	for _, id := range ctx.Project.OwnerIds {
		appAdmins[id] = struct{}{}
	}
	_, isAppAdmin := appAdmins[ctx.Req.UserId]

	for _, urgeIssue := range ctx.UrgeIssues {
		// 跳过：已删除的数据
		if cast.ToInt64(urgeIssue.OldLcData[consts.BasicFieldRecycleFlag]) != consts.AppIsNoDelete {
			continue
		}

		// 跳过：标题为空的任务
		if urgeIssue.OldIssueBo.Title == "" {
			continue
		}

		isIssueOwner, _ := slice.Contain(urgeIssue.OldIssueBo.OwnerIdI64, ctx.Req.UserId)

		if urgeIssue.OldIssueBo.IssueStatusType != consts.StatusTypeComplete {
			// 未完成的任务催办负责人（只有App管理员能催办负责人）
			if isAppAdmin {
				ctx.UrgeOwnerIssueIds = append(ctx.UrgeOwnerIssueIds, urgeIssue.IssueId)
			}
		} else {
			// 已完成的任务催办确认人（需审批确认项目）（只有App管理员或任务负责人能催办确认人）
			if ctx.Project.ProjectTypeId == consts.ProjectTypeNormalId &&
				urgeIssue.OldIssueBo.AuditStatus == consts.AuditStatusNotView {
				if isAppAdmin || isIssueOwner {
					ctx.UrgeAuditorIssueIds = append(ctx.UrgeAuditorIssueIds, urgeIssue.IssueId)
				}
			}
		}
	}

	// 检查催办时间
	if len(ctx.UrgeOwnerIssueIds) > 0 {
		rs, _ := domain.CheckIssuesAllowUrge(ctx.Req.OrgId, ctx.ProjectId, ctx.UrgeOwnerIssueIds, sconsts.CacheUrgeIssue)
		issueIds := make([]int64, 0)
		for i := 0; i < len(rs); i++ {
			if rs[i] {
				issueIds = append(issueIds, ctx.UrgeOwnerIssueIds[i])
			}
		}
		ctx.UrgeOwnerIssueIds = issueIds
	}

	if len(ctx.UrgeAuditorIssueIds) > 0 {
		rs, _ := domain.CheckIssuesAllowUrge(ctx.Req.OrgId, ctx.ProjectId, ctx.UrgeAuditorIssueIds, sconsts.CacheUrgeIssueAudit)
		issueIds := make([]int64, 0)
		for i := 0; i < len(rs); i++ {
			if rs[i] {
				issueIds = append(issueIds, ctx.UrgeAuditorIssueIds[i])
			}
		}
		ctx.UrgeAuditorIssueIds = issueIds
	}

	// 拉取相关信息
	userIdMap := make(map[int64]struct{})
	userIdMap[ctx.Req.UserId] = struct{}{}
	parentIdMap := make(map[int64]struct{})
	for _, issueId := range ctx.UrgeOwnerIssueIds {
		urgeIssue, _ := ctx.UrgeIssues[issueId]
		for _, userId := range urgeIssue.OldIssueBo.OwnerIdI64 {
			userIdMap[userId] = struct{}{}
		}
		if urgeIssue.OldIssueBo.ParentId > 0 {
			parentIdMap[urgeIssue.OldIssueBo.ParentId] = struct{}{}
		}
	}
	for _, issueId := range ctx.UrgeAuditorIssueIds {
		urgeIssue, _ := ctx.UrgeIssues[issueId]
		for _, userId := range urgeIssue.OldIssueBo.AuditorIdsI64 {
			userIdMap[userId] = struct{}{}
		}
		for _, userId := range urgeIssue.OldIssueBo.OwnerIdI64 {
			userIdMap[userId] = struct{}{}
		}
		if urgeIssue.OldIssueBo.ParentId > 0 {
			parentIdMap[urgeIssue.OldIssueBo.ParentId] = struct{}{}
		}
	}
	// 用户信息
	if len(userIdMap) > 0 {
		userIds := make([]int64, 0, len(userIdMap))
		for id, _ := range userIdMap {
			userIds = append(userIds, id)
		}
		userResp := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
			OrgId:   ctx.Req.OrgId,
			UserIds: userIds,
		})
		if userResp.Failure() {
			log.Errorf("[BatchUrgeIssue] 获取Usr失败 org:%d app:%d proj:%d table:%d user:%d, err: %v",
				ctx.Req.OrgId, ctx.AppId, ctx.ProjectId, ctx.TableId, ctx.Req.UserId, errSys)
			return orgResp.Error()
		}
		for i, info := range userResp.BaseUserInfos {
			ctx.UserBaseInfos[info.UserId] = &userResp.BaseUserInfos[i]
		}
	}
	// 父任务信息
	if len(parentIdMap) > 0 {
		parentIds := make([]int64, 0, len(parentIdMap))
		for id, _ := range parentIdMap {
			parentIds = append(parentIds, id)
		}
		data, errSys := domain.GetIssueInfosMapLcByIssueIds(ctx.Req.OrgId, ctx.Req.UserId, parentIds, lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle))
		if errSys != nil {
			return errSys
		}
		for _, d := range data {
			issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
			title := cast.ToString(d[consts.BasicFieldTitle])
			if issueId > 0 {
				ctx.AllParentTitles[issueId] = title
			}
		}
	}

	return nil
}

func (ctx *BatchUrgeIssueContext) fetchOldDatas() errs.SystemErrorInfo {
	oldDatas, errSys := domain.GetIssueInfosMapLcByIssueIds(ctx.Req.OrgId, ctx.Req.UserId, ctx.Req.Input.IssueIds)
	if errSys != nil {
		return errSys
	}
	for _, d := range oldDatas {
		oldIssueBo, errSys := domain.ConvertIssueDataToIssueBo(d)
		if errSys != nil {
			return errSys
		}

		urgeIssue := &UrgeIssueVo{}
		urgeIssue.IssueId = oldIssueBo.Id
		urgeIssue.OldLcData = d
		urgeIssue.OldIssueBo = oldIssueBo
		ctx.UrgeIssues[urgeIssue.IssueId] = urgeIssue
	}
	return nil
}

func (ctx *BatchUrgeIssueContext) process() errs.SystemErrorInfo {
	// 催办负责人
	for _, issueId := range ctx.UrgeOwnerIssueIds {
		urgeIssue, ok := ctx.UrgeIssues[issueId]
		if !ok {
			continue
		}

		var toOpenIds []string
		var owners []*bo.BaseUserInfoBo
		for _, userId := range urgeIssue.OldIssueBo.OwnerIdI64 {
			if userBaseInfo, ok := ctx.UserBaseInfos[userId]; ok {
				// 自己不催办
				if userId != ctx.Req.UserId {
					toOpenIds = append(toOpenIds, userBaseInfo.OutUserId)
				}
				owners = append(owners, userBaseInfo)
			}
		}
		if len(toOpenIds) == 0 {
			continue
		}

		ctx.UrgeOwnerSuccIssueIds = append(ctx.UrgeOwnerSuccIssueIds, issueId)
	}

	// 催办确认人
	for _, issueId := range ctx.UrgeAuditorIssueIds {
		urgeIssue, ok := ctx.UrgeIssues[issueId]
		if !ok {
			continue
		}

		var toOpenIds []string
		var owners []*bo.BaseUserInfoBo
		for _, userId := range urgeIssue.OldIssueBo.AuditorIdsI64 {
			// 自己不催办
			if userId != ctx.Req.UserId {
				// 已经审批过的不催办
				auditStatus, ok := urgeIssue.OldIssueBo.AuditStatusDetail[cast.ToString(userId)]
				if !ok || (auditStatus != consts.AuditStatusPass && auditStatus != consts.AuditStatusReject) {
					if userBaseInfo, ok := ctx.UserBaseInfos[userId]; ok {
						toOpenIds = append(toOpenIds, userBaseInfo.OutUserId)
					}
				}
			}
		}
		for _, userId := range urgeIssue.OldIssueBo.OwnerIdI64 {
			if userBaseInfo, ok := ctx.UserBaseInfos[userId]; ok {
				owners = append(owners, userBaseInfo)
			}
		}
		if len(toOpenIds) == 0 {
			continue
		}

		ctx.UrgeAuditorSuccIssueIds = append(ctx.UrgeAuditorSuccIssueIds, issueId)
	}

	// 记录催办时间
	if len(ctx.UrgeOwnerSuccIssueIds) > 0 {
		errCache := domain.UpdateIssuesUrgeTime(ctx.Req.OrgId, ctx.ProjectId, ctx.UrgeOwnerSuccIssueIds, sconsts.CacheUrgeIssue)
		if errCache != nil {
			log.Errorf("[UrgeAuditIssue] UpdateIssuesUrgeTime err:%v", errCache)
		}
	}
	if len(ctx.UrgeAuditorSuccIssueIds) > 0 {
		errCache := domain.UpdateIssuesUrgeTime(ctx.Req.OrgId, ctx.ProjectId, ctx.UrgeAuditorSuccIssueIds, sconsts.CacheUrgeIssueAudit)
		if errCache != nil {
			log.Errorf("[UrgeAuditIssue] UpdateIssuesUrgeTime err:%v", errCache)
		}
	}
	return nil
}
