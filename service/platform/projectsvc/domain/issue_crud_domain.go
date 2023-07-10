package domain

import (
	"fmt"

	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/businees"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/spf13/cast"
)

// CollectUserDeptIds 收集成员/部门Id
func CollectUserDeptIds(data map[string]interface{}, tableColumns map[string]*projectvo.TableColumnData) ([]int64, []int64) {
	var userIds, deptIds []int64
	for k, v := range data {
		if header, ok := tableColumns[k]; ok {
			if header.Field.Type == consts.LcColumnFieldTypeMember {
				newMemberIds := businees.LcMemberToUserIds(cast.ToStringSlice(v))
				userIds = append(userIds, newMemberIds...)
			} else if header.Field.Type == consts.LcColumnFieldTypeDept {
				newDept := cast.ToStringSlice(v)
				var newDeptIds []int64
				for _, id := range newDept {
					newDeptIds = append(newDeptIds, cast.ToInt64(id))
				}
				deptIds = append(deptIds, newDeptIds...)
			}
		}
	}
	return userIds, deptIds
}

func AssembleUserDeptsByIds(orgId int64, userIds, deptIds []int64) map[string]*uservo.MemberDept {
	users := make(map[int64]*uservo.GetAllUserByIdsRespDataUser)
	depts := make(map[int64]*uservo.GetAllDeptByIdsRespDataUser)

	userIds = slice.SliceUniqueInt64(userIds)
	deptIds = slice.SliceUniqueInt64(deptIds)
	if len(userIds) > 0 {
		resp := userfacade.GetAllUserByIds(orgId, userIds)
		if resp.Failure() {
			log.Errorf("[AssembleUserDeptsByIds] 获取用户失败 org:%d users:%q, err: %v", orgId, userIds, resp.Error())
		} else {
			for _, user := range resp.Data {
				u := user
				users[u.Id] = &u
			}
		}
	}
	if len(deptIds) > 0 {
		resp := userfacade.GetAllDeptByIds(orgId, deptIds)
		if resp.Failure() {
			log.Errorf("[AssembleUserDeptsByIds] 获取部门失败 org:%d depts:%q, err: %v", orgId, deptIds, resp.Error())
		} else {
			for _, dept := range resp.Data {
				d := dept
				depts[d.Id] = &d
			}
		}
	}

	return AssembleUserDepts(users, depts)
}

func AssembleDataIds(data map[string]interface{}) {
	data[consts.BasicFieldDataId] = cast.ToString(data[consts.BasicFieldId])
	data[consts.BasicFieldId] = data[consts.BasicFieldIssueId]
}

func AssembleUserDepts(users map[int64]*uservo.GetAllUserByIdsRespDataUser, depts map[int64]*uservo.GetAllDeptByIdsRespDataUser) map[string]*uservo.MemberDept {
	userDepts := make(map[string]*uservo.MemberDept)
	for _, u := range users {
		key := fmt.Sprintf("U_%d", u.Id)
		userDepts[key] = &uservo.MemberDept{
			Id:       u.Id,
			Name:     u.Name,
			Avatar:   u.Avatar,
			Type:     consts.LcCustomFieldUserType,
			Status:   u.Status,
			IsDelete: u.IsDelete,
		}
	}
	for _, d := range depts {
		key := fmt.Sprintf("D_%d", d.Id)
		userDepts[key] = &uservo.MemberDept{
			Id:       d.Id,
			Name:     d.Name,
			Avatar:   "",
			Type:     consts.LcCustomFieldDeptType,
			Status:   d.Status,
			IsDelete: d.IsDelete,
		}
	}
	return userDepts
}

func getUserDepts(orgId int64, userIds []int64, deptIds []int64) (map[string]*uservo.MemberDept, errs.SystemErrorInfo) {
	users := make(map[int64]*uservo.GetAllUserByIdsRespDataUser, len(userIds))
	depts := make(map[int64]*uservo.GetAllDeptByIdsRespDataUser, len(deptIds))
	if len(userIds) > 0 {
		resp := userfacade.GetAllUserByIds(orgId, userIds)
		if resp.Failure() {
			log.Errorf("[getUserDepts] GetAllUserByIds err: %v", resp.Error())
			return nil, resp.Error()
		}
		for i := 0; i < len(resp.Data); i++ {
			users[resp.Data[i].Id] = &resp.Data[i]
		}
	}

	if len(deptIds) > 0 {
		resp := userfacade.GetAllDeptByIds(orgId, deptIds)
		if resp.Failure() {
			log.Errorf("[BatchCreateIssue] GetAllDeptByIds err: %v", resp.Error())
			return nil, resp.Error()
		}

		for i := 0; i < len(resp.Data); i++ {
			depts[resp.Data[i].Id] = &resp.Data[i]
		}
	}

	return AssembleUserDepts(users, depts), nil
}

// 无码数据：组装成员、部门信息
func AssembleLcDataRelated(data map[string]interface{},
	tableColumns map[string]*projectvo.TableColumnData,
	users map[int64]*uservo.GetAllUserByIdsRespDataUser,
	depts map[int64]*uservo.GetAllDeptByIdsRespDataUser) {

	// 组装
	for k, v := range data {
		if v == nil {
			continue
		}

		if column, ok := tableColumns[k]; ok {
			switch column.Field.Type {
			case consts.LcColumnFieldTypeMember:
				ids, isSlice, _, _ := businees.LcInterfaceToIds(consts.LcCustomFieldUserType, v, true, true)
				us := make([]*uservo.MemberDept, 0)
				for _, id := range ids {
					if u, ok := users[id]; ok {
						us = append(us, &uservo.MemberDept{
							Id:       u.Id,
							Name:     u.Name,
							Avatar:   u.Avatar,
							Type:     consts.LcCustomFieldUserType,
							Status:   u.Status,
							IsDelete: u.IsDelete,
						})
					}
				}
				if isSlice {
					data[k] = us
				} else if len(us) > 0 {
					data[k] = us[0]
				}
			case consts.LcColumnFieldTypeDept:
				ids := cast.ToStringSlice(v)
				ds := make([]*uservo.MemberDept, 0)
				for _, id := range ids {
					i := cast.ToInt64(id)
					if d, ok := depts[i]; ok {
						ds = append(ds, &uservo.MemberDept{
							Id:       d.Id,
							Name:     d.Name,
							Avatar:   "",
							Type:     consts.LcCustomFieldDeptType,
							Status:   d.Status,
							IsDelete: d.IsDelete,
						})
					}
				}
				data[k] = ds
			}
		}
	}

	data[consts.BasicFieldDataId] = cast.ToString(data[consts.BasicFieldId])
	data[consts.BasicFieldId] = data[consts.BasicFieldIssueId]
}
