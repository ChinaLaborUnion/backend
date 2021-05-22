package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/party_course"
)

func CoursePictureRouters(app *iris.Application) {
	coursePictureRouters := app.Party("v1/party_course/picture")

	coursePictureRouters.Post("", hero.Handler(party_course.CreateCoursePicture))
	coursePictureRouters.Put("",hero.Handler(party_course.PutCoursePicture))
	coursePictureRouters.Get("", hero.Handler(party_course.GetCoursePicture))
}