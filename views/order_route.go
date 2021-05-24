package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/order"
)

func RegisterOrderRouters(app *iris.Application){
	orderRouter := app.Party("v1/order")

	orderRouter.Post("", hero.Handler(order.CreateOrder))
	orderRouter.Put("/{oid:int}", hero.Handler(order.CancelOrder))
	orderRouter.Get("/list", hero.Handler(order.ListOrder))
	orderRouter.Post("/_mget", hero.Handler(order.MgetOrders))
}
