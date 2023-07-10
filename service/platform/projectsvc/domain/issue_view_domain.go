package domain

import (
	"fmt"
	"html"
	"strconv"
	"time"

	"github.com/star-table/common/core/util/slice"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/appvo"
	"github.com/star-table/polaris-backend/facade/appfacade"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

// 创建任务视图
func CreateIssueView(orgId, opUserId int64, view *vo.CreateIssueViewReq) (*po.PpmPriIssueView, errs.SystemErrorInfo) {
	owner := int64(0)
	view, err := CheckDataForCreateIssueView(view)
	if err != nil {
		return nil, err
	}
	if *view.IsPrivate {
		owner = opUserId
	}
	if view.Config == "" {
		view.Config = "{}"
	}
	viewPo := po.PpmPriIssueView{
		ID:                  0,
		OrgID:               orgId,
		ProjectID:           *view.ProjectID,
		ProjectObjectTypeId: *view.ProjectObjectTypeID,
		Type:                int8(*view.Type),
		ViewName:            view.ViewName,
		Config:              view.Config,
		Remark:              *view.Remark,
		Owner:               owner,
		Sort:                *view.Sort,
		DelFlag:             consts.AppIsNoDelete,
		Creator:             opUserId,
		CreateTime:          time.Time{},
		Updator:             opUserId,
		UpdateTime:          time.Time{},
	}
	newId, err := idfacade.ApplyPrimaryIdRelaxed(viewPo.TableName())
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	viewPo.ID = newId
	if err := dao.CreateIssueView(viewPo); err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	return &viewPo, nil
}

// 检查数据并对缺省值初始化
func CheckDataForCreateIssueView(view *vo.CreateIssueViewReq) (*vo.CreateIssueViewReq, errs.SystemErrorInfo) {
	viewNameLen := len(view.ViewName)
	if viewNameLen < 1 || viewNameLen > 60 {
		return view, errs.IssueViewNameLenInValid
	}
	if view.IsPrivate == nil {
		defaultVal := false
		view.IsPrivate = &defaultVal
	}
	if view.ProjectID == nil {
		defaultVal := int64(0)
		view.ProjectID = &defaultVal
	}
	if view.Type == nil {
		defaultVal := 1
		view.Type = &defaultVal
	}
	if view.Sort == nil {
		defaultVal := int64(0)
		view.Sort = &defaultVal
	}
	if view.Remark == nil {
		defaultVal := ""
		view.Remark = &defaultVal
	}
	if view.ProjectObjectTypeID == nil {
		defaultVal := int64(0)
		view.ProjectObjectTypeID = &defaultVal
	}
	// 对用户输入进行转义
	view.ViewName = html.EscapeString(view.ViewName)
	*view.Remark = html.EscapeString(*view.Remark)
	return view, nil
}

// 获取一个任务视图
func GetOneTaskView(orgId, viewId int64) (*po.PpmPriIssueView, errs.SystemErrorInfo) {
	_, views, err := dao.SelectIssueViewsByCond(1, 1, db.Cond{
		consts.TcOrgId:   orgId,
		consts.TcId:      viewId,
		consts.TcDelFlag: consts.AppIsNoDelete,
	}, nil)
	if err == db.ErrNoMoreRows || len(views) < 1 {
		return nil, errs.IssueViewNotExist
	}
	return &views[0], nil
}

// 获取任务视图列表
func GetTaskViewList(orgId, userId int64, input *vo.GetIssueViewListReq) (int64, []po.PpmPriIssueView, errs.SystemErrorInfo) {
	if input.Page == nil {
		defaultVal := 1
		input.Page = &defaultVal
	}
	if input.Size == nil {
		defaultVal := 10
		input.Size = &defaultVal
	}
	sortRule := db.Raw(fmt.Sprintf("`%s` asc", consts.TcId))
	if input.SortType != nil {
		// 排序类型。1创建时间顺序，2创建时间倒序，3更新时间顺序，4更新时间倒序。
		switch *input.SortType {
		case 1:
			sortRule = db.Raw(fmt.Sprintf("sort asc, %s asc", consts.TcId))
		case 2:
			sortRule = db.Raw(fmt.Sprintf("sort asc, %s desc", consts.TcId))
		case 3:
			sortRule = db.Raw(fmt.Sprintf("sort asc, %s asc", consts.TcUpdateTime))
		case 4:
			sortRule = db.Raw(fmt.Sprintf("sort asc, %s desc", consts.TcUpdateTime))
		}
	}
	cond := assemblyCondForGetTaskViewList(orgId, userId, input)
	total, views, err := dao.SelectIssueViewsByCond(*input.Page, *input.Size, cond, sortRule)
	if err == db.ErrNoMoreRows {
		return 0, make([]po.PpmPriIssueView, 0), nil
	}
	return total, views, nil
}

// 任务视图条件筛选统计，组装查询条件
func assemblyCondForGetTaskViewList(orgId, userId int64, input *vo.GetIssueViewListReq) db.Cond {
	cond := db.Cond{
		consts.TcOrgId:   orgId,
		consts.TcDelFlag: consts.AppIsNoDelete,
	}
	if input.Ids != nil && len(input.Ids) > 0 {
		cond[consts.TcId] = db.In(input.Ids)
	}
	if input.ProjectID != nil {
		cond[consts.TcProjectId] = *input.ProjectID
	}
	if input.ViewName != nil {
		cond[consts.TcIssueViewName] = db.Like(fmt.Sprintf("%%%s%%", *input.ViewName))
	}
	if input.IsPrivate != nil {
		if *input.IsPrivate {
			cond[consts.TcOwner] = userId
		} else {
			cond[consts.TcOwner] = db.In([]int64{0, userId})
		}
	}
	if input.Type != nil {
		cond[consts.TcType] = *input.Type
	}
	if input.ProjectObjectTypeID != nil {
		cond[consts.TcProjectObjectTypeId] = *input.ProjectObjectTypeID
	}

	return cond
}

// CreateProjectDefaultView 创建项目时创建无码默认视图（isFirst:新手指南表格视图增加筛选项）
func CreateProjectDefaultView(orgId int64, projectId int64, appId int64, projectTypeId int64, isFirst *bool, tableId *int64) errs.SystemErrorInfo {
	//projectBo, err := LoadProjectAuthBo(orgId, projectId)
	//if err != nil {
	//	log.Error(err)
	//	return err
	//}
	tableIds := []int64{}
	if tableId != nil {
		tableIds = append(tableIds, *tableId)
	} else {
		tableList, tableListErr := GetAppTableList(orgId, appId)
		if tableListErr != nil {
			log.Errorf("[CreateProjectDefaultView]GetAppTableList failed:%v, orgId:%d, appId:%d", tableListErr, orgId, appId)
			return tableListErr
		}
		if len(tableList) == 0 {
			return errs.ProjectHasNoTableList
		}
		for _, typeBo := range tableList {
			tableIds = append(tableIds, typeBo.TableId)
		}
	}
	allStatusBos, allStatusErr := GetIssueAllStatus(orgId, []int64{projectId}, tableIds)
	if allStatusErr != nil {
		log.Errorf("[CreateProjectDefaultView]GetIssueAllStatus failed:%v", allStatusErr)
		return allStatusErr
	}

	defaultHiddenColumnIds := []string{
		consts.BasicFieldCode,
		//consts.BasicFieldPlanStartTime,
		consts.BasicFieldRemark,
		consts.BasicFieldFollowerIds,
		consts.BasicFieldPriority,
		consts.BasicFieldDemandSource,
		consts.BasicFieldDemandType,
		consts.BasicFieldBugProperty,
		consts.BasicFieldBugType,
		consts.BasicFieldProjectObjectTypeId,
	}

	defaultConfig := GetDefaultViews(defaultHiddenColumnIds, projectTypeId)

	reqs := []appvo.CreateAppViewData{}
	for _, tId := range tableIds {
		for _, view := range defaultConfig {
			view.Config.TableId = strconv.FormatInt(tId, 10)
			if view.ViewName == consts.OverDueSoonViewName {
				if statusList, ok := allStatusBos[tId]; ok {
					options := []map[string]interface{}{}
					valueIds := []int64{}
					for _, infoBo := range statusList {
						if infoBo.Type != consts.StatusTypeComplete {
							options = append(options, map[string]interface{}{
								"color":     infoBo.BgStyle,
								"fontcolor": infoBo.FontStyle,
								"label":     infoBo.Name,
								"parentId":  infoBo.Type,
								"sort":      infoBo.Sort,
								"value":     infoBo.ID,
							})
							valueIds = append(valueIds, infoBo.ID)
						}
					}
					cur := view.Config.Condition.(bo.AppViewLessShowCondition)
					cur.Conds = append(cur.Conds, bo.AppViewLessShowCondition{
						Column:    consts.BasicFieldIssueStatus,
						FieldType: "multiStatus",
						Type:      "in",
						Value:     valueIds,
						Option:    options,
					})
					view.Config.Condition = cur
					view.Config.RealCondition = bo.AppViewLessCondition{
						Type: "and",
						Conds: []bo.AppViewLessCondition{
							bo.AppViewLessCondition{
								Column: consts.BasicFieldPlanEndTime,
								Type:   "gt",
								Value:  "${today}",
								Values: []interface{}{"${today}"},
							},
							bo.AppViewLessCondition{
								Column: consts.BasicFieldPlanEndTime,
								Type:   "lt",
								Value:  "${afterday:3}",
								Values: []interface{}{"${afterday:3}"},
							},
							bo.AppViewLessCondition{
								Column: consts.BasicFieldIssueStatus,
								Type:   "in",
								Value:  valueIds,
								Values: slice.ToSlice(valueIds),
							},
						},
					}
				}
			}
			reqs = append(reqs, appvo.CreateAppViewData{
				AppId:   appId,
				OrgId:   orgId,
				Config:  json.ToJsonIgnoreError(view.Config),
				Name:    view.ViewName,
				OwnerId: 0,
				Type:    view.Type,
			})
		}
	}
	resp := appfacade.CreateAppViews(appvo.CreateAppViewReq{Reqs: reqs})
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}

	return nil
}

