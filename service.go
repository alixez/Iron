package Iron

// Service class
type Service struct {
	model   string
	context *Context
}

// Init service object function
func (service *Service) Init(ctx *Context) {
	service.context = ctx
}

// GetModel function
func (service *Service) GetModel() string {
	return service.model
}

// ServiceInterface 接口
type ServiceInterface interface {
	Init(context *Context)
	GetModel() string
}
