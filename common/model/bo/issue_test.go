package bo

//
//import (
//	"fmt"
//	"github.com/star-table/common/core/types"
//	"github.com/star-table/common/core/util/json"
//	"testing"
//	"time"
//)
//
//func TestBuildIssueBoFromIssuePo(t *testing.T) {
//	issuePo := po.PpmPriIssue{
//		Id:                  1,
//		OrgId:               1,
//		Code:                "CC-1",
//		ProjectId:           0,
//		TableId: 1234,
//		Title:               "1AAA",
//		Owner:               12,
//		PriorityId:          234,
//		SourceId:            345,
//		IssueObjectTypeId:   456,
//		PlanStartTime:       time.Now(),
//		PlanEndTime:         time.Now(),
//		StartTime:           time.Now(),
//		EndTime:             time.Now(),
//		PlanWorkHour:        4,
//		IterationId:         0,
//		VersionId:           0,
//		ModuleId:            0,
//		ParentId:            0,
//		Status:              567,
//		Creator:             678,
//		CreateTime:          time.Now(),
//		Updator:             789,
//		UpdateTime:          time.Now(),
//		Version:             1,
//		IsDelete:            2,
//	}
//
//	issueBo := &IssueBo{}
//	err := util.ConvertObject(&issuePo, &issueBo)
//
//	fmt.Println(json.ToJson(issueBo))
//	fmt.Println(err)
//
//	/**
//	{
//	    "id": 1,
//	    "orgId": 1,
//	    "code": "CC-1",
//	    "projectId": 0,
//	    "projectObjectTypeId": 1234,
//	    "title": "1AAA",
//	    "owner": 12,
//	    "priorityId": 234,
//	    "sourceId": 345,
//	    "issueObjectTypeId": 456,
//	    "planStartTime": "2019-08-06 14:10:21",
//	    "planEndTime": "2019-08-06 14:10:21",
//	    "startTime": "2019-08-06 14:10:21",
//	    "endTime": "2019-08-06 14:10:21",
//	    "planWorkHour": 4,
//	    "iterationId": 0,
//	    "versionId": 0,
//	    "moduleId": 0,
//	    "parentId": 0,
//	    "status": 567,
//	    "creator": 678,
//	    "createTime": "2019-08-06 14:10:21",
//	    "updator": 789,
//	    "updateTime": "2019-08-06 14:10:21",
//	    "version": 1,
//	    "issueDetail": null,
//	    "ownerInfo": null,
//	    "participantInfos": null,
//	    "followerInfos": null
//	}
//	*/
//}
//
//func TestBuildIssuePoFromBo(t *testing.T) {
//	issueBo := &IssueBo{
//		Id:                  1,
//		OrgId:               1,
//		Code:                "CC-1",
//		ProjectId:           0,
//		TableId: 1234,
//		Title:               "1AAA",
//		Owner:               12,
//		PriorityId:          234,
//		SourceId:            345,
//		IssueObjectTypeId:   456,
//		PlanStartTime:       types.NowTime(),
//		PlanEndTime:         types.NowTime(),
//		StartTime:           types.NowTime(),
//		EndTime:             types.NowTime(),
//		PlanWorkHour:        4,
//		IterationId:         0,
//		VersionId:           0,
//		ModuleId:            0,
//		ParentId:            0,
//		Status:              567,
//		Creator:             678,
//		CreateTime:          types.NowTime(),
//		Updator:             789,
//		UpdateTime:          types.NowTime(),
//		Version:             1,
//	}
//
//	issuePo := &po.PpmPriIssue{}
//	err := util.ConvertObject(&issueBo, &issuePo)
//
//	fmt.Println(json.ToJson(issuePo))
//	fmt.Println(err)
//	/**
//	{
//	    "Id": 1,
//	    "OrgId": 1,
//	    "Code": "CC-1",
//	    "ProjectId": 0,
//	    "TableId": 1234,
//	    "Title": "1AAA",
//	    "Owner": 12,
//	    "PriorityId": 234,
//	    "SourceId": 345,
//	    "IssueObjectTypeId": 456,
//	    "PlanStartTime": "2019-08-06T16:13:28Z",
//	    "PlanEndTime": "2019-08-06T16:13:28Z",
//	    "StartTime": "2019-08-06T16:13:28Z",
//	    "EndTime": "2019-08-06T16:13:28Z",
//	    "PlanWorkHour": 4,
//	    "IterationId": 0,
//	    "VersionId": 0,
//	    "ModuleId": 0,
//	    "ParentId": 0,
//	    "Status": 567,
//	    "Creator": 678,
//	    "CreateTime": "2019-08-06T16:13:28Z",
//	    "Updator": 789,
//	    "UpdateTime": "2019-08-06T16:13:28Z",
//	    "Version": 1,
//	    "IsDelete": 2
//	}
//	*/
//
//}
//
//func TestBuildIssueDetailBoFromIssueDetailPo(t *testing.T) {
//	issueDetailPo := po.PpmPriIssueDetail{
//		Id:         1,
//		OrgId:      2,
//		IssueId:    3,
//		ProjectId:  4,
//		StoryPoint: 0,
//		Tags:       "tags",
//		Remark:     "remark",
//		Status:     123,
//		Creator:    6,
//		CreateTime: time.Now(),
//		Updator:    7,
//		UpdateTime: time.Now(),
//		Version:    8,
//		IsDelete:   9,
//	}
//
//	issueDetailBo := &IssueDetailBo{}
//	err := util.ConvertObject(&issueDetailPo, &issueDetailBo)
//
//	fmt.Println(json.ToJson(issueDetailBo))
//	/**
//	{
//	    "id": 1,
//	    "orgId": 2,
//	    "issueId": 3,
//	    "projectId": 4,
//	    "storyPoint": 0,
//	    "tags": "tags",
//	    "remark": "remark",
//	    "status": 123,
//	    "creator": 6,
//	    "createTime": "2019-08-06 13:54:28",
//	    "updator": 7,
//	    "updateTime": "2019-08-06 13:54:28",
//	    "version": 8
//	}
//	*/
//
//	fmt.Println(err)
//}
//
//func TestBuildIssueRelationBoFromPo(t *testing.T) {
//	issueRelationPo := po.PpmPriIssueRelation{
//		Id:           1,
//		OrgId:        2,
//		IssueId:      3,
//		RelationId:   4,
//		RelationType: 5,
//		Creator:      6,
//		CreateTime:   time.Now(),
//		Updator:      8,
//		UpdateTime:   time.Now(),
//		Version:      10,
//		IsDelete:     2,
//	}
//
//	issueRelationBo := &IssueRelationBo{}
//	err := util.ConvertObject(&issueRelationPo, &issueRelationBo)
//
//	fmt.Println(json.ToJson(issueRelationBo))
//	fmt.Println(err)
//
//	/**
//	{
//	    "id": 1,
//	    "orgId": 2,
//	    "issueId": 3,
//	    "relationId": 4,
//	    "relationType": 5,
//	    "creator": 6,
//	    "createTime": "2019-08-06 14:00:11",
//	    "updator": 8,
//	    "updateTime": "2019-08-06 14:00:11",
//	    "version": 10
//	}
//	*/
//}
//
//func TestBuildIssueUserBoFromPo(t *testing.T) {
//	issueRelationPo := po.PpmPriIssueRelation{
//		Id:           1,
//		OrgId:        2,
//		IssueId:      3,
//		RelationId:   4,
//		RelationType: 5,
//		Creator:      6,
//		CreateTime:   time.Now(),
//		Updator:      8,
//		UpdateTime:   time.Now(),
//		Version:      10,
//		IsDelete:     2,
//	}
//
//	issueUserBo := &IssueUserBo{}
//
//	err := util.ConvertObject(&issueRelationPo, &issueUserBo)
//
//	fmt.Println(json.ToJson(issueUserBo))
//	fmt.Println(err)
//
//	// {"id":1,"orgId":2,"issueId":3,"relationId":4,"relationType":5,"creator":6,"createTime":"2019-08-06 14:06:11","updator":8,"updateTime":"2019-08-06 14:06:11","version":10}
//}
//
//func TestBuildIssueRelationBosFromPos(t *testing.T) {
//
//	//  issueUserBos := make([]IssueUserBo, 5)
//	//
//	//  for i := int64(0); i<5; i++ {
//	//
//	//	  issueUserBo := &IssueUserBo{}
//	//	  issueUserBo.Id = i
//	//	  issueUserBo.Name = "bbb"
//	//	  issueUserBo.OrgId = i
//	//	  issueUserBo.CreateTime = types.NowTime()
//	//	  issueUserBo.RelationId = i
//	//	  issueUserBo.UpdateTime = types.NowTime()
//	//	  issueUserBo.Version = 1
//	//
//	//	  fmt.Println(&issueUserBo.OrgId)
//	//	  fmt.Println(&issueUserBo.IssueRelationBo.OrgId)
//	//
//	//	  issueUserBos[i] = *issueUserBo
//	//  }
//	//  //issueRelationBos := []IssueRelationBo
//	//  fmt.Println(json.ToJson(issueUserBos))
//	//
//	//issueUserBos2 := []IssueUserBo{}
//	//
//	//  issueUserBos = append(issueUserBos, issueUserBos2...)
//	//
//	//  fmt.Println(json.ToJson(issueUserBos))
//
//	issueBo := &IssueBo{
//		Id:                  1,
//		OrgId:               1,
//		Code:                "CC-1",
//		ProjectId:           0,
//		TableId: 1234,
//		Title:               "1AAA",
//		Owner:               12,
//		PriorityId:          234,
//		SourceId:            345,
//		IssueObjectTypeId:   456,
//		PlanStartTime:       types.NowTime(),
//		PlanEndTime:         types.NowTime(),
//		StartTime:           types.NowTime(),
//		EndTime:             types.NowTime(),
//		PlanWorkHour:        4,
//		IterationId:         0,
//		VersionId:           0,
//		ModuleId:            0,
//		ParentId:            0,
//		Status:              567,
//		Creator:             678,
//		CreateTime:          types.NowTime(),
//		Updator:             789,
//		UpdateTime:          types.NowTime(),
//		Version:             1,
//	}
//	issueDetailBo := &IssueDetailBo{
//		Id:    2,
//		OrgId: 3,
//	}
//
//	fmt.Println(json.ToJson(issueBo))
//	fmt.Println(json.ToJson(issueDetailBo))
//
//}