func GetDefaultViews(defaultHiddenColumnIds []string, projectTypeId int64) []bo.LcAppView {
	appDefaultViews := []bo.LcAppView{}
	appDefaultViews = append(appDefaultViews, bo.LcAppView{
		ViewName: consts.TabularViewName,
		Type:     consts.TabularViewType,
		Config: bo.AppViewConfig{
			Orders: []bo.AppViewOrder{
				{
					Asc:    true,
					Column: consts.ProBasicFieldSort,
				},
			},
			Condition:           struct{}{},
			ProjectObjectTypeId: 0,
			RealCondition:       struct{}{},
			HiddenColumnIds:     defaultHiddenColumnIds,
		},
	})
	if projectTypeId == consts.ProjectTypeEmpty {
		return appDefaultViews
	}
	appDefaultViews = append(appDefaultViews, bo.LcAppView{
		ViewName: consts.OwnerViewName,
		Type:     consts.KanbanViewType,
		Config: bo.AppViewConfig{
			Orders: []bo.AppViewOrder{
				{
					Asc:    true,
					Column: consts.ProBasicFieldSort,
				},
			},
			ProjectObjectTypeId: 0,
			Condition: bo.AppViewLessShowCondition{
				Column:    consts.BasicFieldOwnerId,
				FieldType: "member",
				Type:      "values_in",
				Values: []interface{}{[]interface{}{
					bo.MemberConditionValues{
						Id:     "${current_user}",
						Name:   "本人",
						Avatar: consts.DefaultAvatar,
					},
				}},
			},
			RealCondition: bo.AppViewLessCondition{
				Column: consts.BasicFieldOwnerId,
				Type:   "values_in",
				Value:  "",
				Values: []interface{}{consts.KeyCurrentUser},
			},
			HiddenColumnIds: defaultHiddenColumnIds,
		},
	})

	appDefaultViews = append(appDefaultViews, bo.LcAppView{
		ViewName: consts.OverDueSoonViewName,
		Type:     consts.TabularViewType,
		Config: bo.AppViewConfig{
			Orders: []bo.AppViewOrder{
				{
					Asc:    true,
					Column: "sort",
				},
			},
			Condition: bo.AppViewLessShowCondition{
				Type: "and",
				Conds: []bo.AppViewLessShowCondition{
					{
						Column:    consts.BasicFieldPlanEndTime,
						FieldType: "datepicker",
						Type:      "gt",
						Value:     "${today}",
						Props: bo.AppViewLessProp{
							DatePicker: bo.AppViewLessPropSelectType{
								SelectType: "${today}",
							},
						},
					},
					{
						Column:    consts.BasicFieldPlanEndTime,
						FieldType: "datepicker",
						Type:      "lt",
						Value:     "${afterday:3}",
						Props: bo.AppViewLessProp{
							DatePicker: bo.AppViewLessPropSelectType{
								SelectType: "${afterday:N}",
							},
						},
					},
				},
			},
			RealCondition:   struct{}{},
			HiddenColumnIds: defaultHiddenColumnIds,
		},
	})

	return appDefaultViews
}

