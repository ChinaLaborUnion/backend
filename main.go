package main

import (
	"grpc-demo/core/cache"
	viewbase "grpc-demo/core/view"
	"grpc-demo/models/db"
	"grpc-demo/utils"
	"grpc-demo/utils/middlewares"
	"grpc-demo/views"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
)

func initRouter(app *iris.Application) {
	views.PartyCourseRouters(app)
	views.RegisterPartyClassRouters(app)
	views.RegisterSignUpRoute(app)
	views.CoursePictureRouters(app)
	views.RegisterNewsRouters(app)
	views.RegisterNewsLabelRouters(app)
	views.RegisterGoodsRouters(app)
	views.RegisterOrderRouters(app)
	views.HomeWorkRoute(app)
	views.RegisterTransactionRouters(app)
	views.RegisterAccountRouters(app)
}


func main() {
	app := iris.New()
	// 注册控制器
	app.UseGlobal(middlewares.AbnormalHandle, middlewares.RequestLogHandle)
	hero.Register(viewbase.ViewBase)
	// 注册路由
	initRouter(app)
	// 初始化配置
	utils.InitGlobal()
	// 初始化数据库
	db.InitDB()
	// 初始化缓存
	//cache.InitDijan()
	cache.InitRedisPool()
	// 初始化任务队列
	//queue.InitTaskQueue()
	// 启动系统
	app.Run(iris.Addr(":80"), iris.WithoutServerError(iris.ErrServerClosed))
}
