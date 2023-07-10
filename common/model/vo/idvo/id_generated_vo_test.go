package idvo

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"testing"
)

func Test_Convert(t *testing.T){
	obj := ApplyPrimaryIdRespVo{
		Id: 123,
		Err: vo.NewErr(errs.SystemError),
	}
	t.Log(json.ToJsonIgnoreError(obj))

	obj1 := vo.VoidErr{
		Err: vo.NewErr(errs.SystemError),
	}
	t.Log(json.ToJsonIgnoreError(obj1))
}