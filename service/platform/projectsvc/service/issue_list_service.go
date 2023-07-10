package service

import (
	"strconv"
	"strings"

	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/encrypt"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	lang2 "github.com/star-table/polaris-backend/common/core/util/lang"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/bo/status"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/trendsvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/trendsfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

var getProcessError = "proxies.GetProcessStatusId: %v\n"

var log = *logger.GetDefaultLogger()

//func HomeIssues(orgId, currentUserId int64, page int, size int, input *vo.HomeIssueInfoReq, isFromOpen bool) (*vo.HomeIssueInfoResp, errs.SystemErrorInfo) {
//	log.Infof("[HomeIssues] 当前登录用户：%d，组织：%d", currentUserId, orgId)
//	projectBo := &bo.Project{}
//	// 获取项目信息
//	if input.ProjectID != nil && *input.ProjectID > 0 {
//		tmpProBo, err := domain.GetProject(orgId, *input.ProjectID)
//		if err != nil {
//			return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
//		}
//		projectBo = tmpProBo
//	}
//
//	issueCond, err1 := IssueCondAssembly(orgId, currentUserId, input)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueCondAssemblyError, err1)
//	}
//
//	log.Infof("首页任务列表查询条件 %v", issueCond)
//
//	if input == nil {
//		input = &vo.HomeIssueInfoReq{}
//	}
//
//	inputAppId := int64(0)
//	if input.MenuAppID != nil && *input.MenuAppID != "" {
//		menuId, err := strconv.ParseInt(*input.MenuAppID, 10, 64)
//		if err != nil {
//			log.Error(err)
//			return nil, errs.ParamError
//		}
//		inputAppId = menuId
//	}
//
//	var orderBy interface{} = consts.TcCreateTime + " desc"
//	if input.OrderType != nil {
//		orderBy = domain.IssueCondOrderBy(orgId, *input.OrderType, projectBo, input.ProjectObjectTypeID, input.IssueIds)
//	}
//
//	if input.CondOrder != nil && len(input.CondOrder) > 0 {
//		orderStr := ""
//		for _, order := range input.CondOrder {
//			sortType := "asc"
//			if !order.Asc {
//				sortType = "desc"
//			}
//			if order.Column > 0 {
//				orderStr += fmt.Sprintf("json_extract(custom_field, '$.\"%d\".value') %s,", order.Column, sortType)
//			} else {
//				switch order.Column {
//				case -1: //标题
//					orderStr += fmt.Sprintf("title %s,", sortType)
//				case -2: //编号id
//					orderStr += fmt.Sprintf("code %s,", sortType)
//				case -3: //负责人id
//					orderStr += fmt.Sprintf("owner %s,", sortType)
//				case -4: //状态
//				case -6: //截止时间
//					orderStr += fmt.Sprintf("plan_end_time %s,", sortType)
//				case -7: //优先级
//					priorities, err := domain.GetPriorityListByType(orgId, consts.PriorityTypeIssue)
//					if err != nil {
//						log.Error(err)
//						//orderBy = db.Raw("(select sort from ppm_prs_priority p where p.id = priority_id) asc, plan_end_time asc, id desc")
//					} else {
//						if order.Asc {
//							bo.SortPriorityBo(*priorities)
//						} else {
//							bo.SortDescPriorityBo(*priorities)
//						}
//						orderBySort := ""
//						for _, priority := range *priorities {
//							orderBySort += fmt.Sprintf(",%d", priority.Id)
//						}
//						orderStr += "FIELD(priority_id" + orderBySort + ")"
//					}
//				case -8: //关注人
//				case -9: //开始时间
//					orderStr += fmt.Sprintf("plan_start_time %s,", sortType)
//				case -10: //迭代
//				case -11: //任务栏
//				case -12: //需求类型
//				case -13: //需求来源
//				case -14: //缺陷类型
//				case -15: //严重程度
//				case -16: //标签
//				case -26: //确认人
//				case -28: //创建人
//				case -29: //逾期任务
//				case -30: //创建时间
//				case -31: //状态类型
//				case -100: //sort
//					orderStr += fmt.Sprintf("sort %s,", sortType)
//				}
//			}
//
//		}
//		if len(orderStr) > 0 {
//			orderBy = db.Raw(orderStr[0 : len(orderStr)-1])
//		}
//	}
//
//	if input.IsParentBeforeChid != nil && *input.IsParentBeforeChid == 1 {
//		issueCond[consts.TcParentId] = 0
//	}
//	var union []*db.Union
//	if input.SearchCond != nil && *input.SearchCond != "" {
//		union = append(union, db.Or(db.Cond{
//			consts.TcTitle: db.Like("%" + *input.SearchCond + "%"),
//		}).Or(db.Cond{
//			consts.TcCode: db.Like("%" + *input.SearchCond + "%"),
//		}))
//	}
//	if input.Code != nil && *input.Code != "" {
//		issueCond[consts.TcCode] = *input.Code
//	}
//
//	orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
//		OrgId:  orgId,
//		UserId: currentUserId,
//	})
//	if orgInfoResp.Failure() {
//		log.Error(orgInfoResp.Error())
//		return nil, orgInfoResp.Error()
//	}
//	orgRemarkObj := &orgvo.OrgRemarkConfigType{}
//	oriErr := json.FromJson(orgInfoResp.OrganizationInfo.Remark, orgRemarkObj)
//	if oriErr != nil {
//		log.Error(oriErr)
//		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
//	}
//	needAppId := inputAppId
//	if needAppId == 0 {
//		if input.ProjectID != nil && *input.ProjectID != 0 && projectBo.AppId != 0 {
//			needAppId = projectBo.AppId
//		}
//	}
//	appAuthInfo := appauth.GetAppAuthData{}
//	isProAdmin := false
//	if needAppId > 0 {
//		optAuthResp := permissionfacade.GetAppAuth(orgId, needAppId, currentUserId)
//		if optAuthResp.Failure() {
//			log.Infof("[HomeIssues] GetAppAuth orgId:%v,needAppId:%v, userId:%v, message:%v", orgId, needAppId, currentUserId, optAuthResp.Message)
//			return nil, optAuthResp.Error()
//		}
//		appAuthInfo = optAuthResp.NewData
//		isProAdmin, _ = domain.CheckIsProAdmin(orgId, projectBo.AppId, currentUserId, &appAuthInfo)
//	}
//
//	ownerIsContain, _ := slice.Contain(projectBo.OwnerIds, currentUserId)
//	isAdmin := domain.CheckIsAdmin(orgId, currentUserId) || (currentUserId == projectBo.Owner || ownerIsContain) ||
//		isProAdmin
//	if input.ProjectID != nil {
//		if !isAdmin && projectBo.TemplateFlag != consts.TemplateTrue && !isFromOpen {
//			allCollaborateIssueIds, err := domain.GetCollaborateIssues(orgId, currentUserId, orgRemarkObj.OrgSummaryTableAppId, []int64{*input.ProjectID})
//			if err != nil {
//				log.Error(err)
//				return nil, err
//			}
//			//公开项目：如果是我参与的项目展示所有的，如果不是，则只展示我协作的任务
//			isAttendFlag := false
//			if inputAppId != 0 {
//				//镜像应用（校验是否是应用成员和项目成员）
//				isMemberResp := appfacade.IsAppMember(appvo.IsAppMemberReq{
//					AppId:  inputAppId,
//					OrgId:  orgId,
//					UserId: currentUserId,
//				})
//				if isMemberResp.Failure() {
//					log.Error(isMemberResp.Error())
//					return nil, isMemberResp.Error()
//				}
//				isAttendFlag = isMemberResp.NewData
//			} else {
//				isAttend, err := domain.IsProjectParticipant(orgId, currentUserId, *input.ProjectID)
//				if err != nil {
//					log.Error(err)
//					return nil, err
//				}
//				isAttendFlag = isAttend
//			}
//			if !isAttendFlag {
//				issueCond["     "+consts.TcId] = db.In(allCollaborateIssueIds) // 加空格是防止键被覆盖
//			}
//		}
//	}
//
//	////自定义字段筛选
//	//customFieldCond(issueCond, input, &union, orgId, currentUserId)
//
//	if input.StatusList != nil && len(input.StatusList) > 0 {
//		statusUnion, unionErr := domain.IssueCondStatusListAssembly(orgId, input.StatusList)
//		if unionErr != nil {
//			log.Error(unionErr)
//			return nil, unionErr
//		}
//		union = append(union, statusUnion)
//	}
//
//	if input.TrulyStatusIds != nil && len(input.TrulyStatusIds) > 0 {
//		//通用项目需要进行判断
//		allStatus, err := domain.GetIssueAllStatus(orgId, []int64{0}, []int64{0})
//		if err != nil {
//			log.Error(err)
//			return nil, err
//		}
//		//如果有通用项目中的已完成状态（则需要判断审批状态）
//		finishStatusIds := []int64{}
//		for _, bos := range allStatus[0] {
//			if bos.Type == consts.StatusTypeComplete {
//				finishStatusIds = append(finishStatusIds, bos.ID)
//			}
//		}
//		var issueUnion db.Union
//
//		//未完成的状态不需要判断审核状态
//		notNeedOtherFilterStatus := []int64{}
//		trulyFinishStatus := []int64{} //已完成的状态（需要加上审批通过的条件）
//		for _, id := range input.TrulyStatusIds {
//			if ok, _ := slice.Contain(finishStatusIds, id); !ok {
//				if id == int64(-1) {
//					//待确认
//					issueUnion = *(issueUnion.Or(
//						db.And(db.Cond{
//							consts.TcStatus:      db.In(finishStatusIds),
//							consts.TcAuditStatus: consts.AuditStatusNotView,
//						}),
//					))
//				} else {
//					notNeedOtherFilterStatus = append(notNeedOtherFilterStatus, id)
//				}
//			} else {
//				trulyFinishStatus = append(trulyFinishStatus, id)
//			}
//		}
//		if len(notNeedOtherFilterStatus) > 0 {
//			issueUnion = *(issueUnion.Or(
//				db.Cond{consts.TcStatus: db.In(notNeedOtherFilterStatus)}))
//		}
//		if len(trulyFinishStatus) > 0 {
//			issueUnion = *(issueUnion.Or(
//				db.And(db.Cond{
//					consts.TcStatus:      db.In(trulyFinishStatus),
//					consts.TcAuditStatus: consts.AuditStatusPass,
//				}),
//			))
//		}
//		union = append(union, &issueUnion)
//	}
//
//	//是否逾期
//	if input.IsOverdue != nil {
//		//已经逾期筛选条件，并且未完成
//		nowTime := time.Now()
//		//已完成
//		finishedId, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeComplete)
//		if err != nil {
//			log.Errorf(getProcessError, err)
//			return nil, err
//		}
//		if *input.IsOverdue == 1 {
//			//逾期（实际完成时间>计划完成时间  and 未完成且已超过当前时间）
//			union = append(union, db.Or(db.Cond{
//				consts.TcPlanEndTime:       db.Lt(db.Raw(consts.TcEndTime)),
//				" " + consts.TcPlanEndTime: db.Gte(consts.BlankElasticityTime),
//			}).Or(db.And(db.Cond{
//				" " + consts.TcPlanEndTime: db.Between(consts.BlankElasticityTime, date.Format(nowTime)),
//			}, db.Cond{
//				" " + consts.TcStatus: db.NotIn(*finishedId),
//			})))
//		} else if *input.IsOverdue == 2 {
//			//未逾期（已完成且实际完成时间<=计划完成时间 and 未超过当前时间 and 没有设置计划完成时间）
//			union = append(union, db.Or(db.And(db.Cond{
//				consts.TcPlanEndTime: db.Gte(db.Raw(consts.TcEndTime)),
//			}, db.Cond{
//				consts.TcStatus: db.In(*finishedId),
//			})).Or(db.Cond{
//				consts.TcPlanEndTime: db.Gt(date.Format(nowTime)),
//			}).Or(db.Cond{
//				consts.TcPlanEndTime: db.Lt(consts.BlankElasticityTime),
//			}))
//		}
//	}
//
//	if input.IsOnlyPolaris != nil && *input.IsOnlyPolaris {
//		if input.ProjectID != nil {
//			issueCond = db.Cond{}
//			issueCond[consts.TcProjectId] = *input.ProjectID
//			issueCond[consts.TcIsDelete] = consts.AppIsNoDelete
//			issueCond[consts.TcOrgId] = orgId
//			union = []*db.Union{}
//		}
//	}
//	polarisPage := 0
//	polarisSize := 0
//	if input.IsOnlyPolaris != nil && *input.IsOnlyPolaris {
//		//如果仅仅查询极星，就直接分页。否则就用无码的数据进行分页
//		polarisPage = page
//		polarisSize = size
//	}
//	//issueBos, total, err := domain.SelectList(issueCond, union, page, size, orderBy)
//	issueBos, total, err := domain.SelectList(issueCond, union, polarisPage, polarisSize, orderBy, true)
//	if err != nil {
//		log.Error(strs.ObjectToString(err))
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
//	}
//
//	//真正满足条件的任务id
//	meetIds := []int64{}
//	for _, issueBo := range *issueBos {
//		meetIds = append(meetIds, issueBo.Id)
//	}
//
//	if input.IsParentBeforeChid != nil && *input.IsParentBeforeChid == 1 {
//		allNeedIds := []int64{}
//		for _, issueBo := range *issueBos {
//			allNeedIds = append(allNeedIds, issueBo.Id)
//		}
//		issueAndChildren, err := domain.GetIssueAndChildren(orgId, allNeedIds, input.IsFiling)
//		if err != nil {
//			log.Error(err)
//			return nil, err
//		}
//		*issueBos = issueAndChildren
//	}
//
//	var actualTotal = total
//	if input.EnableParentIssues != nil {
//		if *input.EnableParentIssues == 1 {
//			//查询子任务的父任务
//			parentIds := getIssueParentIds(*issueBos)
//			if len(parentIds) > 0 {
//				tempCond := db.And(issueCond)
//				if len(union) > 0 {
//					for _, d := range union {
//						tempCond = tempCond.And(d)
//					}
//				}
//				tempUnion := db.Or(tempCond, db.Cond{
//					consts.TcId: db.In(parentIds),
//				})
//
//				//issueBos, actualTotal, err = domain.SelectList(db.Cond{}, []*db.Union{tempUnion}, page, size, orderBy)
//				issueBos, actualTotal, err = domain.SelectList(db.Cond{}, []*db.Union{tempUnion}, polarisPage, polarisSize, orderBy, true)
//				if err != nil {
//					log.Error(strs.ObjectToString(err))
//					return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
//				}
//			}
//		} else if *input.EnableParentIssues == 2 {
//			issueIds := make([]int64, 0)
//			for _, issueBo := range *issueBos {
//				issueIds = append(issueIds, issueBo.Id)
//			}
//			//查询父任务的子任务
//			tempCond := db.And(issueCond)
//			if len(union) > 0 {
//				for _, d := range union {
//					tempCond = tempCond.And(d)
//				}
//			}
//			if len(issueIds) > 0 {
//				tempUnion := db.Or(tempCond, db.Cond{
//					consts.TcParentId: db.In(issueIds),
//				})
//
//				//issueBos, actualTotal, err = domain.SelectList(db.Cond{}, []*db.Union{tempUnion}, page, size, orderBy)
//				issueBos, actualTotal, err = domain.SelectList(db.Cond{}, []*db.Union{tempUnion}, polarisPage, polarisSize, orderBy, true)
//				if err != nil {
//					log.Error(strs.ObjectToString(err))
//					return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
//				}
//			}
//		}
//	}
//
//	log.Infof("首页任务列表命中总条数 %d", total)
//
//	allNeedIssueBos := []bo.IssueBo{}
//	trulyTotal := actualTotal
//	if input.IsOnlyPolaris != nil && *input.IsOnlyPolaris {
//		allNeedIssueBos = *issueBos
//	} else {
//		//根据所有符合条件的id去查询无码
//		var allUsefulIssueIds []interface{}
//		issueInfoMap := map[int64]bo.IssueBo{}
//		for _, issueBo := range *issueBos {
//			allUsefulIssueIds = append(allUsefulIssueIds, issueBo.Id)
//			issueInfoMap[issueBo.Id] = issueBo
//		}
//		lessReq := vo.LessCondsData{}
//		if input.LessConds != nil {
//			lessReq = *input.LessConds
//			lessReq.Conds = append(lessReq.Conds, &vo.LessCondsData{
//				Type:   "in",
//				Values: allUsefulIssueIds,
//				Column: "issueId",
//				Left:   nil,
//				Right:  nil,
//				Conds:  nil,
//			})
//		} else {
//			lessReq = vo.LessCondsData{
//				Type:   "in",
//				Values: allUsefulIssueIds,
//				Column: "issueId",
//				Left:   nil,
//				Right:  nil,
//				Conds:  nil,
//			}
//		}
//
//		if input.ProjectID != nil && *input.ProjectID != 0 {
//			projectInfo, err := domain.GetProject(orgId, *input.ProjectID)
//			if err != nil {
//				log.Error(err)
//				return nil, err
//			}
//
//			if projectInfo.ProjectTypeId == consts.ProjectTypeNormalId {
//				lessReq = businees.ConvertIssueStatusFilterReqForAll(lessReq)
//			}
//		} else {
//			if input.ProjectID == nil {
//				lessReq = domain.ConvertIssueStatusFilterReqForAll(orgId, lessReq)
//			}
//		}
//
//		orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
//			OrgId:  orgId,
//			UserId: currentUserId,
//		})
//		if orgInfoResp.Failure() {
//			log.Error(orgInfoResp.Error())
//			return nil, orgInfoResp.Error()
//		}
//		orgRemarkObj := &orgvo.OrgRemarkConfigType{}
//		oriErr := json.FromJson(orgInfoResp.OrganizationInfo.Remark, orgRemarkObj)
//		if oriErr != nil {
//			log.Error(oriErr)
//			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
//		}
//
//		lessReqParam := formvo.LessIssueListReq{
//			Condition: lessReq,
//			AppId:     orgRemarkObj.OrgSummaryTableAppId,
//			OrgId:     orgId,
//			UserId:    currentUserId,
//			Page:      int64(page),
//			Size:      int64(size),
//		}
//		if input.LessOrder != nil {
//			lessReqParam.Orders = input.LessOrder
//		}
//		if input.TableID != nil {
//			tableId, tableIdErr := strconv.ParseInt(*input.TableID, 10, 64)
//			if tableIdErr != nil {
//				log.Errorf("[HomeIssues] tableId err:%v. orgId:%d, userId:%d", tableIdErr, orgId, currentUserId)
//				return nil, errs.InvalidTableId
//			}
//			lessReqParam.TableId = tableId
//		}
//		//暂时只在单项目查询的时候传入，因为综合查询是不需要展示自定义字段的
//		if input.ProjectID != nil {
//			appId, err := domain.GetAppIdFromProjectId(orgId, *input.ProjectID)
//			if err != nil {
//				log.Error(err)
//				return nil, err
//			}
//			lessReqParam.RedirectIds = []int64{appId}
//		}
//		lessResp := formfacade.LessIssueList(lessReqParam)
//		if lessResp.Failure() {
//			log.Error(lessResp.Error())
//			return nil, lessResp.Error()
//		}
//		//真正符合条件的数据
//		lessIssueListOrder := make([]bo.IssueBo, 0)
//		selectedIssuesFormLc := map[int64]bo.IssueBo{}
//		for _, m := range lessResp.NewData.List {
//			i, ok := m["issueId"]
//			if !ok {
//				continue
//			}
//			issueId := int64(0)
//			if id, ok := i.(int64); ok {
//				issueId = id
//			} else if id, ok1 := i.(int); ok1 {
//				issueId = int64(id)
//			} else if id, ok1 := i.(float64); ok1 {
//				//理论上 map 解析json 会将int转为float64
//				issueIdStr := strconv.FormatFloat(id, 'f', -1, 64)
//				parseId, err := strconv.ParseInt(issueIdStr, 10, 64)
//				if err != nil {
//					log.Error(err)
//					continue
//				} else {
//					issueId = parseId
//				}
//			} else {
//				continue
//			}
//
//			if issueBo, ok := issueInfoMap[issueId]; ok {
//				issueBo.LessData = m
//				selectedIssuesFormLc[issueBo.Id] = issueBo
//				lessIssueListOrder = append(lessIssueListOrder, issueBo)
//			}
//		}
//		trulyTotal = lessResp.NewData.Total
//
//		//排序处理
//		if input.LessOrder != nil && len(input.LessOrder) > 0 {
//			//有关于无码的排序
//			allNeedIssueBos = lessIssueListOrder
//		} else {
//			//只有极星的排序
//			for _, id := range allUsefulIssueIds {
//				if issueBo, ok := selectedIssuesFormLc[id.(int64)]; ok {
//					allNeedIssueBos = append(allNeedIssueBos, issueBo)
//				}
//			}
//		}
//	}
//
//	homeIssueBos, err3 := domain.ConvertIssueBosToHomeIssueInfos(orgId, allNeedIssueBos)
//	if err3 != nil {
//		log.Error(err3)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
//	}
//
//	homeIssueVos := &[]*vo.HomeIssueInfo{}
//
//	err2 := copyer.Copy(homeIssueBos, homeIssueVos)
//	if err2 != nil {
//		log.Error(err2)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err2)
//	}
//	if input.IsOnlyPolaris != nil && *input.IsOnlyPolaris {
//		//同步历史数据的时候用到任务详情
//		for i, issueBo := range homeIssueBos {
//			remark := issueBo.Issue.IssueDetailBo.Remark
//			(*homeIssueVos)[i].Issue.Remark = &remark
//		}
//	}
//
//	if len(*homeIssueVos) > 0 {
//		for i, info := range *homeIssueVos {
//			if ok, _ := slice.Contain(meetIds, info.Issue.ParentID); ok {
//				(*homeIssueVos)[i].ParentIsMeetCondition = 1
//			} else {
//				(*homeIssueVos)[i].ParentIsMeetCondition = 0
//			}
//
//			if ok, _ := slice.Contain(meetIds, info.Issue.ID); ok {
//				(*homeIssueVos)[i].IsAttach = 0
//			} else {
//				(*homeIssueVos)[i].IsAttach = 1
//			}
//
//			//状态做下特殊逻辑处理 (已完成和待确认的处理状况)
//			if info.Project.ProjectTypeID != consts.ProjectTypeAgileId {
//				if len(info.AuditorsInfo) > 0 && info.Issue.AuditStatus == consts.AuditStatusNotView && info.Status.Type == consts.StatusTypeComplete {
//					//一切都是在有确认人的情况下
//					//如果未审核完成,并且处于完成状态 所有状态里面没有已完成状态
//					for j, status := range info.AllStatus {
//						if status.Name == english.WordTransLate("已完成") {
//							(*homeIssueVos)[i].AllStatus[j] = &vo.HomeIssueStatusInfo{
//								ID:          status.ID,
//								Name:        english.WordTransLate(consts.WaitConfirmStatusName),
//								DisplayName: status.DisplayName,
//								BgStyle:     consts.WaitConfirmStatusBgStyle,
//								FontStyle:   consts.WaitConfirmStatusFontStyle,
//								Type:        status.Type,
//								Sort:        status.Sort,
//							}
//						}
//					}
//
//					//如果处于完成状态 并且 还未确认
//					(*homeIssueVos)[i].Status = &vo.HomeIssueStatusInfo{
//						ID:          info.Status.ID,
//						Name:        english.WordTransLate(consts.WaitConfirmStatusName),
//						DisplayName: info.Status.DisplayName,
//						BgStyle:     consts.WaitConfirmStatusBgStyle,
//						FontStyle:   consts.WaitConfirmStatusFontStyle,
//						Type:        info.Status.Type,
//						Sort:        info.Status.Sort,
//					}
//				}
//			}
//		}
//	}
//	if !isFromOpen && needAppId > 0 {
//		// 鉴权时，检查是否是系统管理员
//		if !appAuthInfo.HasAppRootPermission {
//			for i, info := range *homeIssueVos {
//				for s, _ := range info.LessData {
//					if !appAuthInfo.HasFieldViewAuth(s) {
//						delete((*homeIssueVos)[i].LessData, s)
//					}
//				}
//			}
//		}
//	}
//	resp := &vo.HomeIssueInfoResp{
//		//Total:       total,
//		Total:       int64(trulyTotal),
//		ActualTotal: actualTotal,
//		List:        *homeIssueVos,
//	}
//	return resp, nil
//}

//func IssueReport(orgId, currentUserId int64, reportType int64) (*vo.IssueReportResp, errs.SystemErrorInfo) {
//
//	log.Infof("[IssueReport] 当前登录用户 %d 组织 %d", currentUserId, orgId)
//
//	issueCond := db.Cond{}
//	issueCond[consts.TcIsDelete] = consts.AppIsNoDelete
//	issueCond[consts.TcOrgId] = orgId
//
//	//获取我参与的和我负责的任务
//	issueRelations, issueErr := domain.GetRelatedIssues(currentUserId, orgId)
//	if issueErr != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, issueErr)
//	}
//	if len(issueRelations) > 0 {
//		issueRelationIds := make([]int64, len(issueRelations))
//		for i, entity := range issueRelations {
//			issueRelationIds[i] = entity.IssueId
//		}
//		issueCond[consts.TcId] = db.In(issueRelationIds)
//	} else {
//		issueCond[consts.TcId] = -1
//	}
//
//	//时间范围
//	var startTime, endTime string
//	now := time.Now()
//	switch reportType {
//	case consts.DailyReport:
//		startTime = date.Format(date.GetZeroTime(now))
//		endTime = date.Format((date.GetZeroTime(now)).AddDate(0, 0, 1).Add(-1 * time.Second))
//	case consts.WeeklyReport:
//		startTime = date.Format(date.GetWeekStart(now))
//		endTime = date.Format((date.GetWeekStart(now)).AddDate(0, 0, 7).Add(-1 * time.Second))
//	case consts.MonthlyReport:
//		startTime = date.Format(date.GetMonthStart(now))
//		endTime = date.Format((date.GetMonthStart(now)).AddDate(0, 1, 0).Add(-1 * time.Second))
//	}
//
//	//进行中
//	processingIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeRunning)
//	if err != nil {
//		log.Errorf(getProcessError, err)
//		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//	}
//	//已完成
//	finishedId, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeComplete)
//	if err != nil {
//		log.Errorf(getProcessError, err)
//		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//	}
//
//	var union []*db.Union
//	//获取进行中和在当前时间段内完成的任务
//	union = append(union, db.Or(db.Cond{
//		consts.TcStatus: db.In(processingIds),
//	}).Or(
//		db.And(
//			db.Cond{consts.TcStatus: db.In(finishedId)},
//			db.Cond{consts.TcEndTime: db.Gte(startTime)},
//		),
//	))
//
//	log.Infof("任务分享列表查询条件 %v", issueCond)
//
//	orderBy := consts.TcPlanStartTime + " desc"
//
//	issueBos, total, err := domain.SelectList(issueCond, union, -1, -1, orderBy, true)
//	if err != nil {
//		log.Error(strs.ObjectToString(err))
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
//	}
//	homeIssueBos, err := domain.ConvertIssueBosToHomeIssueInfos(orgId, *issueBos)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
//	}
//
//	resp, insertErr := domain.InsertIssueReport(orgId, currentUserId, total, startTime, endTime, homeIssueBos)
//	if insertErr != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, insertErr)
//	}
//
//	result := &vo.IssueReportResp{}
//
//	copyErr := copyer.Copy(resp, result)
//	if copyErr != nil {
//		log.Error(copyErr)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
//	}
//	return result, nil
//}

func IssueReportDetail(shareID string) (*vo.IssueReportResp, errs.SystemErrorInfo) {
	id, err := encrypt.AesDecrypt(shareID)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.DecryptError, err)
	}

	shareInfo, err := domain.GetIssueReport(id)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	homeIssue := &vo.IssueReportResp{}
	err = json.FromJson(*shareInfo.Content, homeIssue)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
	}

	return homeIssue, nil
}

