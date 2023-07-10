package bo

type ContinentsBo struct {
	Id        int64  `db:"id,omitempty" json:"id"`
	Name      string `db:"name,omitempty" json:"name"`
	Cname     string `db:"cname,omitempty" json:"cname"`
	LowerName string `db:"lower_name,omitempty" json:"lowerName"`
	IsShow    int    `db:"is_show,omitempty" json:"isShow"`
	IsDefault int    `db:"is_default,omitempty" json:"isDefault"`
	IsDelete  int    `db:"is_delete,omitempty" json:"isDelete"`
}

func (*ContinentsBo) TableName() string {
	return "ppm_cmm_continents"
}
