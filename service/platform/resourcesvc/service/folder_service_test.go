package service

import (
	"context"
	"fmt"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/test"
	"github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func TestCreateFolder2(t *testing.T) {
	var a = []int64{0,1}
	if ok,_:= slice.Contain(a, int64(0)); ok {
		fmt.Println(1)
	} else {
		fmt.Println(0)
	}
}
func TestCreateFolder(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		createBo := bo.CreateFolderBo{
			ProjectId: 10116,
			OrgId:     1016,
			Name:      "folder1-2",
			UserId:    1046,
			ParentId:  0,
		}
		res, err := CreateFolder(createBo)
		fmt.Println(res)
		fmt.Println(err)
	}))
}

func TestUpdateFolder(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		//var parentId int64 = 0
		name := "folder656"
		folerBo := bo.UpdateFolderBo{
			UserId:       1046,
			OrgId:        1016,
			FolderID:     1033,
			UpdateFields: []string{"name"},
			ProjectID:    10116,
			Name:         &name,
			//ParentID:&parentId,
		}
		respVo, err := UpdateFolder(folerBo)
		fmt.Println(err)
		fmt.Println(respVo.OldValue)
		fmt.Println(respVo.NewValue)
	}))
}

func TestDeleteFolder(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		folerBo := bo.DeleteFolderBo{
			UserId:    10209,
			OrgId:     10109,
			ProjectId: 10115,
			FolderIds: []int64{1020, 1021},
		}
		res, err := DeleteFolder(folerBo)
		fmt.Println(res)
		fmt.Println(err)
	}))
}

func TestGetFolder(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		var parentId int64 = 0
		folerBo := bo.GetFolderBo{
			UserId:    10209,
			OrgId:     10109,
			ParentId:  &parentId,
			ProjectId: 10115,
			Page:      1,
			Size:      2,
		}
		res, err := GetFolder(folerBo)
		for _, value := range res.List {
			fmt.Println(value)
		}
		fmt.Println(res.Total)
		fmt.Println(err)
	}))
}
