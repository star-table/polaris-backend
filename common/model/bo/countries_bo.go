package bo

type CountriesBo struct {
	Id          int64  `db:"id,omitempty" json:"id"`
	ContinentId int64  `db:"continent_id,omitempty" json:"continentId"`
	Code        string `db:"code,omitempty" json:"code"`
	Name        string `db:"name,omitempty" json:"name"`
	FullName    string `db:"full_name,omitempty" json:"fullName"`
	Cname       string `db:"cname,omitempty" json:"cname"`
	FullCname   string `db:"full_cname,omitempty" json:"fullCname"`
	LowerName   string `db:"lower_name,omitempty" json:"lowerName"`
	Remark      string `db:"remark,omitempty" json:"remark"`
	IsShow      int    `db:"is_show,omitempty" json:"isShow"`
	IsDefault   int    `db:"is_default,omitempty" json:"isDefault"`
	IsDelete    int    `db:"is_delete,omitempty" json:"isDelete"`
}

func (*CountriesBo) TableName() string {
	return "ppm_cmm_countries"
}
