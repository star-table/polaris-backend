package bo

import "time"

type PermissionOperationBo struct {
	Id             int64     `db:"id,omitempty" json:"id"`
	OrgId          int64     `db:"org_id,omitempty" json:"orgId"`
	PermissionId   int64     `db:"permission_id,omitempty" json:"permissionId"`
	LangCode       string    `db:"lang_code,omitempty" json:"langCode"`
	Name           string    `db:"name,omitempty" json:"name"`
	OperationCodes string    `db:"operation_codes,omitempty" json:"operationCodes"`
	Remark         string    `db:"remark,omitempty" json:"remark"`
	IsShow         int       `db:"is_show,omitempty" json:"isShow"`
	Status         int       `db:"status,omitempty" json:"status"`
	Creator        int64     `db:"creator,omitempty" json:"creator"`
	CreateTime     time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator        int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime     time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version        int       `db:"version,omitempty" json:"version"`
	IsDelete       int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PermissionOperationBo) TableName() string {
	return "ppm_rol_permission_operation"
}

type PermissionOperationListBo struct {
	PermissionInfo PermissionBo            `json:"permissionInfo"`
	OperationList  []PermissionOperationBo `json:"operationList"`
	PermissionHave []int64                 `json:"permissionHave"`
}
