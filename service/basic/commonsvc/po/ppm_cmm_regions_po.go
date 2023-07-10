package po

type PpmCmmRegions struct {
	Id        int64  `db:"id,omitempty" json:"id"`
	CityId    int64  `db:"city_id,omitempty" json:"cityId"`
	Code      string `db:"code,omitempty" json:"code"`
	Name      string `db:"name,omitempty" json:"name"`
	Cname     string `db:"cname,omitempty" json:"cname"`
	LowerName string `db:"lower_name,omitempty" json:"lowerName"`
	CodeFull  string `db:"code_full,omitempty" json:"codeFull"`
	IsShow    int    `db:"is_show,omitempty" json:"isShow"`
	IsDefault int    `db:"is_default,omitempty" json:"isDefault"`
	IsDelete  int    `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmCmmRegions) TableName() string {
	return "ppm_cmm_regions"
}
