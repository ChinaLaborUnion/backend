package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/sign_up"
)

func RegisterSignUpRoute(app *iris.Application)  {
	signupRouter := app.Party("v1/party_course/class/sign_up")
	//报名信息路由
	signupRouter.Post("",hero.Handler(sign_up.CreatSignUp))
	signupRouter.Put("/{sid:int}",hero.Handler(sign_up.PutSignUp))
	signupRouter.Get("/list_by_cid/{cid:int}",hero.Handler(sign_up.SignUpListByCid))
	signupRouter.Delete("/{sid:int}",hero.Handler(sign_up.DeleteSignUp))
	signupRouter.Get("/list_by_aid/{aid:int}",hero.Handler(sign_up.SignUpListByAid))
}
