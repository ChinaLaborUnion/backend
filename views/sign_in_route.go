package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/sign_in"
)

func RegisterSignInRoute(app *iris.Application)  {
	signInRouter := app.Party("v1/sign_in")
	//报名信息路由
	signInRouter.Post("",hero.Handler(sign_in.SignIn))
	signInRouter.Get("/list",hero.Handler(sign_in.SignInList))

}
