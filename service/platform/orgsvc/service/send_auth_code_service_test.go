package service

import (
	"testing"
	"time"

	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/dchest/captcha"
	"github.com/smartystreets/goconvey/convey"
	"gotest.tools/assert"
)

func TestSendAuthCode(t *testing.T) {
	convey.Convey("TestSendAuthCode", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		err := SendAuthCode(orgvo.SendAuthCodeReqVo{
			Input: vo.SendAuthCodeReq{
				AuthType:    1,
				AddressType: 2,
				Address:     "ainililia@163.com",
			},
		})
		assert.Equal(t, err, nil)
		time.Sleep(3 * time.Second)

		err = SendAuthCode(orgvo.SendAuthCodeReqVo{
			Input: vo.SendAuthCodeReq{
				AuthType:    2,
				AddressType: 2,
				Address:     "ainililia@163.com",
			},
		})
		assert.Equal(t, err, nil)
		time.Sleep(3 * time.Second)

		err = SendAuthCode(orgvo.SendAuthCodeReqVo{
			Input: vo.SendAuthCodeReq{
				AuthType:    3,
				AddressType: 2,
				Address:     "ainililia@163.com",
			},
		})
		assert.Equal(t, err, nil)
		time.Sleep(3 * time.Second)

		err = SendAuthCode(orgvo.SendAuthCodeReqVo{
			Input: vo.SendAuthCodeReq{
				AuthType:    4,
				AddressType: 2,
				Address:     "ainililia@163.com",
			},
		})
		assert.Equal(t, err, nil)
		time.Sleep(3 * time.Second)

		err = SendAuthCode(orgvo.SendAuthCodeReqVo{
			Input: vo.SendAuthCodeReq{
				AuthType:    5,
				AddressType: 2,
				Address:     "ainililia@163.com",
			},
		})
		assert.Equal(t, err, nil)
		time.Sleep(3 * time.Second)

		err = SendAuthCode(orgvo.SendAuthCodeReqVo{
			Input: vo.SendAuthCodeReq{
				AuthType:    6,
				AddressType: 2,
				Address:     "ainililia@163.com",
			},
		})
		assert.Equal(t, err, nil)
		time.Sleep(3 * time.Second)

	}))

}

func TestGetBaseOrgInfo(t *testing.T) {
	t.Log(captcha.VerifyString("rPZWi7iGRiFavHNK13q5", "899016"))
}
