package yidun

import (
	"errors"
	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/extra/yidun"
)

var verify *yidun.Verifier

func Verify(validate, user string) (*yidun.VerifyResult, error){
	ydConf := config.GetConfig().YiDun
	if ydConf == nil{
		return nil, errors.New("yidun config is nil.")
	}
	var err error
	verify, err = yidun.New(*ydConf)
	if err != nil{
		return nil, err
	}

	return verify.Verify(validate, user)
}
