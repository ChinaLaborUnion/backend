package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/homework"
)

func HomeWorkRoute(app *iris.Application)  {
	HomeWorkRouter := app.Party("/homework")
	//报名信息路由
	HomeWorkRouter.Post("/create",hero.Handler(homework.CreatHomeWork))
	HomeWorkRouter.Put("/{uid:int}",hero.Handler(homework.PutHomeWork))
	HomeWorkRouter.Get("/list/{uid:int}",hero.Handler(homework.HomeWorkList))
	HomeWorkRouter.Delete("/{cid:int}",hero.Handler(homework.DeleteHomeWork))
	HomeWorkRouter.Post("/megt",hero.Handler(homework.HomeWorkMegt))
}
