package po

import "time"

type PpmTreTrends struct {
	Id              int64     `db:"id,omitempty" json:"id"`
	OrgId           int64     `db:"org_id,omitempty" json:"orgId"`
	Uuid            string    `db:"uuid,omitempty" json:"uuid"`
	Module1         string    `db:"module1,omitempty" json:"module1"`
	Module2Id       int64     `db:"module2_id,omitempty" json:"module2Id"`
	Module2         string    `db:"module2,omitempty" json:"module2"`
	Module3Id       int64     `db:"module3_id,omitempty" json:"module3Id"`
	Module3         string    `db:"module3,omitempty" json:"module3"`
	OperCode        string    `db:"oper_code,omitempty" json:"operCode"`
	OperObjId       int64     `db:"oper_obj_id,omitempty" json:"operObjId"`
	OperObjType     string    `db:"oper_obj_type,omitempty" json:"operObjType"`
	OperObjProperty string    `db:"oper_obj_property,omitempty" json:"operObjProperty"`
	RelationObjId   int64     `db:"relation_obj_id,omitempty" json:"relationObjId"`
	RelationObjType string    `db:"relation_obj_type,omitempty" json:"relationObjType"`
	RelationType    string    `db:"relation_type,omitempty" json:"relationType"`
	NewValue        *string   `db:"new_value,omitempty" json:"newValue"`
	OldValue        *string   `db:"old_value,omitempty" json:"oldValue"`
	Ext             string    `db:"ext,omitempty" json:"ext"`
	Creator         int64     `db:"creator,omitempty" json:"creator"`
	CreateTime      time.Time `db:"create_time,omitempty" json:"createTime"`
	IsDelete        int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmTreTrends) TableName() string {
	return "ppm_tre_trends"
}

type PpmTreTrendsCount struct {
	Module3Id int64 `db:"module3_id,omitempty" json:"module3Id"`
	Total     int64 `db:"total,omitempty" json:"total"`
}