//func IssueCondAssembly(orgId, currentUserId int64, input *vo.HomeIssueInfoReq) (db.Cond, errs.SystemErrorInfo) {
//	issueCond := db.Cond{}
//	issueCond[consts.TcIsDelete] = consts.AppIsNoDelete
//	issueCond[consts.TcOrgId] = orgId
//	if input == nil {
//		input = &vo.HomeIssueInfoReq{}
//	}
//
//	//封装基础条件
//	IssueCondBaseAssembly(issueCond, input)
//
//	//审批相关
//	//确认人
//	if input.AuditorIds != nil && len(input.AuditorIds) > 0 {
//		issueCond[consts.TcId] = db.In(db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ? and relation_id in ?", orgId, consts.IssueRelationTypeAuditor, input.AuditorIds))
//	}
//
//	//前置列表封装
//	if input.IssueIDForBefore != nil && *input.IssueIDForBefore != 0 {
//		pos := &[]po.PpmPriIssueRelation{}
//		err := mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
//			consts.TcIsDelete:     consts.AppIsNoDelete,
//			consts.TcOrgId:        orgId,
//			consts.TcRelationId:   *input.IssueIDForBefore,
//			consts.TcRelationType: consts.IssueRelationTypeBeforeAfter,
//		}, pos)
//		if err != nil {
//			log.Error(err)
//			return nil, errs.MysqlOperateError
//		}
//		allIssueIds := []int64{*input.IssueIDForBefore}
//		for _, relation := range *pos {
//			allIssueIds = append(allIssueIds, relation.IssueId)
//		}
//		issueCond["  "+consts.TcId] = db.NotIn(allIssueIds)
//	}
//	//后置列表封装
//	if input.IssueIDForAfter != nil && *input.IssueIDForAfter != 0 {
//		pos := &[]po.PpmPriIssueRelation{}
//		err := mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
//			consts.TcIsDelete:     consts.AppIsNoDelete,
//			consts.TcOrgId:        orgId,
//			consts.TcIssueId:      *input.IssueIDForAfter,
//			consts.TcRelationType: consts.IssueRelationTypeBeforeAfter,
//		}, pos)
//		if err != nil {
//			log.Error(err)
//			return nil, errs.MysqlOperateError
//		}
//		allIssueIds := []int64{*input.IssueIDForAfter}
//		for _, relation := range *pos {
//			allIssueIds = append(allIssueIds, relation.RelationId)
//		}
//		issueCond["   "+consts.TcId] = db.NotIn(allIssueIds)
//	}
//
//	if input.RelatedType != nil {
//		domain.IssueCondRelatedTypeAssembly(issueCond, *input.RelatedType, currentUserId, orgId)
//	} else {
//		// 暂时不过滤私有项目，因为还有协作人的情况。
//		//isAdmin := false
//		//if currentUserId != 0 {
//		//	manageAuthInfoResp := userfacade.GetUserAuthority(orgId, currentUserId)
//		//	if manageAuthInfoResp.Failure() {
//		//		log.Error(manageAuthInfoResp.Message)
//		//		return nil, manageAuthInfoResp.Error()
//		//	}
//		//	adminFlag := manageAuthInfoResp.NewData
//		//	if adminFlag.IsSysAdmin || adminFlag.IsSubAdmin {
//		//		isAdmin = true
//		//	}
//		//} else {
//		//	isAdmin = true
//		//}
//
//		//会对私有项目做过滤，私有项目除外
//		// domain.IssueCondNoRelatedTypeAssembly(issueCond, currentUserId, orgId, isAdmin)
//	}
//
//	if input.ProjectID != nil {
//		issueCond[consts.TcProjectId+" "] = *input.ProjectID
//	}
//
//	if input.IsFiling == nil || *input.IsFiling == 0 || *input.IsFiling > 3 {
//		//默认查询未归档的项目
//		defaultFiling := consts.ProjectIsNotFiling
//		input.IsFiling = &defaultFiling
//	}
//	domain.IssueCondFiling(issueCond, orgId, *input.IsFiling)
//
//	if len(input.IssueTagID) != 0 {
//		domain.IssueCondTagId(issueCond, orgId, input.IssueTagID)
//	}
//
//	if input.ResourceID != nil {
//		domain.IssueCondResourceId(issueCond, orgId, *input.ResourceID)
//	}
//
//	//if input.Status != nil {
//	//	err1 := domain.IssueCondStatusAssembly(issueCond, orgId, *input.Status)
//	//	if err1 != nil {
//	//		log.Error(err1)
//	//		return nil, errs.BuildSystemErrorInfo(errs.IssueCondAssemblyError, err1)
//	//	}
//	//}
//	domain.IssueCondRelationMemberAssembly(issueCond, input)
//
//	//组合筛选类型封装
//	//err2 := domain.IssueCondCombinedCondAssembly(issueCond, input, currentUserId, orgId)
//	//if err2 != nil {
//	//	log.Error(err2)
//	//	return nil, err2
//	//}
//	//增量查询条件封装
//	domain.IssueCondLastUpdateTimeCondAssembly(issueCond, input)
//
//	// 当前任务（用于变更父任务时查询任务列表）
//	if input.CurrentIssueID != nil && *input.CurrentIssueID != 0 {
//		issueCond["   "+consts.TcPath] = db.NotLike("%," + str.ToString(*input.CurrentIssueID) + ",%")
//	}
//
//	return issueCond, nil
//}

