package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
)

// 任务模块：任务状态筛选条件转换，前端其实筛选的是issueStatusType
func ConvertIssueStatusFilterReqForAll(lessReq vo.LessCondsData) vo.LessCondsData {
	if lessReq.Column == consts.BasicFieldIssueStatus {
		lessReq = convertIssueStatusCondForAll(lessReq)
	}
	convertIssueStatusCondsForAll(&lessReq.Conds)
	return lessReq
}

func convertIssueStatusCondsForAll(conds *[]*vo.LessCondsData) {
	if conds == nil || len(*conds) == 0 {
		return
	}
	for i, cond := range *conds {
		if cond.Column == consts.BasicFieldIssueStatus {
			cur := convertIssueStatusCondForAll(*cond)
			(*conds)[i] = &cur
		}
		convertIssueStatusCondsForAll(&cond.Conds)
	}
}

// 已完成 + 待确认 + 进行中 + 未开始 = 100%
// 已完成 + 未完成 = 100%
// 未完成 = 待确认 + 进行中 + 未开始
var waitConfirmType int64 = -1 // 待确认
var completeType int64 = -2    // 已完成
var processType int64 = -3     // 进行中
var notCompleteType int64 = -4 // 未完成 (已经不支持，只兼容老数据)
var notStartType int64 = -5    // 未开始
var allStatusType = []int64{waitConfirmType, completeType, processType, notStartType}

