package main

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"grpc-demo/core/cache"
	viewbase "grpc-demo/core/view"
	"grpc-demo/models/db"
	"grpc-demo/utils"
	logUtils "grpc-demo/utils/log"
	"grpc-demo/utils/middlewares"
	"grpc-demo/views"
	"time"

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
	views.RegisterSignInRoute(app)
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
	go signIn()
	// 启动系统
	app.Run(iris.Addr(":80"), iris.WithoutServerError(iris.ErrServerClosed))
}

func signIn(){
		//timeStr := time.Now().Format("2006-01-02")
		//t, _ := time.Parse("2006-01-02", timeStr)
		//timeNumber := t.Unix()
		//timeNumber += 57600
		//fmt.Println(timeNumber)
		//_time := time.Now().Unix()
		//v := time.NewTimer(time.Second * (time.Duration(timeNumber - _time)))
		//
		//<-v.C
		var timer *time.Timer
		for {
			var ids []int
			db.Driver.Table("account_info").Select("id").Find(&ids)

			nTime := time.Now()
			yesTime := nTime.AddDate(0,0,-1)
			logDay := yesTime.Format("2006-01-02")
			fmt.Println(logDay)
			sql := squirrel.Insert("sign_in").Columns(
				"account_id", "date", "status",
			)
			for _,i :=range ids{
				var signIn db.SignIn
				if err := db.Driver.Where("account_id = ? and date = ?",i,logDay).First(&signIn).Error;err != nil{
					sql = sql.Values(
						i,
						logDay,
						false,
					)
				}
			}
			if s, args, err := sql.ToSql(); err != nil {
				logUtils.Println(err)
			} else {
				if err := db.Driver.Exec(s, args...).Error; err != nil {
					logUtils.Println(err)
					return
				}
			}
			timer = time.NewTimer(time.Hour*24)
			<-timer.C
		}

}