//func customFieldCond(queryCond db.Cond, input *vo.HomeIssueInfoReq, union *[]*db.Union, orgId int64, currentUserId int64) {
//	for i, cond := range input.Conds {
//		blank := ""
//		for x := 0; x <= i; x++ {
//			blank += " "
//		}
//		var valueType string
//		if cond.Value != nil {
//			valueType = reflect.TypeOf(cond.Value).Kind().String()
//		}
//		if cond.Column < 0 {
//			//处理通用字段
//			switch cond.Column {
//			case -1: //标题
//				if cond.Value != nil && valueType == "string" {
//					if cond.Type == "like" {
//						queryCond[blank+consts.TcTitle] = db.Like("%" + fmt.Sprintf("%s", cond.Value) + "%")
//					} else if cond.Type == "not_like" {
//						queryCond[blank+consts.TcTitle] = db.NotLike("%" + fmt.Sprintf("%s", cond.Value) + "%")
//					}
//				}
//			case -2: //标号id
//				if cond.Value != nil && valueType == "string" {
//					if cond.Type == "like" {
//						queryCond[blank+consts.TcCode] = db.Like("%" + fmt.Sprintf("%s", cond.Value) + "%")
//					} else if cond.Type == "not_like" {
//						queryCond[blank+consts.TcCode] = db.NotLike("%" + fmt.Sprintf("%s", cond.Value) + "%")
//					}
//				}
//			case -3: //负责人
//				if cond.Value != nil && valueType == "slice" {
//					raw := db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ? and relation_id in ?", orgId, consts.IssueRelationTypeOwner, cond.Value)
//					if cond.Type == "in" {
//						queryCond[blank+consts.TcId+"   "] = db.In(raw)
//					} else if cond.Type == "not_in" {
//						queryCond[blank+consts.TcId+"    "] = db.NotIn(raw)
//					} else if cond.Type == "values_in" {
//						queryCond[blank+consts.TcId+"   "] = db.In(raw)
//					} else if cond.Type == "all_in" {
//						//目前成员all_in只支持单个
//						tempVal := []interface{}{}
//						_ = copyer.Copy(cond.Value, &tempVal)
//						if len(tempVal) > 0 {
//							queryCond[blank+" "+consts.TcOwner] = tempVal[0]
//						}
//					}
//				}
//				if cond.Type == "is_null" {
//					queryCond[blank+" "+consts.TcOwner] = 0
//				} else if cond.Type == "not_null" {
//					queryCond[blank+" "+consts.TcOwner] = db.NotEq(0)
//				} else if cond.Type == "equal" && cond.Value != nil {
//					if valueType == "string" {
//						if cond.Value == consts.KeyCurrentUser {
//							queryCond[blank+" "+consts.TcOwner] = currentUserId
//						}
//					} else {
//						queryCond[blank+" "+consts.TcOwner] = cond.Value
//					}
//				}
//			case -4: //状态
//				//待确认（通用项目有待确认状态-1）
//				if cond.Value != nil && valueType == "slice" {
//					//通用项目需要进行判断
//					allStatus, err := domain.GetIssueAllStatus(orgId, []int64{0}, []int64{0})
//					if err != nil {
//						log.Error(err)
//						return
//					}
//					//如果有通用项目中的已完成状态（则需要判断审批状态）
//					finishStatusIds := []int64{}
//					for _, bos := range allStatus[0] {
//						if bos.Type == consts.StatusTypeComplete {
//							finishStatusIds = append(finishStatusIds, bos.ID)
//						}
//					}
//					searchStatusIds := &[]string{}
//					copyErr := copyer.Copy(cond.Value, searchStatusIds)
//					if copyErr != nil {
//						log.Error(copyErr)
//						return
//					}
//					if cond.Type == "in" {
//						var issueUnion db.Union
//						//未完成的状态不需要判断审核状态
//						notNeedOtherFilterStatus := []int64{}
//						trulyFinishStatus := []int64{} //已完成的状态（需要加上审批通过的条件）
//						for _, idStr := range *searchStatusIds {
//							id, _ := strconv.ParseInt(idStr, 10, 64)
//							if ok, _ := slice.Contain(finishStatusIds, id); !ok {
//								if id == int64(-1) {
//									//待确认
//									issueUnion = *(issueUnion.Or(
//										db.And(db.Cond{
//											consts.TcStatus:      db.In(finishStatusIds),
//											consts.TcAuditStatus: consts.AuditStatusNotView,
//										}),
//									))
//								} else {
//									notNeedOtherFilterStatus = append(notNeedOtherFilterStatus, id)
//								}
//							} else {
//								trulyFinishStatus = append(trulyFinishStatus, id)
//							}
//						}
//						if len(notNeedOtherFilterStatus) > 0 {
//							issueUnion = *(issueUnion.Or(
//								db.Cond{consts.TcStatus: db.In(notNeedOtherFilterStatus)}))
//						}
//						if len(trulyFinishStatus) > 0 {
//							issueUnion = *(issueUnion.Or(
//								db.And(db.Cond{
//									consts.TcStatus:      db.In(trulyFinishStatus),
//									consts.TcAuditStatus: consts.AuditStatusPass,
//								}),
//							))
//						}
//						*union = append(*union, &issueUnion)
//					} else if cond.Type == "not_in" {
//						notNeedOtherFilterStatus := []int64{}
//						trulyFinishStatus := []int64{} //已完成的状态（需要加上审批通过的条件）
//						for _, idStr := range *searchStatusIds {
//							id, _ := strconv.ParseInt(idStr, 10, 64)
//							if ok, _ := slice.Contain(finishStatusIds, id); !ok {
//								if id == int64(-1) {
//									//待确认
//									//queryCond[db.Raw(blank+"case when status in ? and audit_status != ? then 1 else 0", finishStatusIds, consts.AuditStatusPass)] = 0
//									//这里是not_in 所以要取反
//									queryCond[db.Raw(blank+"(case when status in (?) and audit_status != ? then 1 else 0 end)", strings.Replace(strings.Trim(fmt.Sprint(finishStatusIds), "[]"), " ", ",", -1), consts.AuditStatusPass)] = 0
//								} else {
//									notNeedOtherFilterStatus = append(notNeedOtherFilterStatus, id)
//								}
//							} else {
//								trulyFinishStatus = append(trulyFinishStatus, id)
//							}
//						}
//						if len(notNeedOtherFilterStatus) > 0 {
//							queryCond[blank+consts.TcStatus] = db.NotIn(notNeedOtherFilterStatus)
//						}
//						if len(trulyFinishStatus) > 0 {
//							//queryCond[db.Raw(blank+"(case when status in ? and audit_status = ? then 1 else 0 end)", finishStatusIds, consts.AuditStatusPass)] = 0
//							//这里是not_in 所以要取反
//							queryCond[db.Raw(blank+"(case when status in (?) and audit_status = ? then 1 else 0 end)", strings.Replace(strings.Trim(fmt.Sprint(trulyFinishStatus), "[]"), " ", ",", -1), consts.AuditStatusPass)] = 0
//						}
//					}
//				}
//			case -6: //截止时间
//				//比较的是“天”，所以要区别判断具体
//				begin, _ := time.Parse(consts.AppTimeFormat, fmt.Sprintf("%v", cond.Value))
//				end := begin.AddDate(0, 0, 1)
//				if cond.Type == "gte" {
//					if cond.Value != nil {
//						//大于等于今天零点
//						queryCond[blank+consts.TcPlanEndTime] = db.Gte(cond.Value)
//					}
//				} else if cond.Type == "lte" {
//					if cond.Value != nil {
//						//小于明天零点
//						queryCond[blank+consts.TcPlanEndTime] = db.Between(consts.BlankElasticityTime, end.Add(-1*time.Second).String())
//					}
//				} else if cond.Type == "is_null" {
//					queryCond[blank+consts.TcPlanEndTime] = db.Lt(consts.BlankElasticityTime)
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcPlanEndTime] = db.Gte(consts.BlankElasticityTime)
//				} else if cond.Type == "gt" {
//					if cond.Value != nil {
//						//大于等于明天零点
//						queryCond[blank+consts.TcPlanEndTime] = db.Gte(end.String())
//					}
//				} else if cond.Type == "lt" {
//					if cond.Value != nil {
//						//小于今天零点
//						queryCond[blank+consts.TcPlanEndTime] = db.Between(consts.BlankElasticityTime, begin.Add(-1*time.Second).String())
//					}
//				} else if cond.Type == "between" {
//					if cond.Left != nil && cond.Right != nil {
//						right, _ := time.Parse(consts.AppTimeFormat, fmt.Sprintf("%v", cond.Right))
//						queryCond[blank+consts.TcPlanEndTime] = db.Between(cond.Left, right.Add((86400-1)*time.Second).String())
//					}
//				}
//			case -7: //优先级
//				if cond.Value != nil && valueType == "slice" {
//					if cond.Type == "in" {
//						queryCond[blank+consts.TcPriorityId] = db.In(cond.Value)
//					} else if cond.Type == "not_in" {
//						queryCond[blank+consts.TcPriorityId] = db.NotIn(cond.Value)
//					}
//				}
//			case -8: //关注人
//				if cond.Value != nil && valueType == "slice" {
//					raw := db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ? and relation_id in ?", orgId, consts.IssueRelationTypeFollower, cond.Value)
//					if cond.Type == "in" {
//						queryCond[blank+consts.TcId+"     "] = db.In(raw)
//					} else if cond.Type == "not_in" {
//						queryCond[blank+consts.TcId+"      "] = db.NotIn(raw)
//					} else if cond.Type == "values_in" {
//						queryCond[blank+consts.TcId+"     "] = db.In(raw)
//					} else if cond.Type == "all_in" {
//						//目前成员all_in只支持单个
//						tempVal := []interface{}{}
//						_ = copyer.Copy(cond.Value, &tempVal)
//						if len(tempVal) > 0 {
//							queryCond[blank+consts.TcId+"     "] = db.In(db.Raw("select issue_id from (select issue_id, count(*) as total from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ? and relation_id in ? group by issue_id) as temp where temp.total = 1", orgId, consts.IssueRelationTypeFollower, tempVal))
//						}
//					}
//				}
//				raw := db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ?", orgId, consts.IssueRelationTypeFollower)
//				if cond.Type == "is_null" {
//					queryCond[blank+consts.TcId+"       "] = db.NotIn(raw)
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcId+"       "] = db.In(raw)
//				} else if cond.Type == "equal" && cond.Value != nil {
//					if valueType == "string" {
//						if cond.Value == consts.KeyCurrentUser {
//							queryCond[blank+consts.TcId+"       "] = db.In(db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ? and relation_id = ?", orgId, consts.IssueRelationTypeFollower, currentUserId))
//						}
//					} else {
//						queryCond[blank+consts.TcId+"       "] = db.In(db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ? and relation_id = ?", orgId, consts.IssueRelationTypeFollower, cond.Value))
//					}
//				}
//			case -9: //开始时间
//				//比较的是“天”，所以要区别判断具体
//				begin, _ := time.Parse(consts.AppTimeFormat, fmt.Sprintf("%v", cond.Value))
//				end := begin.AddDate(0, 0, 1)
//				if cond.Type == "gte" {
//					if cond.Value != nil {
//						//大于等于今天零点
//						queryCond[blank+consts.TcPlanStartTime] = db.Gte(cond.Value)
//					}
//				} else if cond.Type == "lte" {
//					if cond.Value != nil {
//						//小于明天零点
//						queryCond[blank+consts.TcPlanStartTime] = db.Between(consts.BlankElasticityTime, end.Add(-1*time.Second).String())
//					}
//				} else if cond.Type == "is_null" {
//					queryCond[blank+consts.TcPlanStartTime] = db.Lt(consts.BlankElasticityTime)
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcPlanStartTime] = db.Gte(consts.BlankElasticityTime)
//				} else if cond.Type == "gt" {
//					if cond.Value != nil {
//						//大于等于明天零点
//						queryCond[blank+consts.TcPlanStartTime] = db.Gte(end.String())
//					}
//				} else if cond.Type == "lt" {
//					if cond.Value != nil {
//						//小于今天零点
//						queryCond[blank+consts.TcPlanStartTime] = db.Between(consts.BlankElasticityTime, begin.Add(-1*time.Second).String())
//					}
//				} else if cond.Type == "between" {
//					if cond.Left != nil && cond.Right != nil {
//						right, _ := time.Parse(consts.AppTimeFormat, fmt.Sprintf("%v", cond.Right))
//						queryCond[blank+consts.TcPlanStartTime] = db.Between(cond.Left, right.Add((86400-1)*time.Second).String())
//					}
//				}
//			case -10: //迭代
//				if cond.Value != nil {
//					if valueType == "slice" {
//						if cond.Type == "in" {
//							queryCond[blank+" "+consts.TcIterationId] = db.In(cond.Value)
//						} else if cond.Type == "not_in" {
//							queryCond[blank+" "+consts.TcIterationId] = db.NotIn(cond.Value)
//						}
//					} else if cond.Type == "equal" {
//						queryCond[blank+" "+consts.TcIterationId] = cond.Value
//					}
//				}
//			case -11: //任务栏
//				//if cond.Value != nil {
//				//	if valueType == "slice" {
//				//		if cond.Type == "in" {
//				//			queryCond[blank+" "+consts.TcProjectObjectTypeId] = db.In(cond.Value)
//				//		} else if cond.Type == "not_in" {
//				//			queryCond[blank+" "+consts.TcProjectObjectTypeId] = db.NotIn(cond.Value)
//				//		}
//				//	} else if cond.Type == "equal" {
//				//		queryCond[blank+" "+consts.TcProjectObjectTypeId] = cond.Value
//				//	}
//				//}
//				//新版改成到issue_relation里面去查询
//				if cond.Value != nil && valueType == "slice" {
//					raw := db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ? and relation_id in ?", orgId, consts.IssueRelationTypeBelongManyPro, cond.Value)
//					if cond.Type == "in" {
//						queryCond[blank+consts.TcId+"      "] = db.In(raw)
//					} else if cond.Type == "not_in" {
//						queryCond[blank+consts.TcId+"       "] = db.NotIn(raw)
//					} else if cond.Type == "values_in" {
//						queryCond[blank+consts.TcId+"      "] = db.In(raw)
//					}
//				}
//				raw := db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ?", orgId, consts.IssueRelationTypeBelongManyPro)
//				if cond.Type == "is_null" {
//					queryCond[blank+consts.TcId+"        "] = db.NotIn(raw)
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcId+"       "] = db.In(raw)
//				} else if cond.Type == "equal" && cond.Value != nil {
//					if valueType == "string" {
//						if cond.Value == consts.KeyCurrentUser {
//							queryCond[blank+consts.TcId+"        "] = db.In(db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ? and relation_id = ?", orgId, consts.IssueRelationTypeBelongManyPro, currentUserId))
//						}
//					} else {
//						queryCond[blank+consts.TcId+"        "] = db.In(db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ? and relation_id = ?", orgId, consts.IssueRelationTypeBelongManyPro, cond.Value))
//					}
//				}
//			case -12: //需求类型
//				if cond.Value != nil && valueType == "slice" {
//					if cond.Type == "in" {
//						queryCond[blank+consts.TcIssueObjectTypeId] = db.In(cond.Value)
//					} else if cond.Type == "not_in" {
//						queryCond[blank+consts.TcIssueObjectTypeId] = db.NotIn(cond.Value)
//					}
//				}
//				if cond.Type == "is_null" {
//					queryCond[blank+consts.TcIssueObjectTypeId] = 0
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcIssueObjectTypeId] = db.NotEq(0)
//				}
//			case -13: //需求来源
//				if cond.Value != nil && valueType == "slice" {
//					if cond.Type == "in" {
//						queryCond[blank+consts.TcSourceId] = db.In(cond.Value)
//					} else if cond.Type == "not_in" {
//						queryCond[blank+consts.TcSourceId] = db.NotIn(cond.Value)
//					}
//				}
//				if cond.Type == "is_null" {
//					queryCond[blank+consts.TcSourceId] = 0
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcSourceId] = db.NotEq(0)
//				}
//			case -14: //缺陷类型
//				if cond.Value != nil && valueType == "slice" {
//					if cond.Type == "in" {
//						queryCond[blank+consts.TcIssueObjectTypeId] = db.In(cond.Value)
//					} else if cond.Type == "not_in" {
//						queryCond[blank+consts.TcIssueObjectTypeId] = db.NotIn(cond.Value)
//					}
//				}
//				if cond.Type == "is_null" {
//					queryCond[blank+consts.TcIssueObjectTypeId] = 0
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcIssueObjectTypeId] = db.NotEq(0)
//				}
//			case -15: //严重程度
//				if cond.Value != nil && valueType == "slice" {
//					if cond.Type == "in" {
//						queryCond[blank+consts.TcPropertyId] = db.In(cond.Value)
//					} else if cond.Type == "not_in" {
//						queryCond[blank+consts.TcPropertyId] = db.NotIn(cond.Value)
//					}
//				}
//				if cond.Type == "is_null" {
//					queryCond[blank+consts.TcPropertyId] = 0
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcPropertyId] = db.NotEq(0)
//				}
//			case -16: //标签
//				if cond.Value != nil && valueType == "slice" {
//					raw := db.Raw("select distinct issue_id from ppm_pri_issue_tag where org_id = ? and is_delete = 2 and tag_id in ?", orgId, cond.Value)
//					if cond.Type == "in" {
//						queryCond[blank+consts.TcId+"         "] = db.In(raw)
//					} else if cond.Type == "not_in" {
//						queryCond[blank+consts.TcId+"          "] = db.NotIn(raw)
//					}
//				}
//
//				raw := db.Raw("select distinct issue_id from ppm_pri_issue_tag where org_id = ? and is_delete = 2", orgId)
//				if cond.Type == "is_null" {
//					queryCond[blank+consts.TcId+"         "] = db.NotIn(raw)
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcId+"          "] = db.In(raw)
//				}
//			case -26: //确认人
//				if cond.Value != nil && valueType == "slice" {
//					raw := db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ? and relation_id in ?", orgId, consts.IssueRelationTypeAuditor, cond.Value)
//					if cond.Type == "in" {
//						queryCond[blank+consts.TcId+"       "] = db.In(raw)
//					} else if cond.Type == "not_in" {
//						queryCond[blank+consts.TcId+"       "] = db.NotIn(raw)
//					} else if cond.Type == "values_in" {
//						queryCond[blank+consts.TcId+"       "] = db.In(raw)
//					} else if cond.Type == "all_in" {
//						//目前成员all_in只支持单个
//						tempVal := []interface{}{}
//						_ = copyer.Copy(cond.Value, &tempVal)
//						if len(tempVal) > 0 {
//							queryCond[blank+consts.TcId+"     "] = db.In(db.Raw("select issue_id from (select issue_id, count(*) as total from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ? and relation_id in ? group by issue_id) as temp where temp.total = 1", orgId, consts.IssueRelationTypeAuditor, tempVal))
//						}
//					}
//				}
//				raw := db.Raw("select distinct issue_id from ppm_pri_issue_relation where org_id = ? and is_delete= 2 and relation_type = ?", orgId, consts.IssueRelationTypeAuditor)
//				if cond.Type == "is_null" {
//					queryCond[blank+consts.TcId+"       "] = db.NotIn(raw)
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcId+"       "] = db.In(raw)
//				}
//			case -28: //创建人
//				if cond.Value != nil && valueType == "slice" {
//					if cond.Type == "in" {
//						queryCond[blank+consts.TcCreator] = db.In(cond.Value)
//					} else if cond.Type == "not_in" {
//						queryCond[blank+consts.TcCreator] = db.NotIn(cond.Value)
//					} else if cond.Type == "values_in" {
//						queryCond[blank+consts.TcCreator] = db.In(cond.Value)
//					}
//				} else if cond.Type == "equal" && cond.Value != nil {
//					if valueType == "string" {
//						if cond.Value == consts.KeyCurrentUser {
//							queryCond[blank+consts.TcCreator] = currentUserId
//						}
//					} else {
//						queryCond[blank+consts.TcCreator] = cond.Value
//					}
//				}
//				if cond.Type == "is_null" {
//					queryCond[blank+consts.TcCreator] = 0
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcCreator] = db.NotEq(0)
//				}
//			case -29: //逾期
//				if cond.Value != nil {
//					nowTime := time.Now()
//					//已完成
//					finishedId, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeComplete)
//					if err != nil {
//						log.Errorf(getProcessError, err)
//						return
//					}
//					if valueType == "string" {
//						if cond.Value == "1" {
//							//逾期（实际完成时间>计划完成时间  and 未完成且已超过当前时间）
//							*union = append(*union, db.Or(db.Cond{
//								consts.TcPlanEndTime:       db.Lt(db.Raw(consts.TcEndTime)),
//								" " + consts.TcPlanEndTime: db.Gte(consts.BlankElasticityTime),
//							}).Or(db.And(db.Cond{
//								" " + consts.TcPlanEndTime: db.Between(consts.BlankElasticityTime, date.Format(nowTime)),
//							}, db.Cond{
//								" " + consts.TcStatus: db.NotIn(*finishedId),
//							})))
//						} else if cond.Value == "2" {
//							//未逾期（已完成且实际完成时间<=计划完成时间 and 未超过当前时间 and 没有设置计划完成时间）
//							*union = append(*union, db.Or(db.And(db.Cond{
//								consts.TcPlanEndTime: db.Gte(db.Raw(consts.TcEndTime)),
//							}, db.Cond{
//								consts.TcStatus: db.In(*finishedId),
//							})).Or(db.Cond{
//								consts.TcPlanEndTime: db.Gt(date.Format(nowTime)),
//							}).Or(db.Cond{
//								consts.TcPlanEndTime: db.Lt(consts.BlankElasticityTime),
//							}))
//						}
//					}
//				}
//			case -30: //创建时间
//				//比较的是“天”，所以要区别判断具体
//				begin, _ := time.Parse(consts.AppTimeFormat, fmt.Sprintf("%v", cond.Value))
//				end := begin.AddDate(0, 0, 1)
//				if cond.Type == "gte" {
//					if cond.Value != nil {
//						//大于等于今天零点
//						queryCond[blank+consts.TcCreateTime] = db.Gte(cond.Value)
//					}
//				} else if cond.Type == "lte" {
//					if cond.Value != nil {
//						//小于明天零点
//						queryCond[blank+consts.TcCreateTime] = db.Between(consts.BlankElasticityTime, end.Add(-1*time.Second).String())
//					}
//				} else if cond.Type == "is_null" {
//					queryCond[blank+consts.TcCreateTime] = db.Lt(consts.BlankElasticityTime)
//				} else if cond.Type == "not_null" {
//					queryCond[blank+consts.TcCreateTime] = db.Gte(consts.BlankElasticityTime)
//				} else if cond.Type == "gt" {
//					if cond.Value != nil {
//						//大于等于明天零点
//						queryCond[blank+consts.TcCreateTime] = db.Gte(end.String())
//					}
//				} else if cond.Type == "lt" {
//					if cond.Value != nil {
//						//小于今天零点
//						queryCond[blank+consts.TcCreateTime] = db.Between(consts.BlankElasticityTime, begin.Add(-1*time.Second).String())
//					}
//				} else if cond.Type == "between" {
//					if cond.Left != nil && cond.Right != nil {
//						right, _ := time.Parse(consts.AppTimeFormat, fmt.Sprintf("%v", cond.Right))
//						queryCond[blank+consts.TcCreateTime] = db.Between(cond.Left, right.Add((86400-1)*time.Second).String())
//					}
//				}
//			case -31: //状态类型
//				if cond.Value != nil && valueType == "slice" && cond.Type == "in" {
//					var searchStatusIds []int
//					copyErr := copyer.Copy(cond.Value, &searchStatusIds)
//					if copyErr != nil {
//						log.Error(copyErr)
//						return
//					}
//					statusUnion, unionErr := domain.IssueCondStatusListAssembly(orgId, searchStatusIds)
//					if unionErr != nil {
//						log.Error(unionErr)
//						return
//					}
//					*union = append(*union, statusUnion)
//				}
//			}
//		} else {
//			fieldInfo, err := domain.GetCustomFieldInfo(orgId, cond.Column)
//			if err != nil {
//				log.Error(err)
//				return
//			}
//			switch cond.Type {
//			case "between":
//				if cond.Left != nil && cond.Right != nil {
//					fmtStr := "json_extract(custom_field, '$.\"%d\".value')"
//					if fieldInfo.FieldType == consts.CustomTypeDate {
//						// 对于日期时间的比较，需要使用 JSON_UNQUOTE 包装一下。参考：https://blog.csdn.net/qq_36527154/article/details/109778744
//						fmtStr = blank + " JSON_UNQUOTE(json_extract(custom_field, '$.\"%d\".value')) "
//					}
//					queryCond[db.Raw(fmt.Sprintf(fmtStr, cond.Column))] = db.Between(cond.Left, cond.Right)
//				}
//			case "equal":
//				if cond.Value != nil {
//					queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = cond.Value
//				}
//			case "un_equal":
//				if cond.Value != nil {
//					queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = db.NotEq(cond.Value)
//				}
//			case "gt":
//				if cond.Value != nil {
//					queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = db.Gt(cond.Value)
//				}
//			case "gte":
//				if cond.Value != nil {
//					queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = db.Gte(cond.Value)
//				}
//			case "lt":
//				if cond.Value != nil {
//					queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = db.Lt(cond.Value)
//				}
//			case "lte":
//				if cond.Value != nil {
//					queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = db.Lte(cond.Value)
//				}
//			case "like":
//				if cond.Value != nil && valueType == "string" {
//					queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = db.Like("%" + fmt.Sprintf("%s", cond.Value) + "%")
//				}
//			case "not_like":
//				if cond.Value != nil && valueType == "string" {
//					queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = db.NotLike("%" + fmt.Sprintf("%s", cond.Value) + "%")
//				}
//			case "is_null":
//				//queryCond[db.Raw(fmt.Sprintf(blank + "json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = db.IsNull()
//				//queryCond[db.Raw(fmt.Sprintf(blank + "json_extract(custom_field, '$.\"%d\".value') ", cond.Column))] = db.Eq("")
//
//				*union = append(*union, db.Or(db.Cond{
//					db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column)): db.IsNull(),
//				}, db.Cond{
//					db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value') ", cond.Column)): "",
//				}))
//			case "not_null":
//				queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = db.IsNotNull()
//				queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value') ", cond.Column))] = db.NotEq("")
//			case "in":
//				if cond.Value != nil && valueType == "slice" {
//					if fieldInfo.FieldType == consts.CustomTypePersonnel {
//						s := []int64{}
//						_ = copyer.Copy(cond.Value, &s)
//						var curUnion db.Union
//						for _, i := range s {
//							curUnion = *(curUnion.Or(db.Cond{
//								db.Raw(fmt.Sprintf("json_contains(custom_field, '{\"userId\": %v}', '$.\"%d\".value')", i, cond.Column)): 1,
//							}))
//						}
//						*union = append(*union, &curUnion)
//					} else if fieldInfo.FieldType == consts.CustomTypeRadio {
//						queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = db.In(cond.Value)
//					} else if fieldInfo.FieldType == consts.CustomTypeCheckBox {
//						s := []string{}
//						_ = copyer.Copy(cond.Value, &s)
//						var curUnion db.Union
//						for _, v := range s {
//							curUnion = *(curUnion.Or(db.Cond{
//								db.Raw(fmt.Sprintf("json_contains(custom_field, '\"%v\"', '$.\"%d\".value')", v, cond.Column)): 1,
//							}))
//						}
//						*union = append(*union, &curUnion)
//					}
//				}
//			case "not_in":
//				if cond.Value != nil && valueType == "slice" {
//					if fieldInfo.FieldType == consts.CustomTypeRadio {
//						queryCond[db.Raw(fmt.Sprintf(blank+"json_extract(custom_field, '$.\"%d\".value')", cond.Column))] = db.NotIn(cond.Value)
//					}
//				}
//			case "all_in":
//				if cond.Value != nil && valueType == "slice" {
//					valueStr := json.ToJsonIgnoreError(cond.Value)
//					length := len(valueStr)
//					if length > 2 {
//						valueStr = valueStr[1 : length-1]
//						queryCond[db.Raw(fmt.Sprintf(blank+"json_contains(json_array(%s), json_extract(custom_field, '$.\"%d\".value'))", valueStr, cond.Column))] = 1
//					}
//				}
//			case "values_in":
//				if cond.Value != nil && valueType == "slice" {
//					s := reflect.ValueOf(cond.Value).Elem()
//					var curUnion db.Union
//					for i := 0; i < s.Len(); i++ {
//						if fieldInfo.FieldType == consts.CustomTypePersonnel {
//							curUnion = *(curUnion.Or(db.Cond{
//								db.Raw(fmt.Sprintf("json_contains(custom_field, '{\"userId\": %v}', '$.\"%d\".value')", s.Index(i), cond.Column)): 1,
//							}))
//						} else {
//							curUnion = *(curUnion.Or(db.Cond{
//								db.Raw(fmt.Sprintf("json_contains(custom_field, '\"%s\"', '$.\"%d\".value')", s.Index(i), cond.Column)): 1,
//							}))
//						}
//					}
//					*union = append(*union, &curUnion)
//				}
//			}
//		}
//	}
//}

