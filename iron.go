package Iron

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
)

type (
	Application struct {
		Echo        *echo.Echo
		Controllers map[string]interface{}
		Services    map[string]ServiceInterface
		Router      *Router
	}

	BootCallBackFunc func(application *Application) error
)

func (app *Application) Use(middleware ...echo.MiddlewareFunc) {
	app.Echo.Use(middleware...)
}

func (app *Application) GetLogger() echo.Logger {
	return app.Echo.Logger
}

func (this *Application) getType(typeOf interface{}) reflect.Type {
	return reflect.Indirect(reflect.ValueOf(typeOf)).Type()
}

func (this *Application) AddController(controller interface{}) {
	fmt.Println(controller)
	fmt.Println(this.getType(controller).Name())
	this.Controllers[this.getType(controller).Name()] = controller
}

func (this *Application) AddService(service ServiceInterface) {
	this.Services[this.getType(service).Name()] = service
}

func (this *Application) initRouter() {
	this.Router.ControllersIndex = this.Controllers
}

func (this *Application) Boot(callback BootCallBackFunc) {
	callback(this)

	// 注入已经注册的service
	this.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(*Context)
			cc.services = this.Services
			return next(cc)
		}
	})

	this.initRouter()
}

func (this *Application) Start(address string) {
	this.Echo.Start(address)
}

func BetterAppContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &Context{
			c,
			nil,
			&APIResponse{
				Code:    0,
				Message: "空",
				SubCode: "default.void.default",
				Data:    map[string]interface{}{},
			},
			map[string]interface{}{},
			nil,
		}

		return next(cc)
	}
}

func AddGormToContext(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			c := ctx.(*Context)
			c.AddDBHelper("gorm", db)
			return next(c)
		}
	}
}

func AddDBHelperToContext(name string, db interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			c := ctx.(*Context)
			c.AddDBHelper(name, db)
			return next(c)
		}
	}
}

func CreateApplication(env *Env) (application *Application) {
	e := echo.New()
	e.Use(BetterAppContext)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			c := ctx.(*Context)
			c.Config = env
			return next(c)
		}
	})

	router := &Router{
		Echo: e,
	}
	application = &Application{
		Echo:        e,
		Router:      router,
		Controllers: make(map[string]interface{}),
		Services:    make(map[string]ServiceInterface),
	}
	return
}
