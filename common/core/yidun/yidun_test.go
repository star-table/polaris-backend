package yidun

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestVerify(t *testing.T) {
	convey.Convey("Test IssueRelationTypeOwner", t, func() {
		t.Log(Verify("1", "1"))
	})
}
