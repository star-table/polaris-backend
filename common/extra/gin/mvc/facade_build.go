package mvc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/star-table/polaris-backend/common/core/consts"

	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/temp"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/core/util/str"
)

const FacadeBuildGoFile = `package {{.Package}}

import ({{.CtxPack}}
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"{{.GinUtilPack}}
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/common/model/vo/{{.VoPackage}}"
	"github.com/star-table/polaris-backend/facade"
)

var _ = commonvo.DataEvent{}
{{.FuncDesc}}
`

const FacadeBuildFunc = `
func {{.MethodName}}({{.ArgsDesc}}) {{.OutArgsDesc}} {
	{{.OutArgsWithoutPtrDesc}}
	reqUrl := {{.Api}}
{{.QueryParamsDesc}}{{.RequestDesc}}
{{.ReturnDesc}}
}`

type FacadeBuilder struct {
	StorageDir string
	Package    string
	VoPackage  string
	Greeters   []interface{}
}

type ReqTypeInfo struct {
	reqType  reflect.Type
	isPtr    bool
	isString bool
}

var skipMethods = []string{"Health"}

func (fb FacadeBuilder) Build() {
	for _, greeter := range fb.Greeters {
		greeterInfo := GetGreeterInfo(greeter)
		goFileContent := fb.greeterRender(greeter, greeterInfo)

		if strings.Trim(goFileContent, " ") == "" {
			continue
		}
		version := strings.ToLower(greeterInfo.Version)
		if version != "" {
			version += "_"
		}
		filePath := fmt.Sprintf("%s/%s_%s_%sfacade.go", fb.StorageDir, greeterInfo.ApplicationName, strings.ToLower(greeterInfo.HttpType), version)

		data := []byte(goFileContent)
		if ioutil.WriteFile(filePath, data, 0644) == nil {
			fmt.Printf("build success: %s\n", filePath)
		}
	}
}

func (fb FacadeBuilder) greeterRender(greeter interface{}, greeterInfo Greeter) string {
	greeterValue := reflect.ValueOf(greeter)
	greeterType := reflect.TypeOf(greeter)

	funcDesc := ""
	methodNum := greeterType.NumMethod()

	if methodNum == 0 {
		return ""
	}
	hasCtx := false
	hasApi := false
	for i := 0; i < methodNum; i++ {
		methodName := greeterType.Method(i).Name
		if exist, _ := slice.Contain(skipMethods, methodName); exist {
			fmt.Println("skip method", methodName)
			continue
		}
		hasApi = true
		method := greeterValue.MethodByName(methodName)
		subFuncDesc, useCtx := fb.greeterMethodRender(method, methodName, greeterInfo)
		if useCtx {
			hasCtx = true
		}
		funcDesc += subFuncDesc + "\n"
	}
	if !hasApi {
		return ""
	}

	goFileTemplateData := map[string]string{}
	goFileTemplateData["Package"] = fb.Package
	goFileTemplateData["VoPackage"] = fb.VoPackage
	goFileTemplateData["FuncDesc"] = funcDesc
	if hasCtx {
		goFileTemplateData["CtxPack"] = "\n\t\"context\""
		goFileTemplateData["GinUtilPack"] = "\n\t\"github.com/star-table/polaris-backend/common/extra/gin/util\""
	} else {
		goFileTemplateData["CtxPack"] = ""
		goFileTemplateData["GinUtilPack"] = ""
	}
	result, _ := temp.Render(FacadeBuildGoFile, goFileTemplateData)
	return result
}

func (fb FacadeBuilder) greeterMethodRender(method reflect.Value, methodName string, greeterInfo Greeter) (string, bool) {
	templateData := map[string]string{}
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("methodName:%v is panic", methodName)
			log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
		}
	}()

	var hasQueryParams, hasRequestBody bool

	//MethodName
	templateData["MethodName"] = methodName + strings.ToUpper(greeterInfo.Version)

	//ArgsDesc
	ptrTab := map[string]ReqTypeInfo{}
	templateData["ArgsDesc"] = fb.buildArgDesc(method, ptrTab)

	_, useCtx := ptrTab["ctx"]

	//ReturnDesc
	templateData["OutArgsDesc"], templateData["OutArgsWithoutPtrDesc"] = fb.buildOutArgsDesc(method, ptrTab)

	//BuildApi
	templateData["Api"] = "fmt.Sprintf(\"%s/" + BuildApi(greeterInfo.ApplicationName, greeterInfo.Version, methodName) + "\", config.GetPreUrl(\"" + greeterInfo.ApplicationName + "\"))"

	//SvcName
	templateData["SvcName"] = greeterInfo.ApplicationName

	//QueryParamsDesc
	templateData["QueryParamsDesc"], hasQueryParams, hasRequestBody = fb.buildQueryParamsDesc(greeterInfo.HttpType, ptrTab)

	//RequestDesc
	templateData["RequestDesc"] = fb.buildRequestDesc(greeterInfo.HttpType, ptrTab, hasQueryParams, hasRequestBody)

	//ReturnDesc
	templateData["ReturnDesc"] = fb.buildReturnDesc(ptrTab)

	result, _ := temp.Render(FacadeBuildFunc, templateData)
	return result, useCtx
}

