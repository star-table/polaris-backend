package openapi

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	authData, err := ParseOpenAuthInfo(c)
	if err != nil {
		Fail(c, err)
		return
	}
	page := 1
	size := 10
	if c.Query("page") != "" {
		page = ParseInt(c.Query("page"))
	}
	if c.Query("size") != "" {
		size = ParseInt(c.Query("size"))
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
		Fail(c, respVo.Error())
	} else {
		Suc(c, respVo.Data)
	}
}
