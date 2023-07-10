package bo

type MobilePrefixBo struct {
	Id        int64  `db:"id,omitempty" json:"id"`
	En        string `db:"en,omitempty" json:"en"`
	Zh        string `db:"zh,omitempty" json:"zh"`
	Code      string `db:"code,omitempty" json:"code"`
	Locale    string `db:"locale,omitempty" json:"locale"`
	Preg      string `db:"preg,omitempty" json:"preg"`
	CountryId int64  `db:"country_id,omitempty" json:"countryId"`
	IsShow    int    `db:"is_show,omitempty" json:"isShow"`
	IsDefault int    `db:"is_default,omitempty" json:"isDefault"`
	IsDelete  int    `db:"is_delete,omitempty" json:"isDelete"`
}

func (*MobilePrefixBo) TableName() string {
	return "ppm_cmm_mobile_prefix"
}