func GetGreeterInfo(greeter interface{}) Greeter {
	greeterValue := reflect.ValueOf(greeter)
	if greeterValue.Kind() == reflect.Ptr {
		greeterValue = greeterValue.Elem()
	}
	return greeterValue.FieldByName("Greeter").Interface().(Greeter)
}

func (fb FacadeBuilder) buildArgDesc(method reflect.Value, ptrTab map[string]ReqTypeInfo) string {
	methodType := method.Type()
	numIn := methodType.NumIn()
	argsDesc := ""
	for i := 0; i < numIn; i++ {
		inType := methodType.In(i)
		inTypeString := inType.String()

		isPtr := false
		if inType.Kind() == reflect.Ptr {
			inType = inType.Elem()
			isPtr = true
		}

		argName := ""
		if inType.String() == "context.Context" {
			argName = "ctx"
		} else if inType.Kind() == reflect.Struct {
			argName = "req"
		} else {
			argName = "arg" + strconv.Itoa(i)
		}
		ptrTab[argName] = ReqTypeInfo{isPtr: isPtr, reqType: inType}

		split := ", "
		if i == numIn-1 {
			split = ""
		}
		argsDesc += fmt.Sprintf("%s %s%s", argName, inTypeString, split)
	}
	return argsDesc
}

func (fb FacadeBuilder) buildOutArgsDesc(method reflect.Value, ptrTab map[string]ReqTypeInfo) (string, string) {
	methodType := method.Type()
	numOut := methodType.NumOut()
	//限制出参只能有一个
	if numOut > 1 {
		panic("Only one value can be returned")
	}

	outArgsDesc := ""
	isPtr := false
	isString := false
	for i := 0; i < numOut; i++ {
		outType := methodType.Out(i)
		outTypeString := outType.String()

		if outType.Kind() == reflect.Ptr {
			outType = outType.Elem()
			isPtr = true
		}

		argName := ""
		if outType.String() == "context.Context" {
			argName = "ctx"
		} else if outType.Kind() == reflect.Struct {
			argName = "resp"
		} else if outType.Kind() == reflect.String {
			isString = true
			isPtr = true
			argName = "resp"
		} else {
			argName = "ret" + strconv.Itoa(i)
		}
		ptrTab[argName] = ReqTypeInfo{isPtr: isPtr, isString: isString, reqType: outType}

		split := ", "
		if i == numOut-1 {
			split = ""
		}
		outArgsDesc += fmt.Sprintf("%s%s", outTypeString, split)
	}
	outArgsDescWithoutPtr := outArgsDesc
	if isPtr {
		outArgsDescWithoutPtr = outArgsDescWithoutPtr[1:]
	}
	if isString {
		outArgsDescWithoutPtr = "respVo := \"\""
	} else {
		outArgsDescWithoutPtr = "respVo := &" + outArgsDescWithoutPtr + "{}"
	}

	return outArgsDesc, outArgsDescWithoutPtr
}

func (fb FacadeBuilder) buildQueryParamsDesc(httpType string, ptrTab map[string]ReqTypeInfo) (string, bool, bool) {
	queryParamsDesc := ""
	requestBodyDesc := ""
	hasQueryParams := false
	hasRequestBody := false

	reqObj, ok := ptrTab["req"]
	if ok {
		reqType := reqObj.reqType
		numField := reqType.NumField()
		hasQueryParams, hasRequestBody = assemblyQueryDesc(numField, reqType, &queryParamsDesc, &requestBodyDesc)
	}

	return queryParamsDesc + requestBodyDesc, hasQueryParams, hasRequestBody
}

func (fb FacadeBuilder) buildRequestBodyLogDesc(httpType string) string {
	if httpType == http.MethodPost {
		return "requestBody"
	}
	return "\"\""
}

