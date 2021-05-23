package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/news"

	//上面这个应该是我写的对应CRUD方法的那个package
)

//设置路由,写完路由后，就要在main.go注册这个路由
func RegisterNewsRouters(app *iris.Application){
	newsRouter := app.Party("v1/news")

	newsRouter.Post("", hero.Handler(news.CreateNews))
	newsRouter.Put("/{nid:int}", hero.Handler(news.PutNews))
	newsRouter.Delete("/{nid:int}", hero.Handler(news.DeleteNews))
	newsRouter.Get("/list", hero.Handler(news.ListNews))
	newsRouter.Post("/_mget", hero.Handler(news.MgetNews))
}

//设置路由,写完路由后，就要在main.go注册这个路由
func RegisterNewsLabelRouters(app *iris.Application){
	newsRouter := app.Party("v1/news_label")

	newsRouter.Post("", hero.Handler(news.CreateNewsLabel))
	newsRouter.Put("/{nlid:int}", hero.Handler(news.PutNewsLabel))
	newsRouter.Delete("/{nlid:int}", hero.Handler(news.DeleteNewsLabel))
	newsRouter.Get("/list", hero.Handler(news.ListNewsLabel))
}