func getIssueParentIds(issues []bo.IssueBo) []int64 {
	////子任务map, 子任务id -> 父任务id
	//childMap := map[int64]int64{}
	//for _, issue := range issues {
	//	childMap[issue.Id] = issue.ParentId
	//}
	//missingParentIds := make([]int64, 0)
	//for _, v := range childMap {
	//	if v <= 0 {
	//		continue
	//	}
	//	if _, ok := childMap[v]; !ok {
	//		missingParentIds = append(missingParentIds, v)
	//	}
	//}

	allIds := []int64{}
	for _, issue := range issues {
		allIds = append(allIds, issue.Id)
		if issue.ParentId != 0 {
			for _, s := range strings.Split(issue.Path[0:len(issue.Path)-1], ",") {
				id, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					log.Error(err)
					break
				}
				if id != 0 {
					allIds = append(allIds, id)
				}
			}
		}
	}
	return slice.SliceUniqueInt64(allIds)
}

// HandleGroupChatAtUserName 项目群聊-查询用户在项目内的任务统计信息。响应用户指令：@用户名a
func HandleGroupChatAtUserName(input *projectvo.HandleGroupChatUserInsAtUserNameReq) (bool, errs.SystemErrorInfo) {
	orgId, projectId, err := domain.GetProjectIdByOpenChatId(input.OpenChatId)
	if err != nil {
		log.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.ParamError, err)
	}
	// 查询操作人信息
	resp1 := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: input.OpUserOpenId,
	})
	if resp1.Failure() {
		return false, resp1.Error()
	}
	// 查询用户信息
	resp := orgfacade.GetBaseUserInfoByEmpId(orgvo.GetBaseUserInfoByEmpIdReqVo{
		OrgId: orgId,
		EmpId: input.AtUserOpenId,
	})
	if resp.Failure() {
		log.Errorf("[HandleGroupChatAtUserName] err: %v, orgId: %d, ourOrgId: %s", err, orgId,
			resp.BaseUserInfo.OutOrgId)
		return false, resp.Error()
	}
	atUser := resp.BaseUserInfo
	// 机器人将信息回复给用户
	//tenant, err := feishu.GetTenant(resp.BaseUserInfo.OutOrgId)
	//if err != nil {
	//	log.Errorf("[HandleGroupChatAtUserName] err: %v, ourOrgId: %s", err, resp.BaseUserInfo.OutOrgId)
	//	return false, err
	//}
	// 查询负责人是被“at的用户”的任务统计信息

	// 检查用户是否是项目成员
	isIn, err := domain.IsProjectParticipant(orgId, atUser.UserId, projectId)
	if err != nil {
		return false, err
	}

	// 如果不是项目成员，则返回提示信息
	if !isIn {
		return true, nil
	}
	//fsResp, oriErr := tenant.SendMessage(fsvo.MsgVo{
	//	ChatId:  input.OpenChatId,
	//	MsgType: "interactive",
	//	Card:    msgCard,
	//})
	//if oriErr != nil {
	//	log.Error(oriErr)
	//	return false, err
	//}
	//if fsResp.Code != 0 {
	//	log.Error("HandleGroupChatAtUserName 发送消息异常")
	//}
	//log.Infof("飞书群聊-用户指令处理-HandleGroupChatAtUserName-响应用户调用结果 code: %v", fsResp.Code)

	return true, nil
}

