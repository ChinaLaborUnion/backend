package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/homework"
)

func HomeWorkRoute(app *iris.Application)  {
	HomeWorkRouter := app.Party("v1/party_course/class/homework")
	//报名信息路由
	HomeWorkRouter.Post("/{cid:int}",hero.Handler(homework.CreateHomeWork))
	HomeWorkRouter.Put("/{hid:int}",hero.Handler(homework.PutHomeWork))
	HomeWorkRouter.Get("/list",hero.Handler(homework.HomeWorkList))
	HomeWorkRouter.Delete("/{hid:int}",hero.Handler(homework.DeleteHomeWork))
	HomeWorkRouter.Post("/_mget",hero.Handler(homework.HomeWorkMegt))
}
