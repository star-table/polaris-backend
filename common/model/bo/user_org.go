package bo

type UserOrganizationBo struct {
	Id          int64     `db:"id,omitempty" json:"id"`
	OrgId       int64     `db:"org_id,omitempty" json:"orgId"`
	UserId      int64     `db:"user_id,omitempty" json:"userId"`
	CheckStatus int       `db:"check_status,omitempty" json:"checkStatus"`
	UseStatus   int       `db:"use_status,omitempty" json:"useStatus"`
	Status      int       `db:"status,omitempty" json:"status"`
	Version     int       `db:"version,omitempty" json:"version"`
	IsDelete    int       `db:"is_delete,omitempty" json:"isDelete"`
}