// CreateDefaultViewsForProject 为项目表单串创建默认的**项目视图**主要包含：进行中、已完成、已归档、全部
func CreateDefaultViewsForProject(orgId, proFormAppId, projectFolderId int64) errs.SystemErrorInfo {
	defaultConfig := []bo.LcAppView{
		{
			ViewName: "进行中",
			Type:     1,
			Config: bo.AppViewConfig{
				Orders: []bo.AppViewOrder{
					{
						Asc:    true,
						Column: "sort",
					},
				},
				Condition: bo.AppViewMultiConditionForProject{
					Type: "and",
					Conds: []bo.AppViewLessCondition{
						{
							Column: consts.ProBasicFieldParticipantIds,
							Type:   "values_in",
							Value:  "${current_dept_and_parents}",
						}, {
							Column: consts.ProBasicFieldStatus,
							Type:   "equal",
							Value:  2,
						}, {
							Column: consts.ProBasicFieldTemplateFlag,
							Type:   "equal",
							Value:  consts.TemplateFalse,
						}, {
							Column: consts.ProBasicFieldIsDelete,
							Type:   "equal",
							Value:  consts.AppIsNoDelete,
						},
					},
				},
			},
		},
		{
			ViewName: "已完成",
			Type:     1,
			Config: bo.AppViewConfig{
				Orders: []bo.AppViewOrder{
					{
						Asc:    true,
						Column: "sort",
					},
				},
				Condition: bo.AppViewMultiConditionForProject{
					Type: "and",
					Conds: []bo.AppViewLessCondition{
						{
							Column: consts.ProBasicFieldParticipantIds,
							Type:   "values_in",
							Value:  "${current_dept_and_parents}",
						}, {
							Column: consts.ProBasicFieldStatus,
							Type:   "equal",
							Value:  3,
						}, {
							Column: consts.ProBasicFieldTemplateFlag,
							Type:   "equal",
							Value:  consts.TemplateFalse,
						}, {
							Column: consts.ProBasicFieldIsDelete,
							Type:   "equal",
							Value:  consts.AppIsNoDelete,
						},
					},
				},
			},
		},
		{
			ViewName: "已归档",
			Type:     1,
			Config: bo.AppViewConfig{
				Orders: []bo.AppViewOrder{
					{
						Asc:    true,
						Column: "sort",
					},
				},
				Condition: bo.AppViewMultiConditionForProject{
					Type: "and",
					Conds: []bo.AppViewLessCondition{
						{
							Column: consts.ProBasicFieldParticipantIds,
							Type:   "values_in",
							Value:  "${current_dept_and_parents}",
						}, {
							Column:    "isFiling",
							FieldType: "select",
							Type:      "in",
							Values:    []interface{}{"1"},
						}, {
							Column: consts.ProBasicFieldTemplateFlag,
							Type:   "equal",
							Value:  consts.TemplateFalse,
						}, {
							Column: consts.ProBasicFieldIsDelete,
							Type:   "equal",
							Value:  consts.AppIsNoDelete,
						},
					},
				},
			},
		},
		{
			ViewName: "全部",
			Type:     1,
			Config: bo.AppViewConfig{
				Orders: []bo.AppViewOrder{
					{
						Asc:    true,
						Column: "sort",
					},
				},
				ProjectObjectTypeId: 0,
				Condition: bo.AppViewMultiConditionForProject{
					Type: "and",
					Conds: []bo.AppViewLessCondition{
						{
							Column: consts.ProBasicFieldParticipantIds,
							Type:   "values_in",
							Value:  "${current_dept_and_parents}",
						}, {
							Column: consts.ProBasicFieldTemplateFlag,
							Type:   "equal",
							Value:  consts.TemplateFalse,
						},
					},
				},
				RealCondition: bo.AppViewLessCondition{},
			},
		},
	}
	viewDataArr := make([]appvo.CreateAppViewData, 0)
	for _, view := range defaultConfig {
		viewDataArr = append(viewDataArr, appvo.CreateAppViewData{
			AppId:   proFormAppId,
			OrgId:   orgId,
			Config:  json.ToJsonIgnoreError(view.Config),
			Name:    view.ViewName,
			OwnerId: 0,
			Type:    view.Type,
		})
	}

	resp := appfacade.CreateAppViews(appvo.CreateAppViewReq{Reqs: viewDataArr})
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}

	return nil
}

