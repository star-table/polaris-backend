package mvc

import (
	"net/http"
	"reflect"
)

type Greeter struct {
	ApplicationName string
	Host string
	Port int
	Version string
	HttpType string
}

func NewGetGreeter(applicationName string, host string, port int, version string) Greeter{
	return NewGreeter(http.MethodGet, applicationName, host, port, version)
}

func NewPostGreeter(applicationName string, host string, port int, version string) Greeter{
	return NewGreeter(http.MethodPost, applicationName, host, port, version)
}

func NewGreeter(httpType string, applicationName string, host string, port int, version string) Greeter{
	return Greeter{
		HttpType: httpType,
		ApplicationName:applicationName,
		Host:host,
		Port:port,
		Version:version,
	}
}

func GetHttpType(greeter interface{}) string{
	greeterValue := reflect.ValueOf(greeter)
	if greeterValue.Kind() == reflect.Ptr{
		greeterValue = greeterValue.Elem()
	}
	return greeterValue.FieldByName("HttpType").String()
}