// 项目群聊，响应用户指令：@用户名a 任务标题1。表示创建负责人为 `用户名a`，标题为 `任务标题1` 的任务。
func HandleGroupChatAtUserNameWithIssueTitle(input *projectvo.HandleGroupChatUserInsAtUserNameWithIssueTitleReq) (bool, errs.SystemErrorInfo) {
	return true, nil
}

func getProjectFirstTableId(orgId, projectId int64) (int64, errs.SystemErrorInfo) {
	appId, err := domain.GetAppIdFromProjectId(orgId, projectId)
	if err != nil {
		return 0, err
	}
	tables, err := domain.GetAppTableList(orgId, appId)
	if err != nil {
		return 0, err
	}
	if len(tables) > 0 {
		return tables[0].TableId, nil
	}

	return 0, nil
}

// 项目群聊，响应用户指令，项目进展
func HandleGroupChatUserInsProProgress(input *projectvo.HandleGroupChatUserInsProProgressReq) (bool, errs.SystemErrorInfo) {
	return true, nil
}

// HandleGroupChatUserInsProjectSettings 项目群聊指令，响应用户指令，群聊推送设置
func HandleGroupChatUserInsProjectSettings(openChatId string, sourceChannel string) (bool, errs.SystemErrorInfo) {
	return true, nil
}

