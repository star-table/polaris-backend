package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/star-table/polaris-backend/common/model/vo/orgvo"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/magiconair/properties/assert"
	"github.com/smartystreets/goconvey/convey"
)

func TestExportInviteTemplate(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		resp, err := ExportInviteTemplate(1574, 24370)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestGenDefaultNameByPhone(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		name1 := domain.GenDefaultNameByPhone("+86-15010011001")
		assert.Equal(t, name1, "用户1001")
		name1 = domain.GenDefaultNameByPhone("15010011002")
		assert.Equal(t, name1, "用户1002")
	}))
}

func TestImportMembers(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(ImportMembers(47964, 1516417, &orgvo.ImportMembersReqVoData{
			ImportUserList: []orgvo.ImportMembersReqVoDataUserItem{
				{
					Name:   "张六",
					Origin: "+86",
					Phone:  "18917640567",
				},
			},
		}))
	}))
}
