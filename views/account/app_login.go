package account

import (
	"github.com/kataras/iris"
	"grpc-demo/constants"
	authbase "grpc-demo/core/auth"
	AccountException "grpc-demo/exceptions/account"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func AppLoginByEmail(ctx iris.Context,auth authbase.AuthAuthorization) {
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	var account db.AccountInfo
	email := params.Str("email","email")
	password := params.Str("password","password")
	if err := db.Driver.Where("email = ? and password = ?",email,password).First(&account).Error;err != nil{
		panic(AccountException.WrongInput())
	}

	token := auth.SetCookie(account.Id)
	mode := ctx.GetHeader(constants.ApiMode)

	if mode == "app" {
		ctx.JSON(iris.Map{
			"token":token,
			"id":account.Id,
		})
	}else{
		ctx.JSON(iris.Map{
			"id":account.Id,
		})
	}
}

func LoginByPhone(ctx iris.Context,auth authbase.AuthAuthorization){
	//params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	//phone := params.Str("phone","phone")
	//password := params.Str("password","password")
	//
	//var account db.Account
	//
	//if err := db.Driver.Where("phone = ? and pawword = ?",phone,password).First(&account);err == nil{
	//	panic("account not exist")
	//}
	//
	//auth.SetCookie(account.Id)
	//ctx.JSON(iris.Map{
	//	"id":account.Id,
	//})
}

