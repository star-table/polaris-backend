package bo

type StatesBo struct {
	Id        int64  `db:"id,omitempty" json:"id"`
	CountryId int64  `db:"country_id,omitempty" json:"countryId"`
	Code      string `db:"code,omitempty" json:"code"`
	Name      string `db:"name,omitempty" json:"name"`
	Cname     string `db:"cname,omitempty" json:"cname"`
	LowerName string `db:"lower_name,omitempty" json:"lowerName"`
	CodeFull  string `db:"code_full,omitempty" json:"codeFull"`
	AreaId    int    `db:"area_id,omitempty" json:"areaId"`
	IsShow    int    `db:"is_show,omitempty" json:"isShow"`
	IsDefault int    `db:"is_default,omitempty" json:"isDefault"`
	IsDelete  int    `db:"is_delete,omitempty" json:"isDelete"`
}

func (*StatesBo) TableName() string {
	return "ppm_cmm_states"
}
