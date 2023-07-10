package errs

import (
	"fmt"
	"github.com/star-table/common/core/errors"
)

type SystemErrorInfo errors.SystemErrorInfo

func BuildSystemErrorInfo(resultCodeInfo errors.ResultCodeInfo, e ...error) errors.SystemErrorInfo {
	return errors.BuildSystemErrorInfo(resultCodeInfo, e...)
}

func BuildSystemErrorInfoWithMessage(resultCodeInfo errors.ResultCodeInfo, message string) errors.SystemErrorInfo {
	resultCodeInfo.SetMessage(message)
	return resultCodeInfo
}

func BuildSystemErrorInfoWithPanicRecover(r interface{}, stack string) errors.SystemErrorInfo {
	resultCodeInfo := errors.PanicError
	resultCodeInfo.SetMessage(fmt.Sprintf("捕获到的错误：%s, 堆栈:%s\n", r, stack))
	return resultCodeInfo
}

func BuildSystemErrorInfoWithMessageAndCode(resultCodeInfo errors.ResultCodeInfo, code int, message string) errors.SystemErrorInfo {
	resultCodeInfo.SetCode(code)
	resultCodeInfo.SetMessage(message)
	return resultCodeInfo
}

//add system ResultCodeInfo
func AddResultCodeInfo(code int, message string, langCode string) errors.ResultCodeInfo {
	return errors.AddResultCodeInfo(code, message, langCode)
}

func GetResultCodeInfoByCode(code int) errors.ResultCodeInfo {
	return errors.GetResultCodeInfoByCode(code)
}
