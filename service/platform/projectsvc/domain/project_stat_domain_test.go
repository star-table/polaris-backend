package domain

import (
	"github.com/star-table/common/core/util/tests"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAppendProjectDayStat(t *testing.T) {

	convey.Convey("Test 项目任务燃尽图", t, tests.StartUp(func() {
		convey.Convey("Test 项目任务燃尽图", func() {
			projectBo := bo.ProjectBo{
				Id:    1010,
				OrgId: 17,
			}
			convey.So(AppendProjectDayStat(projectBo, consts.BlankDate), convey.ShouldEqual, nil)
		})
	}))

}
