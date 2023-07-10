package domain

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

func AssemblyAttachmentList(resourceList []*vo.Resource, issueBoList []*bo.IssueBo, relationMap map[int64][]int64) ([]bo.AttachmentBo, errs.SystemErrorInfo) {
	issueBoMap := map[int64]*bo.IssueBo{}
	for _, value := range issueBoList {
		issueBoMap[value.Id] = value
	}
	bos := make([]bo.AttachmentBo, 0)
	for _, resource := range resourceList {
		attachmentBo := bo.AttachmentBo{
			Resource:  *resource,
			IssueList: []bo.IssueBo{},
		}

		if issueIds, ok := relationMap[resource.ID]; ok {
			for _, id := range issueIds {
				if value, ok := issueBoMap[id]; ok {
					attachmentBo.IssueList = append(attachmentBo.IssueList, *value)
				}
			}
		}

		bos = append(bos, attachmentBo)
	}

	//for resourceId, issueIdList := range relationMap {
	//	attachmentBo := bo.AttachmentBo{}
	//	if value, ok := resourceBoMap[resourceId]; ok {
	//		attachmentBo.Resource = value
	//	} else {
	//		log.Errorf("ResourceId %d not in resourceBoMap", resourceId)
	//		return nil, errs.ResourceNotExist
	//	}
	//	issueBoList := make([]bo.IssueBo, 0)
	//	for _, issueId := range issueIdList {
	//		if value, ok := issueBoMap[issueId]; ok {
	//			issueBo := value
	//			issueBoList = append(issueBoList, issueBo)
	//		} else {
	//			log.Errorf("IssueId %d not in issueBoMap", issueId)
	//			//return nil, errs.IssueNotExist
	//			continue
	//		}
	//	}
	//	attachmentBo.IssueList = issueBoList
	//	bos = append(bos, attachmentBo)
	//}
	return bos, nil
}

func AssemblyAttachment(resourceList *vo.Resource, issueBo *bo.IssueBo, relationMap map[int64][]int64) (*bo.AttachmentBo, errs.SystemErrorInfo) {
	return nil, nil
}
