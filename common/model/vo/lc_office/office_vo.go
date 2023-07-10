package lc_office

import "github.com/star-table/polaris-backend/common/model/vo"

type GetOfficeConfig struct {
	Type           string `json:"Type,omitempty"`           // 类型 type: collabora, microsoft
	UrlPrefix      string `json:"UrlPrefix,omitempty"`      // OfficeUrl前缀
	FileExtensions string `json:"FileExtensions,omitempty"` // 文件后嘴
}

type GetOfficeConfigRespVo struct {
	vo.Err
	Data *GetOfficeConfig `json:"data"`
}
