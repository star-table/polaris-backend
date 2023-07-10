package tablefacade

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/star-table/common/core/threadlocal"

	middlewareMeta "github.com/star-table/go-common/pkg/middleware/meta"
	v1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/spf13/cast"
	goGrpc "google.golang.org/grpc"
)

const (
	metaUserId    = "x-md-userid"
	metaOrgId     = "x-md-orgid"
	metaPmTraceId = "pm-trace-id"
)

var grpcClient v1.RowsClient

func InitGrpcClient(discover registry.Discovery) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///go-table.grpc"),
		grpc.WithTimeout(15*time.Second),
		grpc.WithDiscovery(discover),
		grpc.WithMiddleware(
			recovery.Recovery(),
			validate.Validator(),
			middlewareMeta.Client(),
		),
		grpc.WithOptions(goGrpc.WithDefaultCallOptions(goGrpc.MaxCallRecvMsgSize(50000000))),
	)
	if err != nil {
		panic(err)
	}
	grpcClient = v1.NewRowsClient(conn)
}

// 改造 objectType、process 时需要用的接口调用

// GetTableColumns 查询表配置中的列信息。接口地址待确认
func GetTableColumns(req projectvo.GetTableColumnsReq) *projectvo.TablesColumnsResp {
	respVo := &projectvo.TablesColumnsResp{Data: &projectvo.TablesColumnsRespData{}}
	reqUrl := fmt.Sprintf("%s%s/read/tableSchemas", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetTablesOfApp 获取应用的 table list 等待 tablesvc 实现
func GetTablesOfApp(req projectvo.GetTablesOfAppReq) *projectvo.GetTablesOfAppRespVo {
	respVo := &projectvo.GetTablesOfAppRespVo{Data: &v1.ReadTablesByAppsReply{}}
	reqUrl := fmt.Sprintf("%s%s/read/tablesByApps", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func GetTablesByOrg(req projectvo.GetTablesByOrgReq) *projectvo.GetTablesByOrgRespVo {
	respVo := &projectvo.GetTablesByOrgRespVo{Data: &projectvo.TableData{}}
	reqUrl := fmt.Sprintf("%s%s/read/org/tables", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// CreateSummeryTable 创建汇总表
func CreateSummeryTable(req projectvo.CreateSummeryTableReq) *projectvo.CreateSummeryTableRespVo {
	respVo := &projectvo.CreateSummeryTableRespVo{Data: &v1.CreateSummeryTableReply{}}
	reqUrl := fmt.Sprintf("%s%s/read/createSummery", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// CreateTable 给应用创建 table（表）
func CreateTable(req projectvo.CreateTableReq) *projectvo.CreateTableRespVo {
	respVo := &projectvo.CreateTableRespVo{Data: &projectvo.CreateTableReply{}}
	reqUrl := fmt.Sprintf("%s%s/table/create", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// RenameTable
func RenameTable(req projectvo.RenameTableReq) *projectvo.RenameTableResp {
	respVo := &projectvo.RenameTableResp{Data: &projectvo.TableMetaData{}}
	reqUrl := fmt.Sprintf("%s%s/table/rename", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// DeleteTable
func DeleteTable(req projectvo.DeleteTableReq) *projectvo.DeleteTableResp {
	respVo := &projectvo.DeleteTableResp{Data: &projectvo.TableMetaData{}}
	reqUrl := fmt.Sprintf("%s%s/table/delete", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// SetAutoSchedule
func SetAutoSchedule(req projectvo.SetAutoScheduleReq) *projectvo.SetAutoScheduleResp {
	respVo := &projectvo.SetAutoScheduleResp{Data: &projectvo.TableAutoSchedule{}}
	reqUrl := fmt.Sprintf("%s%s/table/setAutoSchedule", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// ReadTables
func ReadTables(req projectvo.GetTablesReqVo) *projectvo.GetTablesDataResp {
	respVo := &projectvo.GetTablesDataResp{Data: &projectvo.TableData{}}
	reqUrl := fmt.Sprintf("%s%s/read/tables", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// ReadOneTable
func ReadOneTable(req projectvo.GetTableInfoReq) *projectvo.GetTableInfoResp {
	respVo := &projectvo.GetTableInfoResp{Data: &projectvo.ReadTableReply{}}
	reqUrl := fmt.Sprintf("%s%s/read/table", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetSummeryTableId 获取汇总表的tableId
func GetSummeryTableId(req projectvo.GetSummeryTableIdReqVo) *projectvo.GetSummeryTableIdRespVo {
	respVo := &projectvo.GetSummeryTableIdRespVo{Data: &v1.ReadSummeryTableIdReply{}}
	reqUrl := fmt.Sprintf("%s%s/read/summeryTableId", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// ReadTablesByApps 通过 appIds 获取对应的 table list
func ReadTablesByApps(req projectvo.ReadTablesByAppsReqVo) *projectvo.ReadTablesByAppsRespVo {
	respVo := &projectvo.ReadTablesByAppsRespVo{Data: &projectvo.ReadTablesByAppsData{}}
	reqUrl := fmt.Sprintf("%s%s/read/tablesByApps", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// InitOrgColumns 初始化团队字段
func InitOrgColumns(req orgvo.InitOrgColumnsReq) *orgvo.InitOrgColumnRespVo {
	respVo := &orgvo.InitOrgColumnRespVo{Data: &orgvo.InitOrgColumnsReply{}}
	reqUrl := fmt.Sprintf("%s%s/org/columns/init", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetOrgColumns  读取组织字段列表
func GetOrgColumns(req orgvo.GetOrgColumnsReq) *orgvo.GetOrgColumnsRespVo {
	respVo := &orgvo.GetOrgColumnsRespVo{Data: &orgvo.ReadOrgColumnsReply{}}
	reqUrl := fmt.Sprintf("%s%s/read/org/columns", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// CreateOrgColumn 创建组织字段
func CreateOrgColumn(req orgvo.CreateOrgColumnReq) *orgvo.CreateOrgColumnRespVo {
	respVo := &orgvo.CreateOrgColumnRespVo{Data: &v1.CreateOrgColumnReply{}}
	reqUrl := fmt.Sprintf("%s%s/org/column/create", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// DeleteOrgColumn 删除组织字段
func DeleteOrgColumn(req orgvo.DeleteOrgColumnReq) *orgvo.DeleteOrgColumnRespVo {
	respVo := &orgvo.DeleteOrgColumnRespVo{Data: &v1.DeleteOrgColumnReply{}}
	reqUrl := fmt.Sprintf("%s%s/org/column/delete", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	if respVo.Err.Code == 403 {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.LcCanNotDeleteByUse))
	}
	return respVo
}

// CreateColumn 创建普通表头
func CreateColumn(req projectvo.CreateColumnReqVo) *projectvo.CreateColumnRespVo {
	respVo := &projectvo.CreateColumnRespVo{Data: &projectvo.CreateColumnReply{}}
	reqUrl := fmt.Sprintf("%s%s/column/create", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// CopyColumn 拷贝列,包含数据拷贝
func CopyColumn(req projectvo.CopyColumnReqVo) *projectvo.CopyColumnRespVo {
	respVo := &projectvo.CopyColumnRespVo{Data: &v1.CopyColumnReply{}}
	reqUrl := fmt.Sprintf("%s%s/column/copy", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// UpdateColumn 更新普通表头
func UpdateColumn(req projectvo.UpdateColumnReqVo) *projectvo.UpdateColumnRespVo {
	respVo := &projectvo.UpdateColumnRespVo{Data: &projectvo.UpdateColumnReply{}}
	reqUrl := fmt.Sprintf("%s%s/column/update", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// UpdateColumnDescription 更新表头字段描述
func UpdateColumnDescription(req projectvo.UpdateColumnDescriptionReqVo) *projectvo.UpdateColumnDescriptionRespVo {
	respVo := &projectvo.UpdateColumnDescriptionRespVo{Data: &v1.UpdateColumnDescriptionReply{}}
	reqUrl := fmt.Sprintf("%s%s/column/description/update", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// DeleteColumn 删除表头
func DeleteColumn(req projectvo.DeleteColumnReqVo) *projectvo.DeleteColumnRespVo {
	respVo := &projectvo.DeleteColumnRespVo{Data: &v1.DeleteColumnReply{}}
	reqUrl := fmt.Sprintf("%s%s/column/delete", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// ReadTableSchemasByAppId  获取某些列的配置
func ReadTableSchemasByAppId(req projectvo.GetTableSchemasByAppIdReq) *projectvo.GetTableSchemasByAppIdRespVo {
	respVo := &projectvo.GetTableSchemasByAppIdRespVo{Data: &projectvo.TablesColumnsRespData{}}
	reqUrl := fmt.Sprintf("%s%s/read/tableSchemasByAppId", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req.Input, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// ReadTableSchemasByOrgId  获取整个组织某些列
func ReadTableSchemasByOrgId(req projectvo.GetTableSchemasByOrgIdReq) *projectvo.GetTableSchemasByOrgIdRespVo {
	respVo := &projectvo.GetTableSchemasByOrgIdRespVo{Data: &projectvo.TablesColumnsRespData{}}
	reqUrl := fmt.Sprintf("%s%s/read/orgTableSchemas", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(req.OrgId, req.UserId, consts.HttpMethodPost, reqUrl, nil, nil, req, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func RecycleAttachment(orgId, userId int64, req *v1.RecycleAttachmentRequest) *projectvo.RecycleAttachmentResp {
	respVo := &projectvo.RecycleAttachmentResp{Data: &v1.RecycleAttachmentReply{}}
	reqUrl := fmt.Sprintf("%s%s/rows/attachment/recycle", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(orgId, userId, consts.HttpMethodPost, reqUrl, nil, nil, req, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func RecoverAttachment(orgId, userId int64, req *v1.RecoverAttachmentRequest) *projectvo.RecoverAttachmentResp {
	respVo := &projectvo.RecoverAttachmentResp{Data: &v1.RecoverAttachmentReply{}}
	reqUrl := fmt.Sprintf("%s%s/rows/attachment/recover", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(orgId, userId, consts.HttpMethodPost, reqUrl, nil, nil, req, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func List(orgId, userId int64, req *v1.ListRequest) *projectvo.ListRowsReply {
	respVo := &projectvo.ListRowsReply{Data: &v1.ListReply{}}
	traceId := threadlocal.GetTraceId()
	ctx := metadata.NewClientContext(context.Background(), metadata.Metadata{metaPmTraceId: traceId, metaUserId: cast.ToString(userId), metaOrgId: cast.ToString(orgId)})
	reply, err := grpcClient.List(ctx, req)
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.TableDomainError, err))
		return respVo
	}
	respVo.Data = reply

	return respVo
}

func ListRaw(orgId, userId int64, req *v1.ListRawRequest) *projectvo.ListRawRowsReply {
	respVo := &projectvo.ListRawRowsReply{Data: &v1.ListRawReply{}}
	traceId := threadlocal.GetTraceId()
	ctx := metadata.NewClientContext(context.Background(), metadata.Metadata{metaPmTraceId: traceId, metaUserId: cast.ToString(userId), metaOrgId: cast.ToString(orgId)})
	reply, err := grpcClient.ListRaw(ctx, req)
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.TableDomainError, err))
		return respVo
	}
	respVo.Data = reply

	return respVo
}

func DeleteRows(orgId, userId int64, req *v1.DeleteRequest) *projectvo.DeleteReply {
	respVo := &projectvo.DeleteReply{Data: &v1.DeleteReply{}}
	reqUrl := fmt.Sprintf("%s%s/rows/delete", config.GetPreUrl(consts.ServiceTable), ApiV1Prefix)
	err := facade.RequestWithCommonHeader(orgId, userId, consts.HttpMethodPost, reqUrl, nil, nil, req, respVo.Data)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

func ListRawExpand(orgId, userId int64, req *v1.ListRawRequest) (*formvo.LessIssueRawListResp, errs.SystemErrorInfo) {
	respVo := ListRaw(orgId, userId, req)
	if respVo.Failure() {
		return nil, respVo.Error()
	}

	result := &formvo.LessIssueRawListResp{Data: []map[string]interface{}{}}
	err2 := json.Unmarshal(respVo.Data.Data, &result.Data)
	if err2 != nil {
		log.Errorf("[ListRawExpand] Unmarshal, err: %v", err2)
		return nil, errs.JSONConvertError
	}

	return result, nil
}
