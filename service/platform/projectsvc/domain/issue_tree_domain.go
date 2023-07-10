package domain

import (
	"fmt"
	"github.com/star-table/polaris-backend/common/model/bo"
)

func has(v1 bo.IssueNode, vs []*bo.IssueNode) bool {
	var hasChildPath bool
	hasChildPath = false
	for _, v2 := range vs {
		v3 := *v2
		if fmt.Sprintf("%s%d,", v1.Path, v1.IssueId) == v3.Path {
			hasChildPath = true
			break
		}
	}
	return hasChildPath

}

func MakeTree(vs []*bo.IssueNode, node *bo.IssueNode) {
	childs := findChild(node, vs)
	for _, child := range childs {
		node.Child = append(node.Child, child)
		if has(*child, vs) {
			MakeTree(vs, child)
		}
	}

}

func findChild(v *bo.IssueNode, vs []*bo.IssueNode) (ret []*bo.IssueNode) {
	for _, v2 := range vs {
		if fmt.Sprintf("%s%d,", v.Path, v.IssueId) == v2.Path {
			ret = append(ret, v2)
		}
	}
	return
}
