package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"

	"grpc-demo/views/account"

)

func RegisterAccountRouters(app *iris.Application) {

	//登录路由
	accountRouter := app.Party("v1/account/login")

	//app登陆路由
	accountRouter.Post("/register_by_email", hero.Handler(account.RegisterByEmail))
	accountRouter.Post("/email_valid", hero.Handler(account.IsEmailSend))
	accountRouter.Post("/login_by_email", hero.Handler(account.AppLoginByEmail))

}
