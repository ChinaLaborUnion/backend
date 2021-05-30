package account

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func GetAccountInfo(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckLogin()
	account := auth.AccountModel()

	ctx.JSON(iris.Map{
		"id": account.Id,
		"email": account.Email,
		"avator": account.Avator,
		"nickname":  account.Nickname,
		"phone":account.Phone,
	})
}

func PutAccountInfo(ctx iris.Context,auth authbase.AuthAuthorization)  {
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	account := auth.AccountModel()

	if account.EmailValidated != true{
		panic("email no check")
	}

	if params.Has("avator"){
		avator := params.Str("avator", "avator")
		account.Avator = avator
	}

	if params.Has("nickname") {
		nickname := params.Str("nickname", "nickName")
		account.Nickname = nickname
	}

	db.Driver.Save(&account)
	ctx.JSON(
		iris.Map{
			"id":    account.Id,
	})
}

