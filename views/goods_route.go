package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/goods"
)

func RegisterGoodsRouters(app *iris.Application){
	goodsRouter := app.Party("v1/goods")

	goodsRouter.Post("", hero.Handler(goods.CreateGoods))
	goodsRouter.Put("/{gid:int}", hero.Handler(goods.PutGoods))
	goodsRouter.Delete("/{gid:int}", hero.Handler(goods.DeleteGoods))
	goodsRouter.Get("/list", hero.Handler(goods.ListGoods))
	goodsRouter.Post("/_mget", hero.Handler(goods.MgetGoods))
}
