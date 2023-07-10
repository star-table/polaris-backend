package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) GetOneTableColumns(req projectvo.GetTableColumnReq) projectvo.TableColumnsResp {
	resp, err := service.GetOneTableColumns(req)
	return projectvo.TableColumnsResp{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (PostGreeter) GetTablesColumns(req projectvo.GetTablesColumnsReq) projectvo.TablesColumnsResp {
	resp, err := service.GetTablesColumns(req)
	return projectvo.TablesColumnsResp{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (PostGreeter) CreateTable(req projectvo.CreateTableReq) projectvo.CreateTableRespVo {
	res, err := service.CreateTable(req)
	return projectvo.CreateTableRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) RenameTable(req projectvo.RenameTableReq) projectvo.RenameTableResp {
	res, err := service.RenameTable(req)
	return projectvo.RenameTableResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) DeleteTable(req projectvo.DeleteTableReq) projectvo.DeleteTableResp {
	res, err := service.DeleteTable(req)
	return projectvo.DeleteTableResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) SetAutoSchedule(req projectvo.SetAutoScheduleReq) projectvo.SetAutoScheduleResp {
	res, err := service.SetAutoSchedule(req)
	return projectvo.SetAutoScheduleResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetTable(req projectvo.GetTableInfoReq) projectvo.GetTableInfoResp {
	res, err := service.GetTable(req)
	return projectvo.GetTableInfoResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetTables(req projectvo.GetTablesReqVo) projectvo.GetTablesDataResp {
	res, err := service.GetTables(req.OrgId, req.UserId, req.Input.AppId)
	return projectvo.GetTablesDataResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetTablesByApps(req projectvo.ReadTablesByAppsReqVo) projectvo.ReadTablesByAppsRespVo {
	res, err := service.GetTablesByApps(req.OrgId, req.UserId, req.Input.AppIds)
	return projectvo.ReadTablesByAppsRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetTablesByOrg(req projectvo.GetTablesByOrgReq) projectvo.GetTablesByOrgRespVo {
	res, err := service.GetTablesByOrg(req.OrgId, req.UserId)
	return projectvo.GetTablesByOrgRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
