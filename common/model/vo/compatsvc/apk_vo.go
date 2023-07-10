package compatsvc

import "github.com/star-table/polaris-backend/common/model/vo"

type UploadApkInfoReqVo struct {
	Uid     int64  `json:"uid"`
	OrgId   int64  `json:"orgId"`
	ApkName string `json:"apkName"`
	Url     string `json:"url"`
	Size    int64  `json:"size"`
}

type GetAllApkReqVo struct {
	Uid   int64 `json:"uid"`
	OrgId int64 `json:"orgId"`
	Page  int64 `json:"page"`
	Size  int64 `json:"size"`
}

type DeleteApkReqVo struct {
	Uid   int64 `json:"uid"`
	OrgId int64 `json:"orgId"`
	ApkId int64 `json:"apkId"`
}

type UploadApkInfoRespVo struct {
	vo.Err
}

type DeleteApkRespVo struct {
	vo.Err
}

//APK返回结果映射
type GetAllApkRespListVo struct {
	vo.Err
	GetAllApkRespVo GetAllApkRespVo `json:"data"`
}

type GetAllApkRespVo struct {
	ApkArray []GetAllApkArrayRespVo `json:"apkArray"`
	Total    int                    `json:"total"`
}

type GetAllApkArrayRespVo struct {
	ApkId          int64  `json:"apkId"`
	ApkName        string `json:"apkName"`
	PackageName    string `json:"packageName"`
	PackageVersion string `json:"packageVersion"`
	MainActivity   string `json:"mainActivity"`
	Creator        int64  `json:"creator"`
	CreatorName    string `json:"creatorName"`
	CreateTime     string `json:"createTime"`
}
