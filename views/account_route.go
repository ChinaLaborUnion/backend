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
	accountRouter.Post("/email_valid", hero.Handler(account.EmailRegistered))
	accountRouter.Post("/login_by_email", hero.Handler(account.AppLoginByEmail))

	//用户信息路由
	infoRouter := app.Party("v1/account/info")
	infoRouter.Post("/reset_email", hero.Handler(account.ResetEmail))
	infoRouter.Get("/get", hero.Handler(account.GetAccountInfo))
	infoRouter.Put("/put", hero.Handler(account.PutAccountInfo))
	infoRouter.Get("/list", hero.Handler(account.ListAccount))
	infoRouter.Post("/_mget", hero.Handler(account.MgetAccounts))

}
