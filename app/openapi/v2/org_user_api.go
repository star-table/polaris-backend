package v2

import (
	"github.com/star-table/polaris-backend/app/openapi"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	authData, err := openapi.ParseOpenAuthInfo(c)
	if err != nil {
		openapi.Fail(c, err)
		return
	}
	page := 1
	size := 10
	if c.Query("page") != "" {
		page = openapi.ParseInt(c.Query("page"))
	}
	if c.Query("size") != "" {
		size = openapi.ParseInt(c.Query("size"))
	}
	input := &vo.OrgUserListReq{
		CheckStatus: []int{2},
	}
	nameTemp := c.Query("name")
	if nameTemp != "" {
		input.Name = &nameTemp
	}
	emailTemp := c.Query("email")
	if emailTemp != "" {
		input.Email = &emailTemp
	}
	mobileTemp := c.Query("mobile")
	if mobileTemp != "" {
		input.Mobile = &mobileTemp
	}

	respVo := orgfacade.OpenOrgUserList(orgvo.OpenOrgUserListReqVo{
		Page:   page,
		Size:   size,
		UserId: 0,
		OrgId:  authData.OrgID,
		Input:  input,
	})
	if respVo.Failure() {
		openapi.Fail(c, respVo.Error())
	} else {
		openapi.Suc(c, respVo.Data)
	}
}
