package compatsvc

import "github.com/star-table/polaris-backend/common/model/vo"

type GetReportReqVo struct {
	Uid   int64 `json:"uid"`
	OrgId int64 `json:"orgId"`
	Page  int64 `json:"page"`
	Size  int64 `json:"size"`
}

type GetReportListRespVo struct {
	vo.Err
	GetReportRespVo GetReportRespVo `json:"data"`
}

type GetReportRespVo struct {
	Total       int                    `json:"total"`
	ReportArray []GetReportArrayRespVo `json:"reportArray"`
}

type GetReportArrayRespVo struct {
	ReportId    int64  `json:"reportId"`
	PackageName string `json:"packageName"`
	DeviceNum   int64  `json:"deviceNum"`
	FinishNum   int64  `json:"finishNum"`
	PassNum     int64  `json:"passNum"`
	Status      int64  `json:"status"`
	Creator     int64  `json:"creator"`
	CreatorName string `json:"creatorName"`
	CreateTime  string `json:"createTime"`
	updateTime  string `json:"updateTime"`
}

//删除
type DeleteReportReqVo struct {
	Uid      int64 `json:"uid"`
	OrgId    int64 `json:"orgId"`
	ReportId int64 `json:"reportId"`
}

type DeleteReportRespVo struct {
	vo.Err
}

//获取指定浏览报告的应用信息

type GetReportApkInfoReqVo struct {
	Uid      int64 `json:"uid"`
	OrgId    int64 `json:"orgId"`
	ReportId int64 `json:"reportId"`
}

type GetReportApkInfoRespVo struct {
	vo.Err
	GetReportApkInfoVo GetReportApkInfoVo `json:"data"`
}

type GetReportApkInfoVo struct {
	ReportId       int64  `json:"reportId"`
	Status         int64  `json:"status"`
	DeviceNum      int64  `json:"deviceNum"`
	PackageName    string `json:"packageName"`
	PackageVersion string `json:"packageVersion"`
	Size           string `json:"size"`
	CreateTime     string `json:"createTime"`
}

//获取详情报告
type ReportDetailOverViewReqVo struct {
	Uid      int64 `json:"uid"`
	OrgId    int64 `json:"orgId"`
	ReportId int64 `json:"reportId"`
}

type ReportDetailOverViewRespListVo struct {
	vo.Err
	ReportDetailOverViewRespVo ReportDetailOverViewRespVo `json:"data"`
}

type ReportDetailOverViewRespVo struct {
	Total                                  int                                    `json:"total"`
	ReportDetailOverViewArrayRespVo        []ReportDetailOverViewArrayRespVo      `json:"detailArray"`
	ReportDetailOverViewTestResultRespVo   ReportDetailOverViewTestResultRespVo   `json:"testResult"`
	ReportDetailOverViewDistributionRespVo ReportDetailOverViewDistributionRespVo `json:"distribution"`
}

type ReportDetailOverViewArrayRespVo struct {
	DetailId        int64  `json:"detailId"`
	ReportId        int64  `json:"reportId"`
	DeviceId        int64  `json:"deviceId"`
	DeviceName      string `json:"deviceName"`
	Manufacturer    string `json:"manufacturer"`
	Model           string `json:"model"`
	DeviceVersion   string `json:"deviceVersion"`
	ScreenSize      string `json:"screenSize"`
	Status          int64  `json:"status"`
	IsPass          int64  `json:"isPass"`
	InstallResult   int64  `json:"installResult"`
	LaunchResult    int64  `json:"launchResult"`
	MonkeyResult    int64  `json:"monkeyResult"`
	UninstallResult int64  `json:"uninstallResult"`
	ErrorMsg        string `json:"errorMsg"`
	ErrorStatus     int64  `json:"errorStatus"`
	CrashNum        int64  `json:"crashNum"`
	AnrNum          int64  `json:"anrNum"`
	ExceptionNum    int64  `json:"exceptionNum"`
}

type ReportDetailOverViewTestResultRespVo struct {
	FinishNum        int64 `json:"finishNum"`
	TotalNum         int64 `json:"totalNum"`
	PassNum          int64 `json:"passNum"`
	InstallFailedNum int64 `json:"installFailedNum"`
	LaunchFailedNum  int64 `json:"launchFailedNum"`
	MonkeyFailedNum  int64 `json:"monkeyFailedNum"`
	UninstallFailNum int64 `json:"uninstallFailNum"`
}

type ReportDetailOverViewDistributionRespVo struct {
	Manufacturer  map[string]int64 `json:"manufacturer"`
	DeviceVersion map[string]int64 `json:"deviceVersion"`
	ScreenSize    map[string]int64 `json:"screenSize"`
}

