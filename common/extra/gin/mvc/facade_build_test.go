package mvc

import (
	"context"
	"github.com/star-table/polaris-backend/common/model/vo/idvo"
	"testing"
)

type PostGreeter struct {
	Greeter
}

func (PostGreeter) ApplyPrimaryId(ctx *context.Context, req idvo.ApplyPrimaryIdReqVo) idvo.ApplyPrimaryIdRespVo {
	return idvo.ApplyPrimaryIdRespVo{}
}

func TestFacadeBuilder_Build(t *testing.T) {

	postGreeter := PostGreeter{Greeter: NewPostGreeter("idsvc", "127.0.0.1", 8080, "v1")}

	facadeBuilder := FacadeBuilder{
		StorageDir: "F:\\workspace-test",
		Package: "facade",
		VoPackage: "idvo",
		Greeters: []interface{}{&postGreeter},
	}

	facadeBuilder.Build()
}