// 项目群聊，响应用户指令，项目任务，通过 chat id 等参数，获取对应项目下任务的统计信息
func HandleGroupChatUserInsProIssue(input *projectvo.HandleGroupChatUserInsProIssueReq) (bool, errs.SystemErrorInfo) {
	return true, nil
}

func IssueStatusTypeStat(orgId, currentUserId int64, input *vo.IssueStatusTypeStatReq) (*vo.IssueStatusTypeStatResp, errs.SystemErrorInfo) {
	if input.ProjectID != nil && !domain.JudgeProjectIsExist(orgId, *input.ProjectID) {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}
	if input.IterationID != nil {
		iterationBo, err := domain.GetIterationBoByOrgId(*input.IterationID, orgId)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.IterationNotExist)
		}
		input.ProjectID = &iterationBo.ProjectId
	}
	// 兼容一下查询项目下的各个状态任务的数量，而非只查询某项目下某个成员的任务统计。
	issueStatusStatBos, err1 := domain.GetIssueStatusStatWithLc(bo.IssueStatusStatCondBo{
		OrgId:        orgId,
		ProjectId:    input.ProjectID,
		IterationId:  input.IterationID,
		RelationType: input.RelationType,
		UserId:       currentUserId,
	})
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}
	result := &vo.IssueStatusTypeStatResp{}

	for _, issueStatusStatBo := range issueStatusStatBos {
		result.NotStartTotal += int64(issueStatusStatBo.IssueWaitCount)
		result.ProcessingTotal += int64(issueStatusStatBo.IssueRunningCount)
		result.CompletedTotal += int64(issueStatusStatBo.IssueEndCount)
		result.CompletedTodayTotal += int64(issueStatusStatBo.IssueEndTodayCount)
		result.Total += int64(issueStatusStatBo.IssueCount)
		result.OverdueCompletedTotal += int64(issueStatusStatBo.IssueOverdueEndCount)
		result.OverdueTodayTotal += int64(issueStatusStatBo.IssueOverdueTodayCount)
		result.OverdueTotal += int64(issueStatusStatBo.IssueOverdueCount)
		result.OverdueTomorrowTotal += int64(issueStatusStatBo.IssueOverdueTomorrowCount)
		result.TodayCount += int64(issueStatusStatBo.TodayCount)
		result.TodayCreateCount += int64(issueStatusStatBo.TodayCreateCount)
		result.WaitConfirmedTotal += int64(issueStatusStatBo.WaitConfirmedCount)
	}

	//即将到期 今天到期的主子任务数+明日逾期的主子任务数
	result.BeAboutToOverdueSum = result.OverdueTomorrowTotal + result.OverdueTodayTotal

	//@我的数量
	if currentUserId > 0 {
		callMeCountResp := trendsfacade.CallMeCount(trendsvo.CallMeCountReqVo{
			ProjectId: input.ProjectID,
			UserId:    currentUserId,
			OrgId:     orgId,
		})
		if callMeCountResp.Failure() {
			log.Error(callMeCountResp.Error())
			return nil, callMeCountResp.Error()
		}
		result.CallMeTotal = callMeCountResp.Count
	}

	result.List = append(result.List, &vo.StatCommon{
		Name:  "已逾期",
		Count: result.OverdueTotal,
	})
	result.List = append(result.List, &vo.StatCommon{
		Name:  "进行中",
		Count: result.ProcessingTotal,
	})
	result.List = append(result.List, &vo.StatCommon{
		Name:  "未完成",
		Count: result.NotStartTotal + result.ProcessingTotal,
	})
	result.List = append(result.List, &vo.StatCommon{
		Name:  "已完成",
		Count: result.CompletedTotal,
	})
	lang := lang2.GetLang()
	isOtherLang := lang2.IsEnglish()
	if isOtherLang {
		for index, item := range result.List {
			if tmpMap, ok1 := consts.LANG_ISSUE_STAT_DESC_MAP[lang]; ok1 {
				if tmpVal, ok2 := tmpMap[item.Name]; ok2 {
					(*result.List[index]).Name = tmpVal
				}
			}
		}
	}

	return result, nil
}

