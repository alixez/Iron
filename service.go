package Iron

import "github.com/labstack/gommon/log"

// Service class
type Service struct {
	Logger  Logger
	context *Context
}

// ServiceInterface 接口
type ServiceInterface interface {
	Init(context *Context)
	// GetModel() string
	GetContext() *Context
	GetService(name string) interface{}
	GetDB(name string) interface{}
}

// Init service object function
func (service *Service) Init(ctx *Context) {
	service.context = ctx
	if service.Logger == nil {
		service.Logger = log.New("Service")
	}
}

// GetContext 获取应用上下文
func (service *Service) GetContext() *Context {
	return service.context
}

// // GetModel function
// func (service *Service) GetModel() string {
// 	return service.model
// }

// GetDB is a function to get DB
func (service *Service) GetDB(name string) interface{} {
	return service.context.GetDB(name)
}

// GetService 获取某个service
func (service *Service) GetService(name string) interface{} {
	currentSrv := service.context.GetService(name)
	if nil == currentSrv {
		return nil
	}
	return currentSrv
}