func assemblyQueryDesc(numField int, reqType reflect.Type, queryParamsDesc, queryBodyDesc *string) (bool, bool) {
	hasQueryParams := false
	hasRequestBody := false
	for i := 0; i < numField; i++ {
		field := reqType.Field(i)
		fieldName := field.Name

		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if _, ok := BasicTypeConverter[fieldType.String()]; ok {
			if *queryParamsDesc == "" {
				*queryParamsDesc += "\tqueryParams := map[string]interface{}{}\n"
			}
			*queryParamsDesc += fmt.Sprintf("\tqueryParams[\"%s\"] = %s.%s\n", str.LcFirst(fieldName), "req", fieldName)
			hasQueryParams = true
		} else if reqType.Kind() == reflect.Struct {
			if *queryBodyDesc != "" {
				panic("req body more than one")
			}
			*queryBodyDesc = fmt.Sprintf("\trequestBody := &%s.%s\n", "req", fieldName)
			hasRequestBody = true
		}
	}
	return hasQueryParams, hasRequestBody
}

func (fb FacadeBuilder) buildRequestDesc(httpType string, ptrTab map[string]ReqTypeInfo, hasQueryParams, hasRequestBody bool) string {
	headerOptionsDesc := ""
	requestDesc := ""
	httpRequestHeaderDesc := ""

	queryParamsDesc := "nil"
	if hasQueryParams {
		queryParamsDesc = "queryParams"
	}
	requestBodyDesc := "nil"
	if hasRequestBody {
		requestBodyDesc = "requestBody"
	}

	ctx, ok := ptrTab["ctx"]
	if ok {
		//respObj, ok := ptrTab["resp"]
		//if !ok {
		//	panic("no resp respObj")
		//}

		//respObjPtrFlag := "*"
		//if respObj.isPtr {
		//	respObjPtrFlag = ""
		//}

		ctxPtrFlag := ""
		if ctx.isPtr {
			ctxPtrFlag = "*"
		}
		headerOptionsDesc += fmt.Sprintf("\theaderOptions, _ := util.BuildHeaderOptions(%s%s)\n", ctxPtrFlag, "ctx")
		//headerOptionsDesc += fmt.Sprint("\tif err != nil{\n")
		//headerOptionsDesc += fmt.Sprint("\t\trespVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))\n")
		//headerOptionsDesc += fmt.Sprintf("\t\treturn %srespVo\n", respObjPtrFlag)
		//headerOptionsDesc += fmt.Sprint("\t}\n")
		httpRequestHeaderDesc += ", headerOptions"
	} else {
		httpRequestHeaderDesc += ", nil"
	}

	httpMethod := strings.ToUpper(httpType)
	httpMethodDesc := ""
	switch httpMethod {
	case consts.HttpMethodGet:
		httpMethodDesc = "consts.HttpMethodGet"
	case consts.HttpMethodPost:
		httpMethodDesc = "consts.HttpMethodPost"
	case consts.HttpMethodPut:
		httpMethodDesc = "consts.HttpMethodPut"
	case consts.HttpMethodDelete:
		httpMethodDesc = "consts.HttpMethodDelete"
	}

	respObj, ok := ptrTab["resp"]
	if !ok {
		panic("no resp respObj")
	}
	if respObj.isString {
		requestDesc = fmt.Sprintf("\terr := facade.Request(%s, reqUrl, %s%s, %s, &respVo)\n", httpMethodDesc, queryParamsDesc, httpRequestHeaderDesc, requestBodyDesc)
	} else {
		requestDesc = fmt.Sprintf("\terr := facade.Request(%s, reqUrl, %s%s, %s, respVo)\n", httpMethodDesc, queryParamsDesc, httpRequestHeaderDesc, requestBodyDesc)
	}

	requestDesc += "\tif err.Failure() {\n"
	if respObj.isString {
		requestDesc += "\t\treturn facade.ToJsonString(err)\n"
	} else {
		requestDesc += "\t\trespVo.Err = err\n"
	}
	requestDesc += "\t}"

	return headerOptionsDesc + requestDesc
}

func (fb FacadeBuilder) buildReturnDesc(ptrTab map[string]ReqTypeInfo) string {
	respObj, ok := ptrTab["resp"]
	if !ok {
		panic("no resp respObj")
	}
	respObjPtrFlag := "*"
	if respObj.isPtr {
		respObjPtrFlag = ""
	}

	return fmt.Sprintf("\treturn %srespVo", respObjPtrFlag)
}