func IssueStatusTypeStatDetail(orgId, currentUserId int64, input *vo.IssueStatusTypeStatReq) (*vo.IssueStatusTypeStatDetailResp, errs.SystemErrorInfo) {
	projectId := input.ProjectID

	if projectId != nil && !domain.JudgeProjectIsExist(orgId, *projectId) {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}

	issueStatusStatBos, err1 := domain.GetIssueStatusStatWithLc(bo.IssueStatusStatCondBo{
		OrgId:        orgId,
		UserId:       currentUserId,
		ProjectId:    input.ProjectID,
		IterationId:  input.IterationID,
		RelationType: input.RelationType,
	})
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	result := &vo.IssueStatusTypeStatDetailResp{
		NotStart:   []*vo.IssueStatByObjectType{},
		Processing: []*vo.IssueStatByObjectType{},
		Completed:  []*vo.IssueStatByObjectType{},
	}

	for _, issueStatusStatBo := range issueStatusStatBos {
		temp := issueStatusStatBo
		if issueStatusStatBo.IssueWaitCount > 0 {
			result.NotStart = append(result.NotStart, &vo.IssueStatByObjectType{
				ProjectObjectTypeID:   &temp.ProjectTypeId,
				ProjectObjectTypeName: &temp.ProjectTypeName,
				Total:                 int64(temp.IssueWaitCount),
			})
		}
		if issueStatusStatBo.IssueRunningCount > 0 {
			result.Processing = append(result.Processing, &vo.IssueStatByObjectType{
				ProjectObjectTypeID:   &temp.ProjectTypeId,
				ProjectObjectTypeName: &temp.ProjectTypeName,
				Total:                 int64(temp.IssueRunningCount),
			})
		}
		if issueStatusStatBo.IssueEndCount > 0 {
			result.Completed = append(result.Completed, &vo.IssueStatByObjectType{
				ProjectObjectTypeID:   &temp.ProjectTypeId,
				ProjectObjectTypeName: &temp.ProjectTypeName,
				Total:                 int64(temp.IssueEndCount),
			})
		}
	}

	return result, nil
}

//func GetSimpleIssueInfoBatch(orgId int64, ids []int64) (*[]vo.Issue, errs.SystemErrorInfo) {
//	list, _, err := domain.SelectList(db.Cond{
//		consts.TcOrgId: orgId,
//		consts.TcId:    db.In(ids),
//		//consts.TcIsDelete: consts.AppIsNoDelete,
//	}, nil, 0, 0, nil, false)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	issueVo := &[]vo.Issue{}
//	copyErr := copyer.Copy(list, issueVo)
//	if copyErr != nil {
//		log.Error(copyErr)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
//	}
//
//	return issueVo, nil
//}

func GetLcIssueInfoBatch(orgId int64, issueIds []int64) ([]*bo.IssueBo, errs.SystemErrorInfo) {
	issueBos, err := domain.GetIssueInfosLc(orgId, 0, issueIds)
	if err != nil {
		log.Errorf("[GetSimpleIssueInfoBatch] domain.GetIssueInfosLc err:%v, orgId:%v, issueIds:%v", err, orgId, issueIds)
		return nil, err
	}
	return issueBos, nil
}

//func GetIssueRemindInfoList(reqVo projectvo.GetIssueRemindInfoListReqVo) (*projectvo.GetIssueRemindInfoListRespData, errs.SystemErrorInfo) {
//	if reqVo.Page < 0 {
//		return nil, errs.BuildSystemErrorInfo(errs.PageInvalidError)
//	}
//	if reqVo.Size < 0 || reqVo.Size > 100 {
//		return nil, errs.BuildSystemErrorInfo(errs.PageSizeInvalidError)
//	}
//
//	selectIssueIdsCondBo := bo.SelectIssueIdsCondBo{}
//
//	input := reqVo.Input
//	//计划结束时间条件
//	selectIssueIdsCondBo.BeforePlanEndTime = input.BeforePlanEndTime
//	selectIssueIdsCondBo.AfterPlanEndTime = input.AfterPlanEndTime
//	selectIssueIdsCondBo.BeforePlanStartTime = input.BeforePlanStartTime
//	selectIssueIdsCondBo.AfterPlanStartTime = input.AfterPlanStartTime
//
//	issueRemindInfos, total, err := domain.SelectIssueRemindInfoList(selectIssueIdsCondBo, reqVo.Page, reqVo.Size)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	return &projectvo.GetIssueRemindInfoListRespData{
//		Total: total,
//		List:  issueRemindInfos,
//	}, nil
//}

