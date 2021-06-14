package views

import (

	"grpc-demo/views/transaction"
	"grpc-demo/views/transaction/payment/ali"
	"grpc-demo/views/transaction/payment/wx"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
)

//设置路由,写完路由后，就要在main.go注册这个路由
func RegisterTransactionRouters(app *iris.Application){
	transactionRouter := app.Party("v1/transaction")

	transactionRouter.Post("/aio/{oid:int}/{tid:int}",hero.Handler(transaction.PaymentMiddleware))

	//微信回调
	transactionRouter.Post("/wx/callback" ,hero.Handler(wx.CallbackReceiver))
	//支付宝回调
	transactionRouter.Post("/ali/callback", hero.Handler(ali.PaymentCallbackReceiver))

	//微信回调
	transactionRouter.Post("wx/paycallback",hero.Handler(wx.PaymentCallbackReceiver))
}

