package compatsvc

import "github.com/star-table/polaris-backend/common/model/vo"

type GetAllDevicesRespListVo struct {
	vo.Err
	GetAllDevicesRespVo GetAllDevicesRespVo `json:"data"`
}

type GetAllDevicesRespVo struct {
	DeviceArray []GetAllDevicesArrayRespVo `json:"deviceArray"`
	Total       int                        `json:"total"`
}

type GetAllDevicesArrayRespVo struct {
	DeviceId      int64  `json:"deviceId"`
	Serial        string `json:"serial"`
	DeviceName    string `json:"deviceName"`
	Width         int64  `json:"width"`
	Height        int64  `json:"height"`
	Manufacturer  string `json:"manufacturer"`
	Model         string `json:"model"`
	Platform      string `json:"platform"`
	DeviceVersion string `json:"deviceVersion"`
	IsPresent     int64  `json:"isPresent"`
	IsInUse       int64  `json:"isInUse"`
}

//
type GetDeviceFiltrateRespVo struct {
	vo.Err
	GetDeviceFiltrateVo GetDeviceFiltrateVo `json:"data"`
}

type GetDeviceFiltrateVo struct {
	Manufacturer  []string `json:"manufacturer"`
	DeviceVersion []string `json:"deviceVersion"`
	ScreenSize    []string `json:"screenSize"`
}

//
type GetAllDevicesStatusListRespVo struct {
	vo.Err
	GetAllDevicesStatusRespVo GetAllDevicesStatusRespVo `json:"data"`
}

type GetAllDevicesStatusRespVo struct {
	DeviceArray []GetAllDevicesStatusArrayRespVo `json:"deviceArray"`
	Total       int                              `json:"total"`
}

type GetAllDevicesStatusArrayRespVo struct {
	DeviceId  int64 `json:"deviceId"`
	IsPresent int64 `json:"isPresent"`
	IsInUse   int64 `json:"isInUse"`
}

//提测的req
type StartCompatReqVo struct {
	Uid           int64   `json:"uid"`
	OrgId         int64   `json:"orgId"`
	ApkId         int64   `json:"apkId"`
	DeviceIdArray []int64 `json:"deviceIdArray"`
}

type StartCompatRespVo struct {
	vo.Err
	StartCompatVo StartCompatVo `json:"data"`
}

type StartCompatVo struct {
	ReportId int64 `json:"reportId"`
}
