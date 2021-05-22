package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/party_class"
)

func RegisterPartyClassRouters(app *iris.Application){

	//先理解为前缀吧
	classesRouter := app.Party("v1/party_course")

	classesRouter.Post("/{pid:int}/class", hero.Handler(party_class.CreatePartyClass))
	classesRouter.Put("/class/{cid:int}", hero.Handler(party_class.PutPartyClass))
	classesRouter.Delete("/class/{cid:int}", hero.Handler(party_class.DeletePartyClass))
	classesRouter.Get("/class/list", hero.Handler(party_class.ListPartyClasses))
	classesRouter.Post("/class/_mget", hero.Handler(party_class.MgetPartyClass))
}
