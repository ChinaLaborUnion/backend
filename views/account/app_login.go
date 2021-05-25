package account

import (
	"fmt"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func AppLoginByEmail(ctx iris.Context,auth authbase.AuthAuthorization) {
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	var account db.Account
	email := params.Str("email","email")
	password := params.Str("password","password")
	if err := db.Driver.Where("email = ? and password = ?",email,password).First(&account).Error;err != nil{
		panic("account not exsit")
	}
	fmt.Println(account.Id)
	auth.SetCookie(account.Id)
	ctx.JSON(iris.Map{
		"id":account.Id,
	})
}

func LoginByPhone(ctx iris.Context,auth authbase.AuthAuthorization){
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	phone := params.Str("phone","phone")
	password := params.Str("password","password")

	var account db.Account

	if err := db.Driver.Where("phone = ? and pawword = ?",phone,password).First(&account);err == nil{
		panic("account not exist")
	}

	auth.SetCookie(account.Id)
	ctx.JSON(iris.Map{
		"id":account.Id,
	})
}

