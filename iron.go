package Iron

import (
	"reflect"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type (
	// Logger interface
	// 重建日志对象，可以设置日志输出格式
	Logger interface {
		echo.Logger
		SetHeader(header string)
	}

	// Application class
	Application struct {
		Echo        *echo.Echo
		Controllers map[string]interface{}
		Services    map[string]ServiceInterface
		Router      *Router
		Logger      Logger
	}

	// BootCallBackFunc 启动器的回掉方法
	BootCallBackFunc func(application *Application) error
)

// Use ...
// 使用某个中间件
func (app *Application) Use(middleware ...echo.MiddlewareFunc) {
	app.Echo.Use(middleware...)
}

// GetEchoLogger is a function to get echo logger
func (app *Application) GetEchoLogger() echo.Logger {
	// l := log.New("-");
	// app.Echo.Logger =
	// lo := app.Echo.Logger.(*log.Logger)
	// lo.SetHeader("hello")
	// app.Echo.Logger.(*log.Logger).SetHeader(`{"时间":"${time_rfc3339_nano}","level":"${level}","prefix":"${prefix}",` +
	// 	`"file":"${short_file}","line":"${line}"}`)
	return app.Echo.Logger
}

/*
getType 获取 Type
*/
func (app *Application) getType(typeOf interface{}) reflect.Type {
	return reflect.Indirect(reflect.ValueOf(typeOf)).Type()
}

/*
AddController 添加注册控制器
*/
func (app *Application) AddController(controller interface{}) {
	// fmt.Println(controller)
	// fmt.Println(this.getType(controller).Name())
	app.Controllers[app.getType(controller).Name()] = controller
}

/*
AddService 添加服务
*/
func (app *Application) AddService(service ServiceInterface) {
	app.Services[app.getType(service).Name()] = service
}

/*
initRouter 初始化路由
*/
func (app *Application) initRouter() {
	app.Router.ControllersIndex = app.Controllers
}

/*
Boot is Application booter
*/
func (app *Application) Boot(callback BootCallBackFunc) {
	callback(app)

	// 注入已经注册的service
	app.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(*Context)
			cc.services = app.Services
			return next(cc)
		}
	})

	app.initRouter()
}

/*
Start ...
启动 Echo 应用
*/
func (app *Application) Start(address string) {
	app.Echo.Start(address)
}

/*
BetterAppContext ...
封装了一个更好的上下文
*/
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

/*
AddGormToContext ...
将Gorm添加到上下文
*/
func AddGormToContext(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			c := ctx.(*Context)
			c.AddDBHelper("gorm", db)
			return next(c)
		}
	}
}

/*
AddDBHelperToContext ...
添加数据库帮助库到上下文
*/
func AddDBHelperToContext(name string, db interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			c := ctx.(*Context)
			c.AddDBHelper(name, db)
			return next(c)
		}
	}
}

/*
CreateApplication ...
创建应用程序
*/
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
		Logger:      log.New("-"),
		Controllers: make(map[string]interface{}),
		Services:    make(map[string]ServiceInterface),
	}
	return
}
