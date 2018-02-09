package Iron

import (
	"net/http"
)

type Controller struct {
	context *Context
}

type ControllerInterface interface {
	Init()
	GetContext() *Context
	APISuccess() error
}

// Init funciton
func (this *Controller) Init(context *Context) {
	this.context = context
}

func (this *Controller) GetContext() *Context {
	return this.context
}

func (this *Controller) APISuccess(data interface{}) error {
	this.context.apiResponse.Code = 1
	this.context.apiResponse.SubCode = "default.void.success"
	this.context.apiResponse.Message = "接口请求成功"
	if data != nil {

		this.context.apiResponse.Data = data
	} else {
		this.context.apiResponse.Data = map[string]string{}
	}

	return this.context.JSON(http.StatusOK, this.context.apiResponse)
}

func (this *Controller) APIReturnWithOptional(code int, subCode string, message string, data interface{}) error {
	this.context.apiResponse.Code = code
	this.context.apiResponse.SubCode = subCode
	this.context.apiResponse.Message = message
	if data != nil {
		this.context.apiResponse.Data = data
	} else {
		this.context.apiResponse.Data = map[string]string{}
	}

	return this.context.JSON(http.StatusOK, this.context.apiResponse)
}

func (this *Controller) APIFail() error {
	this.context.apiResponse.Code = 0
	this.context.apiResponse.SubCode = "default.void.fail"
	this.context.apiResponse.Message = "接口请求失败"
	this.context.apiResponse.Data = map[string]string{}

	return this.context.JSON(http.StatusOK, this.context.apiResponse)
}
