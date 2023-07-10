package lang

import (
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/polaris-backend/common/core/consts"
)

func GetLang() string {
	return threadlocal.GetValue(consts.AppHeaderLanguage)
}

func IsEnglish() bool {
	return GetLang() == consts.LangEnglish
}
