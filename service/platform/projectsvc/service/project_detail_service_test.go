package service

import (
	"fmt"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/consts"
	bo2 "github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"testing"
	"time"
)

func TestUpdateProjectDetail(t *testing.T) {

	bo := &bo2.ProjectDetailBo{}
	bo.UpdateTime = time.Now()
	po := &po.PpmProProjectDetail{}
	mayBlank, _ := time.Parse(consts.AppTimeFormat, "")
	fmt.Println(mayBlank.Format(consts.AppTimeFormat))
	copyer.Copy(bo, po)

	t.Log(po)

}
