package domain

import (
	"context"
	"testing"

	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestClearUserRoleList(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {

		t.Log(ClearUserRoleList(10101, 10201, 0))
		//t.Log(cache.HSet("polaris:rolesvc:org_10101:user_10201:user_role_list_hash:1", "hello", "world"))
		//t.Log(cache.HGet("polaris:rolesvc:org_10101:user_10201:user_role_list_hash:1", "hello"))
		//
		//t.Log(cache.HDel("polaris:rolesvc:org_10101:user_10201:user_role_list_hash:1", "hello"))
	}))
}