func IssueListStat(orgId, userId, projectId int64) (*vo.IssueListStatResp, errs.SystemErrorInfo) {
	projectInfo, projectErr := domain.GetProject(orgId, projectId)
	if projectErr != nil {
		log.Error(projectErr)
		return nil, errs.ProjectNotExist
	}
	projectTypeList, err := domain.GetAppTableList(orgId, projectInfo.AppId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	data, err := domain.GetIssueStatusStatWithLc(bo.IssueStatusStatCondBo{
		OrgId:     orgId,
		UserId:    userId,
		ProjectId: &projectId,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	midStat := map[int64]vo.IssueListStatData{}
	for _, datum := range data {
		midStat[datum.ProjectTypeId] = vo.IssueListStatData{
			Total:                 int64(datum.IssueCount),
			FinishedCount:         int64(datum.IssueEndCount),
			OverdueCount:          int64(datum.IssueOverdueCount),
			ProjectObjectTypeID:   datum.ProjectTypeId,
			ProjectObjectTypeName: datum.ProjectTypeName,
		}
	}

	resStat := []*vo.IssueListStatData{}

	for _, objectType := range projectTypeList {
		if _, ok := midStat[objectType.TableId]; ok {
			mid := midStat[objectType.TableId]
			resStat = append(resStat, &mid)
		} else {
			resStat = append(resStat, &vo.IssueListStatData{
				ProjectObjectTypeID:   objectType.TableId,
				ProjectObjectTypeName: objectType.Name,
				Total:                 0,
				FinishedCount:         0,
				OverdueCount:          0,
			})
		}
	}

	return &vo.IssueListStatResp{List: resStat}, nil
}

//func HomeIssuesGroup(orgId, currentUserId int64, page int, size int, input *vo.HomeIssueInfoReq) (*vo.HomeIssueInfoGroupResp, errs.SystemErrorInfo) {
//	if input.ProjectID == nil || *input.ProjectID == 0 {
//		return nil, errs.ParamError
//	}
//	projectId := *input.ProjectID
//	projectInfo, err := domain.GetProject(orgId, projectId)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	list, err := HomeIssuesForTableMode(orgId, currentUserId, page, size, input, false)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	timeSpan := int64(0)
//	if len(list.List) > 0 {
//		var timeRange []time.Time
//		for _, info := range list.List {
//			start := time.Time(info.Issue.PlanStartTime)
//			end := time.Time(info.Issue.PlanEndTime)
//
//			if start.After(consts.BlankTimeObject) && end.After(consts.BlankTimeObject) {
//				timeRange = append(timeRange, start)
//				timeRange = append(timeRange, end)
//			}
//		}
//
//		if len(timeRange) > 0 {
//			minTime := timeRange[0]
//			maxTime := timeRange[0]
//			for _, t := range timeRange {
//				if t.Before(minTime) {
//					minTime = t
//				}
//				if t.After(maxTime) {
//					maxTime = t
//				}
//			}
//
//			timeSpan = maxTime.Unix() - minTime.Unix()
//		}
//	}
//
//	commonRes := &vo.HomeIssueInfoGroupResp{
//		Total:       list.Total,
//		ActualTotal: list.ActualTotal,
//		TimeSpan:    timeSpan,
//		Group: []*vo.HomeIssueGroup{
//			&vo.HomeIssueGroup{
//				ID:        0,
//				Name:      "",
//				Avatar:    "",
//				BgStyle:   "",
//				FontStyle: "",
//				TimeSpan:  0,
//				List:      list.List,
//			},
//		},
//	}
//	if input.GroupType == nil {
//		return commonRes, nil
//	}
//	lang := lang2.GetLang()
//	isOtherLang := lang2.IsEnglish()
//	//任务栏
//	//projectTypeList, err := ProjectObjectTypesWithProject(orgId, projectId)
//	//if err != nil {
//	//	log.Error(err)
//	//	return nil, err
//	//}
//	trulyProjectObjectTypeList := []vo.ProjectObjectType{}
//	//projectObjectTypeIds := []int64{}
//	//for _, objectType := range projectTypeList.List {
//	//	if projectInfo.ProjectTypeId == consts.ProjectTypeAgileId && objectType.Name == "迭代" {
//	//		continue
//	//	}
//	//	projectObjectTypeIds = append(projectObjectTypeIds, objectType.ID)
//	//	temp := objectType
//	//	trulyProjectObjectTypeList = append(trulyProjectObjectTypeList, *temp)
//	//}
//	//if input.ProjectObjectTypeID != nil && *input.ProjectObjectTypeID != 0 {
//	//	projectObjectTypeIds = []int64{*input.ProjectObjectTypeID}
//	//}
//	//获取项目所有的状态
//	tableList, tableListErr := domain.GetAppTableList(orgId, projectInfo.AppId)
//	if tableListErr != nil {
//		log.Errorf("[HomeIssuesGroup] GetAppTableList failed:%d, orgId:%d, appId:%d", tableListErr, orgId, projectInfo.AppId)
//		return nil, tableListErr
//	}
//	tableIds := []int64{}
//	for _, typeBo := range tableList {
//		tableIds = append(tableIds, typeBo.Id)
//	}
//	//allStatus, err := domain.GetIssueAllStatusNew(orgId, []int64{projectId}, tableIds) //这里每个表中的statusId可能会一样，下面的分组需要注意
//	//if err != nil {
//	//	log.Error(err)
//	//	return nil, err
//	//}
//	//statusMap := map[int][]int64{}
//	//for _, bos := range allStatus {
//	//	for _, infoBo := range bos {
//	//		if temp, ok := statusMap[infoBo.Type]; ok {
//	//			if ok1, _ := slice.Contain(temp, infoBo.ID); !ok1 {
//	//				statusMap[infoBo.Type] = append(statusMap[infoBo.Type], infoBo.ID)
//	//			}
//	//		} else {
//	//			statusMap[infoBo.Type] = append(statusMap[infoBo.Type], infoBo.ID)
//	//		}
//	//	}
//	//}
//	group := []*vo.HomeIssueGroup{}
//	switch *input.GroupType {
//	case 1:
//		//获取所有任务的负责人
//		//userIds := []int64{}
//		//for _, info := range list.List {
//		//	userIds = append(userIds, info.Issue.Owner)
//		//}
//		//userIds = slice.SliceUniqueInt64(userIds)
//		//ownerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed("", orgId, userIds)
//		//if err != nil {
//		//	log.Error(err)
//		//	return nil, err
//		//}
//		//userMap := maps.NewMap("UserId", ownerInfos)
//		//for _, id := range userIds {
//		//	if _, ok := userMap[id]; ok {
//		//		temp := userMap[id].(bo.BaseUserInfoBo)
//		//		group = append(group, &vo.HomeIssueGroup{
//		//			ID:       temp.UserId,
//		//			Name:     temp.Name,
//		//			Avatar:   temp.Avatar,
//		//			TimeSpan: 0,
//		//			List:     []*vo.HomeIssueInfo{},
//		//		})
//		//	}
//		//}
//	case 2:
//		//状态
//		//if len(statusMap) == 3 {
//		//	group = append(group, &vo.HomeIssueGroup{
//		//		ID:        1,
//		//		Name:      "未开始",
//		//		Avatar:    "#FFFFFF",
//		//		BgStyle:   "#DBDBDB",
//		//		FontStyle: "",
//		//		TimeSpan:  0,
//		//		List:      []*vo.HomeIssueInfo{},
//		//	})
//		//	//敏捷包含进行中
//		//	group = append(group, &vo.HomeIssueGroup{
//		//		ID:        2,
//		//		Name:      "进行中",
//		//		Avatar:    "",
//		//		BgStyle:   "#FFCD1C",
//		//		FontStyle: "#FFFFFF",
//		//		TimeSpan:  0,
//		//		List:      []*vo.HomeIssueInfo{},
//		//	})
//		//} else {
//		//	group = append(group, &vo.HomeIssueGroup{
//		//		ID:        1,
//		//		Name:      "未完成",
//		//		Avatar:    "#FFFFFF",
//		//		BgStyle:   "#DBDBDB",
//		//		FontStyle: "",
//		//		TimeSpan:  0,
//		//		List:      []*vo.HomeIssueInfo{},
//		//	})
//		//}
//		group = append(group, &vo.HomeIssueGroup{
//			ID:        3,
//			Name:      "已完成",
//			Avatar:    "",
//			BgStyle:   "#69A922",
//			FontStyle: "#FFFFFF",
//			TimeSpan:  0,
//			List:      []*vo.HomeIssueInfo{},
//		})
//		if isOtherLang {
//			otherLanguageMap := make(map[string]string, 0)
//			if tmpMap, ok1 := consts.LANG_ISSUE_STAT_DESC_MAP[lang]; ok1 {
//				otherLanguageMap = tmpMap
//			}
//			for index, item := range group {
//				if tmpVal, ok2 := otherLanguageMap[item.Name]; ok2 {
//					group[index].Name = tmpVal
//				}
//			}
//		}
//	case 3:
//		//优先级
//		//priorities, err := domain.GetPriorityListByType(orgId, consts.PriorityTypeIssue)
//		//if err != nil {
//		//	log.Error(err)
//		//	return nil, err
//		//}
//		//bo.SortPriorityBo(*priorities)
//		//for _, priorityBo := range *priorities {
//		//	group = append(group, &vo.HomeIssueGroup{
//		//		ID:        priorityBo.Id,
//		//		Name:      priorityBo.Name,
//		//		Avatar:    "",
//		//		BgStyle:   priorityBo.BgStyle,
//		//		FontStyle: priorityBo.FontStyle,
//		//		TimeSpan:  0,
//		//		List:      []*vo.HomeIssueInfo{},
//		//	})
//		//}
//	case 4:
//		//任务栏
//		for _, objectType := range trulyProjectObjectTypeList {
//			group = append(group, &vo.HomeIssueGroup{
//				ID:        objectType.ID,
//				Name:      objectType.Name,
//				Avatar:    "",
//				BgStyle:   objectType.BgStyle,
//				FontStyle: objectType.FontStyle,
//				TimeSpan:  0,
//				List:      []*vo.HomeIssueInfo{},
//			})
//		}
//	case 5:
//		//迭代
//		iterationList, _, err := domain.GetIterationBoList(0, 0, db.Cond{
//			consts.TcIsDelete:  consts.AppIsNoDelete,
//			consts.TcProjectId: projectId,
//			consts.TcOrgId:     orgId,
//		}, nil)
//		if err != nil {
//			log.Error(err)
//			return nil, err
//		}
//		for _, iterationBo := range *iterationList {
//			group = append(group, &vo.HomeIssueGroup{
//				ID:        iterationBo.Id,
//				Name:      iterationBo.Name,
//				Avatar:    "",
//				BgStyle:   "",
//				FontStyle: "",
//				TimeSpan:  0,
//				List:      []*vo.HomeIssueInfo{},
//			})
//		}
//		// 迭代的"未规划"多语言
//		notPlanName := "未规划"
//		if isOtherLang {
//			otherLanguageMap := make(map[string]string, 0)
//			if tmpMap, ok1 := consts.LANG_ROLE_NAME_MAP[lang]; ok1 {
//				otherLanguageMap = tmpMap
//			}
//			if tmpVal, ok2 := otherLanguageMap[notPlanName]; ok2 {
//				notPlanName = tmpVal
//			}
//		}
//		group = append(group, &vo.HomeIssueGroup{
//			ID:        0,
//			Name:      notPlanName,
//			Avatar:    "",
//			BgStyle:   "",
//			FontStyle: "",
//			TimeSpan:  0,
//			List:      []*vo.HomeIssueInfo{},
//		})
//	case 6:
//		//具体状态
//		//statusArr := []int64{}
//		//for _, bos := range allStatus {
//		//	for _, infoBo := range bos {
//		//		//去重
//		//		if ok, _ := slice.Contain(statusArr, infoBo.ID); ok {
//		//			continue
//		//		}
//		//		group = append(group, &vo.HomeIssueGroup{
//		//			ID:        infoBo.ID,
//		//			Name:      infoBo.Name,
//		//			Avatar:    "",
//		//			BgStyle:   infoBo.BgStyle,
//		//			FontStyle: infoBo.FontStyle,
//		//			TimeSpan:  0,
//		//			List:      []*vo.HomeIssueInfo{},
//		//		})
//		//		statusArr = append(statusArr, infoBo.ID)
//		//	}
//		//}
//	}
//	if len(group) == 0 {
//		return commonRes, nil
//	}
//
//	for i, issueGroup := range group {
//		for _, info := range list.List {
//			switch *input.GroupType {
//			case 1:
//				//负责人
//				//if info.Issue.Owner == issueGroup.ID {
//				//	group[i].List = append(group[i].List, info)
//				//}
//			case 2:
//				//状态
//				//if status, ok := statusMap[int(issueGroup.ID)]; ok {
//				//	if ok1, _ := slice.Contain(status, info.Issue.Status); ok1 {
//				//		group[i].List = append(group[i].List, info)
//				//	}
//				//}
//			case 3:
//				//优先级
//				if info.Issue.PriorityID == issueGroup.ID {
//					group[i].List = append(group[i].List, info)
//				}
//			case 4:
//				//任务栏
//				if info.Issue.ProjectObjectTypeID == issueGroup.ID {
//					group[i].List = append(group[i].List, info)
//				}
//			case 5:
//				//迭代
//				if info.Issue.IterationID == issueGroup.ID {
//					group[i].List = append(group[i].List, info)
//				}
//			case 6:
//				//具体状态
//				if info.Issue.Status == issueGroup.ID {
//					group[i].List = append(group[i].List, info)
//				}
//			}
//		}
//	}
//
//	for i, issueGroup := range group {
//		var timeRange []time.Time
//		for _, info := range issueGroup.List {
//			start := time.Time(info.Issue.PlanStartTime)
//			end := time.Time(info.Issue.PlanEndTime)
//
//			if start.After(consts.BlankTimeObject) && end.After(consts.BlankTimeObject) {
//				timeRange = append(timeRange, start)
//				timeRange = append(timeRange, end)
//				group[i].FitTotal++
//			}
//		}
//
//		if len(timeRange) == 0 {
//			continue
//		}
//		minTime := timeRange[0]
//		maxTime := timeRange[0]
//		for _, t := range timeRange {
//			if t.Before(minTime) {
//				minTime = t
//			}
//			if t.After(maxTime) {
//				maxTime = t
//			}
//		}
//		group[i].TimeSpan = maxTime.Unix() - minTime.Unix()
//	}
//
//	return &vo.HomeIssueInfoGroupResp{
//		Total:       list.Total,
//		ActualTotal: list.ActualTotal,
//		TimeSpan:    timeSpan,
//		Group:       group,
//	}, nil
//}

// GetIssueIdsByOrgId 根据一个组织的一批任务 id 列表
//func GetIssueIdsByOrgId(orgId, userId int64, input *projectvo.GetIssueIdsByOrgIdReq) (*projectvo.GetIssueIdsByOrgIdResp, errs.SystemErrorInfo) {
//	result := &projectvo.GetIssueIdsByOrgIdResp{
//		List: make([]int64, 0),
//	}
//	ids, total, err := domain.GetIssueIdsByOrgId(orgId, input.Page, input.Size)
//	if err != nil {
//		log.Error(err)
//		return result, err
//	}
//	result.List = ids
//	result.Total = total
//
//	return result, nil
//}

// InsertIssueProRelation 根据任务 id，将其与项目的关联，新增一条关联数据到 issue_relation 中。
//func InsertIssueProRelation(orgId, userId int64, input *projectvo.InsertIssueProRelationReq) errs.SystemErrorInfo {
//	issueIds := input.IssueIds
//	// 检查该任务是否已存在关联关系，已存在则跳过，不存在，则插入关联关系
//	existRelations, err := domain.GetIssueRelationsByCond(issueIds, []int{consts.IssueRelationTypeBelongManyPro}, db.Cond{
//		consts.TcOrgId: orgId,
//	})
//	if err != nil {
//		log.Errorf("[InsertIssueProRelation] err: %v", err)
//		return err
//	}
//	existIssueIds := make([]int64, 0)
//	for _, item := range existRelations {
//		existIssueIds = append(existIssueIds, item.IssueId)
//	}
//	notExistIssueIds := int642.ArrayDiff(issueIds, existIssueIds)
//	if len(notExistIssueIds) < 1 {
//		// 没有要处理的任务
//		return nil
//	}
//	// 查询任务信息
//	cond1 := db.Cond{
//		consts.TcOrgId: orgId,
//		consts.TcId:    notExistIssueIds,
//	}
//	bos, _, err := domain.SelectList(cond1, nil, 1, 10_0000, "id desc", false)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	waitRelationBos := make([]*po.PpmPriIssueRelation, 0)
//	// 组装 issue relation
//	oids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, len(*bos))
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	for index, issueBo := range *bos {
//		oid := oids.Ids[index].Id
//		waitRelationBos = append(waitRelationBos, &po.PpmPriIssueRelation{
//			Id:           oid,
//			OrgId:        issueBo.OrgId,
//			ProjectId:    issueBo.ProjectId,
//			IssueId:      issueBo.Id,
//			RelationId:   issueBo.ProjectObjectTypeId, // 这里的值是 TableId，也就是任务栏。
//			RelationCode: consts.IssueRelationBelongManyProCode,
//			RelationType: consts.IssueRelationTypeBelongManyPro,
//			Creator:      issueBo.Creator,
//			CreateTime:   time.Time(issueBo.CreateTime),
//			Updator:      issueBo.Updator,
//			UpdateTime:   time.Time(issueBo.UpdateTime),
//			Version:      issueBo.Version,
//			IsDelete:     issueBo.IsDelete,
//		})
//	}
//	insertErr := mysql.BatchInsert(&po.PpmPriIssueRelation{}, slice.ToSlice(waitRelationBos))
//	if insertErr != nil {
//		log.Error(insertErr)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertErr)
//	}
//
//	return nil
//}

func GetTableStatus(orgId int64, tableId int64) ([]status.StatusInfoBo, errs.SystemErrorInfo) {
	tableStatus, err := domain.GetTableStatus(orgId, tableId)
	if err != nil {
		log.Errorf("[GetTableStatus]错误, orgId:%v, tableId:%v, err:%v", orgId, tableId, err)
		return nil, err
	}
	return tableStatus, nil
}