//获取指定类型

type GetReportDetailSingleReqVo struct {
	Uid      int64 `json:"uid"`
	OrgId    int64 `json:"orgId"`
	ReportId int64 `json:"reportId"`
	DetailId int64 `json:"detailId"`
}

type GetReportDetailSingleRespVo struct {
	vo.Err
	GetReportDetailSingleVo GetReportDetailSingleVo `json:"data"`
}

type GetReportDetailSingleVo struct {
	DeviceId        int64         `json:"deviceId"`
	DetailId        int64         `json:"detailId"`
	Platform        string        `json:"platform"`
	Model           string        `json:"model"`
	DeviceVerion    string        `json:"deviceVerion"`
	ScreenSize      string        `json:"screenSize"`
	Status          int64         `json:"status"`
	IsPass          int64         `json:"isPass"`
	InstallResult   int64         `json:"installResult"`
	InstallTime     int64         `json:"installTime"`
	LaunchResult    int64         `json:"launchResult"`
	LaunchTime      int64         `json:"launchTime"`
	MonkeyResult    int64         `json:"monkeyResult"`
	MonkeyTime      int64         `json:"monkeyTime"`
	UninstallResult int64         `json:"uninstallResult"`
	UninstallTime   int64         `json:"uninstallTime"`
	ErrorStatus     int64         `json:"errorStatus"`
	Cpu             []float64     `json:"cpu"`
	Memory          []int64       `json:"memory"`
	Flow            []interface{} `json:"flow"`
	Anr             []string      `json:"anr"`
	Crash           []string      `json:"crash"`
	Exception       []string      `json:"exception"`
}

//详细报告错误信息
type GetReportDetailErrorReqVo struct {
	Uid      int64 `json:"uid"`
	OrgId    int64 `json:"orgId"`
	ReportId int64 `json:"reportId"`
}

type GetReportDetailErrorListRespVo struct {
	vo.Err
	GetReportDetailErrorRespVo GetReportDetailErrorRespVo `json:"data"`
}

type GetReportDetailErrorRespVo struct {
	Total     int              `json:"total"`
	Anr       []GetReportDetailErrorInfo `json:"anr"`
	Crash     []GetReportDetailErrorInfo `json:"crash"`
	Exception []GetReportDetailErrorInfo `json:"exception"`
}

type GetReportDetailErrorInfo struct {
	ErrorMsg string `json:"errorMsg"`
	ErrorNum int `json:"errorNum"`
}

type GetReportDetailErrorVo struct {
	ErrorMsg string `json:"errorMsg"`
	ErrorNum int    `json:"errorNum"`
}

//
type GetReportDetailPerformanceReqVo struct {
	Uid      int64 `json:"uid"`
	OrgId    int64 `json:"orgId"`
	ReportId int64 `json:"reportId"`
}

type GetReportDetailPerformanceListRespVo struct {
	vo.Err
	GetReportDetailPerformanceRespVo GetReportDetailPerformanceRespVo `json:"data"`
}

type GetReportDetailPerformanceRespVo struct {
	Total                        int                            `json:"total"`
	GetReportDetailPerformanceVo []GetReportDetailPerformanceVo `json:"devicePerformanceArray"`
	PerformanceOverViewVo        PerformanceOverViewVo        `json:"performanceOverView"`
}

type GetReportDetailPerformanceVo struct {
	DeviceId      int64   `json:"deviceId"`
	DeviceName    string  `json:"deviceName"`
	InstallTime   int64   `json:"installTime"`
	LaunchTime    int64   `json:"launchTime"`
	UninstallTime int64   `json:"uninstallTime"`
	CpuAvgRate    float64 `json:"cpuAvgRate"`
	MemAvg        int64   `json:"memAvg"`
	FlowAvg       int64   `json:"flowAvg"`
}

type PerformanceOverViewVo struct {
	InstallTime   PerformanceOverViewInnerVo `json:"installTime"`
	LaunchTime    PerformanceOverViewInnerVo `json:"launchTime"`
	UninstallTime PerformanceOverViewInnerVo `json:"uninstallTime"`
	CpuRate       PerformanceOverViewInnerVo `json:"cpuRate"`
	Memory        PerformanceOverViewInnerVo `json:"memory"`
	Flow          PerformanceOverViewInnerVo `json:"flow"`
}

type PerformanceOverViewInnerVo struct {
	TopName     string `json:"topName"`
	TopValue    float64  `json:"topValue"`
	AvgValue    float64  `json:"avgValue"`
	BoottomName string `json:"boottomName"`
	BottomValue float64  `json:"bottomValue"`
}