func convertIssueStatusCondForAll(cond vo.LessCondsData) vo.LessCondsData {
	//not_in转化成in, 因为值是固定已知的（-1，-2, -3, -4, -5）
	newValues := make([]int64, 0)
	intValues := make([]int64, 0)
	_ = copyer.Copy(cond.Values, &intValues)
	// 兼容下"未完成"
	if ok, _ := slice.Contain(intValues, notCompleteType); ok {
		values := make([]int64, 0)
		for _, value := range intValues {
			if value != notCompleteType {
				values = append(values, value)
			}
		}
		values = append(values, waitConfirmType, processType, notStartType)
		intValues = values
	}
	if cond.Type == "not_in" {
		for _, value := range allStatusType {
			if ok, _ := slice.Contain(intValues, value); !ok {
				newValues = append(newValues, value)
			}
		}
	} else {
		newValues = intValues
	}
	if len(newValues) > 0 {
		res := vo.LessCondsData{
			Type:  "or",
			Conds: nil,
		}
		for _, value := range newValues {
			switch value {
			case notStartType:
				res.Conds = append(res.Conds, &vo.LessCondsData{
					Type:   "equal",
					Value:  consts.StatusTypeNotStart,
					Column: consts.BasicFieldIssueStatusType,
					Conds:  nil,
				})
			//case notCompleteType:
			//	res.Conds = append(res.Conds, &vo.LessCondsData{
			//		Type: "or",
			//		Conds: []*vo.LessCondsData{
			//			&vo.LessCondsData{
			//				Type:   "equal",
			//				Value:  consts.StatusTypeNotStart,
			//				Column: consts.BasicFieldIssueStatusType,
			//				Conds:  nil,
			//			},
			//			&vo.LessCondsData{
			//				Type:   "equal",
			//				Value:  consts.StatusTypeRunning,
			//				Column: consts.BasicFieldIssueStatusType,
			//				Conds:  nil,
			//			},
			//			&vo.LessCondsData{
			//				Type: "and",
			//				Conds: []*vo.LessCondsData{
			//					&vo.LessCondsData{
			//						Type:   "equal",
			//						Value:  consts.StatusTypeComplete,
			//						Column: consts.BasicFieldIssueStatusType,
			//						Conds:  nil,
			//					},
			//					&vo.LessCondsData{
			//						Type:   "in",
			//						Values: []int{consts.AuditStatusNotView, consts.AuditStatusView, consts.AuditStatusReject},
			//						Column: consts.BasicFieldAuditStatus,
			//						Conds:  nil,
			//					},
			//				},
			//			},
			//		},
			//	})
			case processType:
				res.Conds = append(res.Conds, &vo.LessCondsData{
					Type:   "equal",
					Value:  consts.StatusTypeRunning,
					Column: consts.BasicFieldIssueStatusType,
					Conds:  nil,
				})
			case completeType:
				res.Conds = append(res.Conds, &vo.LessCondsData{
					Type: "and",
					Conds: []*vo.LessCondsData{
						&vo.LessCondsData{
							Type:   "equal",
							Value:  consts.StatusTypeComplete,
							Column: consts.BasicFieldIssueStatusType,
							Conds:  nil,
						},
						&vo.LessCondsData{
							Type: "or",
							Conds: []*vo.LessCondsData{
								&vo.LessCondsData{
									Type:   "in",
									Values: []int{consts.AuditStatusNoNeed, consts.AuditStatusPass}, // -1是非通用项目确认状态的初始值
									Column: consts.BasicFieldAuditStatus,
									Conds:  nil,
								},
								&vo.LessCondsData{
									Type:   "is_null",
									Column: consts.BasicFieldAuditorIds,
									Conds:  nil,
								},
							},
						},
					},
				})
			case waitConfirmType:
				res.Conds = append(res.Conds, &vo.LessCondsData{
					Type: "and",
					Conds: []*vo.LessCondsData{
						&vo.LessCondsData{
							Type:   "equal",
							Value:  consts.StatusTypeComplete,
							Column: consts.BasicFieldIssueStatusType,
							Conds:  nil,
						},
						&vo.LessCondsData{
							Type:   "in",
							Values: []int{consts.AuditStatusNotView, consts.AuditStatusView, consts.AuditStatusReject},
							Column: consts.BasicFieldAuditStatus,
							Conds:  nil,
						},
						&vo.LessCondsData{
							Type:   "not_null",
							Column: consts.BasicFieldAuditorIds,
							Conds:  nil,
						},
					},
				})
			}
		}
		return res
	} else {
		return vo.LessCondsData{
			Type:   "raw_sql",
			Column: "1 <> 1",
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 单个老通用项目：任务状态筛选条件转换，主要处理待确认状态-1
func ConvertIssueStatusFilterReqForNormalProject(lessReq *vo.LessCondsData) *vo.LessCondsData {
	if lessReq.Column == consts.BasicFieldIssueStatus {
		lessReq = convertIssueStatusCondForNormalProject(lessReq)
	}
	convertIssueStatusCondsForNormalProject(&lessReq.Conds)
	return lessReq
}

func convertIssueStatusCondsForNormalProject(conds *[]*vo.LessCondsData) {
	if conds == nil || len(*conds) == 0 {
		return
	}
	for i, cond := range *conds {
		if cond.Column == consts.BasicFieldIssueStatus {
			cur := convertIssueStatusCondForNormalProject(cond)
			(*conds)[i] = cur
		}
		convertIssueStatusCondsForNormalProject(&cond.Conds)
	}
}

var waitConfirm int64 = -1
var notStart int64 = 7
var process int64 = 16
var finish int64 = 26
var allStatus = []int64{waitConfirm, notStart, process, finish}

func convertIssueStatusCondForNormalProject(cond *vo.LessCondsData) *vo.LessCondsData {
	//not_in转化成in,因为值是固定已知的（-1, 7, 16, 26）
	newValues := make([]int64, 0)
	intValues := make([]int64, 0)
	newValuesInterface := make([]interface{}, 0)
	_ = copyer.Copy(cond.Values, &intValues)
	if cond.Type == "not_in" {
		for _, value := range allStatus {
			if ok, _ := slice.Contain(intValues, value); !ok {
				newValues = append(newValues, value)
			}
		}
	} else {
		newValues = intValues
	}

	// 查询条件优化
	if len(newValues) == 0 || len(newValues) == len(allStatus) {
		copyer.Copy(newValues, &newValuesInterface)
		// 等于0或者等于4直接放原来的就行，4的话表示所有的都筛出来（7，16, 26）-1自然也包括其中了
		return &vo.LessCondsData{
			Type:   "in",
			Values: newValuesInterface,
			Column: consts.BasicFieldIssueStatus,
			Conds:  nil,
		}
	}

	res := &vo.LessCondsData{
		Type:  "or",
		Conds: nil,
	}
	for _, newValue := range newValues {
		switch newValue {
		case waitConfirm:
			res.Conds = append(res.Conds, &vo.LessCondsData{
				Type:   "and",
				Column: "",
				Conds: []*vo.LessCondsData{
					&vo.LessCondsData{
						Type:   "equal",
						Value:  finish,
						Column: consts.BasicFieldIssueStatus,
						Conds:  nil,
					},
					&vo.LessCondsData{
						Type:   "un_equal",
						Value:  consts.AuditStatusPass,
						Column: consts.BasicFieldAuditStatus,
						Conds:  nil,
					},
				},
			})
		case finish:
			res.Conds = append(res.Conds, &vo.LessCondsData{
				Type:   "and",
				Column: "",
				Conds: []*vo.LessCondsData{
					&vo.LessCondsData{
						Type:   "equal",
						Value:  finish,
						Column: consts.BasicFieldIssueStatus,
						Conds:  nil,
					},
					&vo.LessCondsData{
						Type:   "equal",
						Value:  consts.AuditStatusPass,
						Column: consts.BasicFieldAuditStatus,
						Conds:  nil,
					},
				},
			})
		case notStart:
			res.Conds = append(res.Conds, &vo.LessCondsData{
				Type:   "equal",
				Value:  newValue,
				Column: consts.BasicFieldIssueStatus,
				Conds:  nil,
			})
		case process:
			res.Conds = append(res.Conds, &vo.LessCondsData{
				Type:   "equal",
				Value:  newValue,
				Column: consts.BasicFieldIssueStatus,
				Conds:  nil,
			})
		}
	}
	return res
}

// ConvertConditionOfHomeIssues 任务状态改造后，将 homeIssues 接口的入参转换成无码支持的条件筛选 todo
// 根据咨询，openAPI 目前只用到了 3 个查询条件，因此这里暂且只对需要用到的条件进行转换
// 3 个条件：orderType,projectId,lastUpdateTime
// 转换后的结果存入 `input` 的 `LessConds` 中
//func ConvertConditionOfHomeIssues(input *vo.HomeIssueInfoReq) {
//	condArr := make([]*vo.LessCondsData, 0)
//	// 排序规则转换
//	if input.OrderType != nil {
//		input.LessOrder = ConvertConditionOfHomeIssuesForOrder(*input.OrderType)
//	}
//	if input.SearchCond != nil {
//		condArr = append(condArr, &vo.LessCondsData{
//			Type:   "like",
//			Value:  *input.SearchCond,
//			Column: consts.BasicFieldTitle,
//		})
//	}
//	if input.Status != nil {
//		if len(input.StatusList) < 1 {
//			input.StatusList = []int{*input.Status}
//		} else {
//			input.StatusList = append(input.StatusList, *input.Status)
//		}
//	}
//	if len(input.StatusList) > 0 {
//		// status interface type
//		statusIfArr := make([]interface{}, 0)
//		for _, status := range input.StatusList {
//			statusIfArr = append(statusIfArr, status)
//		}
//		condArr = append(condArr, &vo.LessCondsData{
//			Type:   "in",
//			Values: statusIfArr,
//			Column: consts.BasicFieldIssueStatusType,
//		})
//	}
//	if input.ProjectID != nil {
//		condArr = append(condArr, &vo.LessCondsData{
//			Type:   "equal",
//			Value:  *input.ProjectID,
//			Column: consts.BasicFieldProjectId,
//		})
//	}
//	if input.LastUpdateTime != nil {
//		condArr = append(condArr, &vo.LessCondsData{
//			Type:   "gte",
//			Value:  input.LastUpdateTime,
//			Column: consts.ProBasicFieldUpdateTime,
//		})
//	}
//	if input.IssueIds != nil {
//		condArr = append(condArr, &vo.LessCondsData{
//			Type:   "in",
//			Values: slice2.Int64Slice2IfSlice(input.IssueIds),
//			Column: consts.BasicFieldIssueId,
//		})
//	}
//	if input.OwnerChangeTimeStart != nil && input.OwnerChangeTimeEnd != nil {
//		condArr = append(condArr, &vo.LessCondsData{
//			Type:   "between",
//			Values: []interface{}{*input.OwnerChangeTimeStart, *input.OwnerChangeTimeEnd},
//			Column: "ownerChangeTime",
//		})
//	} else if input.OwnerChangeTimeStart != nil && input.OwnerChangeTimeEnd == nil {
//		condArr = append(condArr, &vo.LessCondsData{
//			Type:   "gte",
//			Value:  *input.OwnerChangeTimeStart,
//			Column: "ownerChangeTime",
//		})
//	} else if input.OwnerChangeTimeStart == nil && input.OwnerChangeTimeEnd != nil {
//		condArr = append(condArr, &vo.LessCondsData{
//			Type:   "lte",
//			Value:  *input.OwnerChangeTimeEnd,
//			Column: "ownerChangeTime",
//		})
//	}
//
//	resCondObj := &vo.LessCondsData{
//		Type:  "and",
//		Conds: condArr,
//	}
//	input.LessConds = resCondObj
//}

// ConvertConditionOfHomeIssuesForOrder
// 排序类型，1：项目分组，2：优先级分组，3：创建日期降序，4：最后更新日期降序, 5: 按开始时间最早, 6：按开始时间最晚,
// 8：按截止时间最近，9：按创建时间最早, 10: sort排序（正序）11：sort排序（倒序）12:截止时间（正序）13：优先级正序14：优先级倒序
// 15：负责人正序16：负责人倒序17：编号正序18：编号倒序19：标题正序20：标题倒序21：状态正序（必须传项目id，敏捷必须指定任务栏）
// 22：状态倒序（必须传项目id，敏捷必须指定任务栏）23:完成时间倒序24:按照传入id排序25:按照父任务正序26：按照父任务倒序
//func ConvertConditionOfHomeIssuesForOrder(orderFlag int) []*vo.LessOrder {
//	orderArr := make([]*vo.LessOrder, 0)
//	switch orderFlag {
//	case 1:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: consts.BasicFieldProjectId,
//		})
//	case 2:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: consts.BasicFieldPriority,
//		})
//	case 3:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: "createTime",
//		})
//	case 4:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: "updateTime",
//		})
//	case 5:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    true,
//			Column: consts.BasicFieldPlanStartTime,
//		})
//	case 6:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: consts.BasicFieldPlanStartTime,
//		})
//	case 8:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: consts.BasicFieldPlanEndTime,
//		})
//	case 9:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    true,
//			Column: "createTime",
//		})
//	case 10:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    true,
//			Column: "sort",
//		})
//	case 11:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: "sort",
//		})
//	case 12:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    true,
//			Column: consts.BasicFieldPlanEndTime,
//		})
//	case 13:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    true,
//			Column: consts.BasicFieldPriority,
//		})
//	case 14:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: consts.BasicFieldPriority,
//		})
//	case 15:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    true,
//			Column: consts.BasicFieldOwnerId,
//		})
//	case 16:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: consts.BasicFieldOwnerId,
//		})
//	case 17:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    true,
//			Column: consts.BasicFieldCode,
//		})
//	case 18:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: consts.BasicFieldCode,
//		})
//	case 19:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    true,
//			Column: consts.BasicFieldTitle,
//		})
//	case 20:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: consts.BasicFieldTitle,
//		})
//	case 21:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    true,
//			Column: consts.BasicFieldIssueStatus,
//		})
//	case 22:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    false,
//			Column: consts.BasicFieldIssueStatus,
//		})
//	case 24:
//		orderArr = append(orderArr, &vo.LessOrder{
//			Asc:    true,
//			Column: consts.BasicFieldIssueId,
//		})
//	}
//
//	return orderArr
//}
