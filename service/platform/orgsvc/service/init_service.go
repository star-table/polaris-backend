package service

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/appvo"
	"github.com/star-table/polaris-backend/facade/appfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
)

// NewbieGuideInit 初始化lark通用数据
func NewbieGuideInit(orgId, userId int64, sourceChannel string) errs.SystemErrorInfo {
	//departmentInfo, err := domain.GetTopDepartmentInfoList(orgId)
	//if err != nil {
	//	log.Error("获取部门信息错误 " + strs.ObjectToString(err))
	//	return err
	//}
	//var departmentId int64
	//for _, v := range departmentInfo {
	//	departmentId = v.Id
	//	break
	//}

	////用户初始化
	//zhangsanId, lisiId, err := domain.LarkUserInit(orgId, sourceChannel, sourcePlatform, departmentId)
	//if err != nil {
	//	return err
	//}
	//log.Info("用户初始化成功")

	////项目初始化
	//preCode := "XSZN"
	//remark := "在这个项目中，我们会向您依次介绍极星的基础架构和产品功能，让你更快上手这款简洁易用的协同办公软件。"
	//start := types.NowTime()
	//endTime, _ := time.Parse(consts.AppTimeFormat, "2099-12-12 12:00:00")
	//end := types.Time(endTime)
	//// 创建新手项目，项目类型设为新项目类型
	//projectTypeId := int64(consts.ProjectTypeCommon2022V47)
	//
	//trueFlag := true
	//projectInfo := projectfacade.CreateProject(projectvo.CreateProjectReqVo{Input: vo.CreateProjectReq{
	//	Name:         "新手指南",
	//	PreCode:      &preCode,
	//	PublicStatus: consts.PublicProject,
	//	Remark:       &remark,
	//	Owner:        userId,
	//	//MemberIds:     []int64{userId},
	//	MemberForDepartmentID: []int64{0},
	//	ResourcePath:          "https://polaris-hd2.oss-cn-shanghai.aliyuncs.com/project/undraw_Projectpicture_update_jjgk.png",
	//	ResourceType:          consts.OssResource,
	//	PlanStartTime:         &start,
	//	PlanEndTime:           &end,
	//	ProjectTypeID:         &projectTypeId,
	//	IsFirst:               &trueFlag,
	//},
	//	UserId: userId,
	//	OrgId:  orgId,
	//})
	//if projectInfo.Failure() {
	//	log.Errorf("[NewbieGuideInit] err: %v", projectInfo.Error())
	//	return errs.BuildSystemErrorInfo(errs.BaseDomainError, projectInfo.Error())
	//}
	//log.Infof("项目初始化成功 orgId: %d", orgId)
	//proAppId, oriErr := strconv.ParseInt(projectInfo.Project.AppID, 10, 64)
	//if oriErr != nil {
	//	log.Errorf("[NewbieGuideInit] orgId: %d, proAppId: %v, parse proAppId err: %v, ", orgId, projectInfo.Project.AppID, oriErr)
	//	return errs.BuildSystemErrorInfo(errs.TypeConvertError, oriErr)
	//}
	//// 1:迭代；2表示需求、任务、缺陷这一类
	//// objectType := 2
	//// 初始化项目是一个通用项目，设置它的任务栏
	//tableId, err := InitTaskBarList(orgId, userId, projectInfo.Project.ID, proAppId, []string{
	//	"极星是什么？",
	//	"快速入门",
	//	"进阶",
	//	"高阶",
	//	"出神入化",
	//})
	//if err != nil {
	//	log.Errorf("[NewbieGuideInit] orgId: %d, err: %v", orgId, err)
	//	return err
	//}
	//
	////任务初始化
	//issueInitErr := projectfacade.IssueLarkInit(projectvo.NewbieGuideIssuesInitReqVo{
	//	//ZhangsanId: int64(0),
	//	//LisiId:     int64(0),
	//	OrgId:     orgId,
	//	ProjectId: projectInfo.Project.ID,
	//	UserId:    userId,
	//	TableId:   tableId,
	//})
	//if issueInitErr.Failure() {
	//	log.Errorf("[NewbieGuideInit] orgId: %d, err: %v", orgId, issueInitErr.Error())
	//	return issueInitErr.Error()
	//}

	newbieGuideTemplateId := domain.GetNewbieGuideTemplateId(sourceChannel)
	if newbieGuideTemplateId != 0 {
		resp := appfacade.ApplyTemplate(appvo.ApplyTemplateReq{
			OrgId:         orgId,
			UserId:        userId,
			TemplateId:    newbieGuideTemplateId,
			IsNewbieGuide: true,
		})
		if resp.Failure() {
			log.Errorf("[NewbieGuideInit] ApplyTemplate orgId: %d, userId: %d, templateId: %d, err: %v", orgId, userId, newbieGuideTemplateId, resp.Error())
			return nil
		}
	}

	log.Infof("任务初始化成功 orgId: %d", orgId)

	return nil
}