// CreateDefaultViewsForIssue 为汇总表创建默认的视图
func CreateDefaultViewsForIssue(orgId, summaryAppId int64) errs.SystemErrorInfo {
	defaultConfig := []bo.LcAppViewJson{
		{
			ViewName: consts.MyOverDueView,
			Type:     1,
			Config:   consts.MyOverDueViewConfig,
		},
		{
			ViewName: consts.MyDailyDueOfToday,
			Type:     1,
			Config:   consts.MyDailyDueOfTodayViewConfig,
		},
		{
			ViewName: consts.MyOverDueSoonView,
			Type:     1,
			Config:   consts.OverDueSoonViewConfig,
		},
		{
			ViewName: consts.MyPending,
			Type:     1,
			Config:   consts.MyPendingViewConfig,
		},
		{
			ViewName: consts.OwnerViewName,
			Type:     1,
			Config:   consts.OwnerViewConfig,
		},
		{
			ViewName: consts.AllOverDueSoonView,
			Type:     1,
			Config:   consts.AllOverDueSoonViewConfig,
		},
		{
			ViewName: consts.AllOverDueView,
			Type:     1,
			Config:   consts.AllOverDueViewConfig,
		},
		{
			ViewName: consts.AllUnFinishedView,
			Type:     1,
			Config:   consts.AllUnFinishedViewConfig,
		},

		//{
		//	ViewName: consts.AllIssueView,
		//	Type:     1,
		//	Config:   consts.AllIssueViewConfig,
		//},
	}
	viewDataArr := make([]appvo.CreateAppViewData, 0)
	for _, view := range defaultConfig {
		ownerId := int64(0)
		if view.ViewName != consts.AllIssueView {
			ownerId = -1
		}
		viewDataArr = append(viewDataArr, appvo.CreateAppViewData{
			AppId:   summaryAppId,
			OrgId:   orgId,
			Config:  view.Config,
			Name:    view.ViewName,
			OwnerId: ownerId,
			Type:    view.Type,
		})
	}

	resp := appfacade.CreateAppViews(appvo.CreateAppViewReq{Reqs: viewDataArr})
	if resp.Failure() {
		log.Errorf("[CreateDefaultViewsForIssue] appfacade.CreateAppViews err:%v, orgId:%v, summaryAppId:%v",
			resp.Error(), orgId, summaryAppId)
		return resp.Error()
	}
	return nil
}
