package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/party_course"
)

func PartyCourseRouters(app *iris.Application) {
	partyCourseRouters := app.Party("v1/party_course/info")

	partyCourseRouters.Post("", hero.Handler(party_course.CreatePartyCourse))
	partyCourseRouters.Put("/{cid:int}",hero.Handler(party_course.PutPartyCourse))
	partyCourseRouters.Delete("/{cid:int}",hero.Handler(party_course.DeletePartyCourse))
	partyCourseRouters.Get("/list", hero.Handler(party_course.ListPartyCourse))
	partyCourseRouters.Post("/_mget", hero.Handler(party_course.MgetPartyCourse))




}