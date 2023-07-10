package bo

import "time"

type ResourceBo struct {
	Id          int64     `db:"id,omitempty" json:"id"`
	OrgId       int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId   int64     `db:"project_id,omitempty" json:"projectId"`
	Type        int       `db:"type,omitempty" json:"type"`
	Bucket      string    `db:"bucket,omitempty" json:"bucket"`
	Host        string    `db:"host,omitempty" json:"host"`
	Path        string    `db:"path,omitempty" json:"path"`
	Name        string    `db:"name,omitempty" json:"name"`
	Suffix      string    `db:"suffix,omitempty" json:"suffix"`
	Md5         string    `db:"md5,omitempty" json:"md5"`
	Size        int64     `db:"size,omitempty" json:"size"`
	FileType    int       `db:"file_type,omitempty" json:"fileType"`
	SourceType  int       `db:"source_type,omitempty" json:"sourceType"`
	Creator     int64     `db:"creator,omitempty" json:"creator"`
	CreateTime  time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator     int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime  time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version     int       `db:"version,omitempty" json:"version"`
	IsDelete    int       `db:"is_delete,omitempty" json:"isDelete"`
	CreatorName string    `json:"creatorName"`
}

func (*ResourceBo) TableName() string {
	return "ppm_res_resource"
}
