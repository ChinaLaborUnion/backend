package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"

	"grpc-demo/views/account"

)

func RegisterAccountRouters(app *iris.Application) {

	//登录路由
	loginRouter := app.Party("account/login")

	//app登陆路由
	loginRouter.Post("/register_by_email", hero.Handler(account.RegisterByEmail))
	loginRouter.Post("/app_register", hero.Handler(account.Register))
	loginRouter.Post("/email_valid", hero.Handler(account.IsEmailSend))
	loginRouter.Post("/login_by_email", hero.Handler(account.AppLoginByEmail))
	loginRouter.Post("/register_by_phone", hero.Handler(account.RegisterByPhone))
	loginRouter.Post("/login_by_phone", hero.Handler(account.LoginByPhone))

}
