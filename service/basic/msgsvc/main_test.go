package main

import (
	"github.com/star-table/polaris-backend/common/extra/gin/mvc"
	"github.com/star-table/polaris-backend/service/basic/msgsvc/api"
	"testing"
)

func TestBuild(t *testing.T){
	port := 9090
	host := "127.0.0.1"
	version := "v1"
	applicationName := "msgsvc"

	postGreeter := api.PostGreeter{Greeter: mvc.NewPostGreeter(applicationName, host, port, version)}
	getGreeter := api.GetGreeter{Greeter: mvc.NewGetGreeter(applicationName, host, port, version)}
	facadeBuilder := mvc.FacadeBuilder{
		StorageDir: "F:\\workspace-test",
		Package: "facade",
		VoPackage: "idvo",
		Greeters: []interface{}{&postGreeter, &getGreeter},
	}

	facadeBuilder.Build()
}