//func initProjectObjectType(orgId, userId int64, projectId int64, objectType int, name string) errs.SystemErrorInfo {
//	projectType1 := projectfacade.CreateProjectObjectType(projectvo.CreateProjectObjectTypeReqVo{
//		Input: vo.CreateProjectObjectTypeReq{
//			ProjectID:  projectId,
//			Name:       name,
//			ObjectType: objectType,
//			BeforeID:   0,
//		},
//		OrgId:  orgId,
//		UserId: userId,
//	})
//	if projectType1.Failure() {
//		return errs.BuildSystemErrorInfo(errs.BaseDomainError, projectType1.Error())
//	}
//	log.Infof("项目对象类型-%s初始化成功", name)
//	return nil
//}

// initProjectTableList 在一个项目下，初始化一批 table
//func initProjectTableList(orgId, userId int64, proAppId int64, tableNames []string) errs.SystemErrorInfo {
//	for _, tableName := range tableNames {
//		if err := initOneProjectTable(orgId, userId, proAppId, tableName); err != nil {
//			log.Errorf("[initProjectTableList] err: %v", err)
//			return err
//		}
//	}
//
//	return nil
//}

// initOneProjectTable 初始化一个项目的 table
//func initOneProjectTable(orgId, userId, proAppId int64, tableName string) errs.SystemErrorInfo {
//	createResp := tablefacade.CreateTable(projectvo.CreateTableReq{
//		OrgId:  orgId,
//		UserId: userId,
//		Input: &tableV1.CreateTableRequest{
//			AppId:            proAppId,
//			Name:             tableName,
//			IsNeedStoreTable: false,
//			IsNeedColumn:     true,
//		},
//	})
//	if createResp.Failure() {
//		log.Errorf("[initOneProjectTable] err: %v", createResp.Error())
//		return createResp.Error()
//	}
//
//	return nil
//}

// InitTaskBarList 编辑项目下，表的任务栏列，增加一些特殊任务栏
//func InitTaskBarList(orgId, userId int64, projectId, proAppId int64, barNames []string) (int64, errs.SystemErrorInfo) {
//	// 查询应用下的 table
//	tableResp := tablefacade.ReadTables(projectvo.GetTablesReqVo{
//		OrgId:  orgId,
//		UserId: userId,
//		Input: &tableV1.ReadTablesRequest{
//			AppId: proAppId,
//		},
//	})
//	if tableResp.Failure() {
//		log.Errorf("[InitTaskBarList] err: %v", tableResp.Error())
//		return 0, tableResp.Error()
//	}
//	if len(tableResp.Data.Tables) == 0 {
//		return 0, errs.TableNotExist
//	}
//	curTable := tableResp.Data.Tables[0]
//
//	// 查询 table 下的“任务栏”列
//	columnArrResp := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
//		OrgId:  orgId,
//		UserId: userId,
//		Input: &tableV1.ReadTableSchemasRequest{
//			TableIds:          []int64{curTable.TableId},
//			ColumnIds:         []string{consts.BasicFieldProjectObjectTypeId}, // 任务栏
//			IsNeedDescription: false,
//		},
//	})
//	if columnArrResp.Failure() {
//		log.Errorf("[InitTaskBarList] err: %v", columnArrResp.Error())
//		return 0, columnArrResp.Error()
//	}
//	// 任务栏列
//	taskBarColumn := &projectvo.TableColumnData{}
//	for _, tableSchema := range columnArrResp.Data.Tables {
//		if tableSchema.TableId == curTable.TableId {
//			for _, tmpColumn := range tableSchema.Columns {
//				if tmpColumn.Name == consts.BasicFieldProjectObjectTypeId {
//					taskBarColumn = tmpColumn
//				}
//			}
//			break
//		}
//	}
//	// 更新列，配置新的选项值
//	assemblyUpdateTaskBarColumn(taskBarColumn, barNames)
//	updateResp := tablefacade.UpdateColumn(projectvo.UpdateColumnReqVo{
//		OrgId:  orgId,
//		UserId: userId,
//		Input: &projectvo.UpdateColumnReqVoInput{
//			ProjectId: projectId,
//			AppId:     proAppId,
//			TableId:   curTable.TableId,
//			Column:    taskBarColumn,
//		},
//	})
//	if updateResp.Failure() {
//		log.Errorf("[InitTaskBarList] err: %v", updateResp.Error())
//		return 0, updateResp.Error()
//	}
//
//	return curTable.TableId, nil
//}

// assemblyUpdateTaskBarColumn 组装新的任务栏结构
//func assemblyUpdateTaskBarColumn(oldColumn *projectvo.TableColumnData, barNames []string) {
//	optionList := make([]projectvo.ColumnSelectOption, 0)
//	for i, barName := range barNames {
//		optionList = append(optionList, projectvo.ColumnSelectOption{
//			Id:    consts.DefaultProjectTaskBarIds[i],
//			Value: barName,
//		})
//	}
//
//	oldColumn.Field.Props["select"] = projectvo.FormConfigColumnFieldMultiselectPropsMultiselectForInterfaceId{
//		Options: optionList,
//	}
//}
