package account

import (
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris"
	"grpc-demo/constants"
	authbase "grpc-demo/core/auth"
	"grpc-demo/core/cache"
	accountEnums "grpc-demo/enums/account"
	AccountException "grpc-demo/exceptions/account"
	"grpc-demo/models/db"
	"grpc-demo/utils/hash"
	mailUtils "grpc-demo/utils/mail"
	paramsUtils "grpc-demo/utils/params"
)

func Register(ctx iris.Context,auth authbase.AuthAuthorization)  {
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	account := auth.AccountModel()

	if account.EmailValidated != true{
		panic("email no check")
	}

	if params.Has("city") {
		city := params.Str("city", "city")
		account.City = city
	}
	if params.Has("avator"){
		avator := params.Str("avator", "avator")
		account.Avator = avator
	}
	if params.Has("country") {
		country := params.Str("country", "country")
		account.Country = country
	}
	if params.Has("province") {
		province := params.Str("province", "province")
		account.Province = province
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

func RegisterByEmail(ctx iris.Context)  {
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	email := params.Str("email", "email")

	v := hash.GetRandomString(5)
		//存入缓存
	_, _ = cache.Redis.Do(constants.DbNumberModel, "set", v, email,60*5)
	err := mailUtils.Send(v,email);if err != nil{
		panic(AccountException.AccountNotFount())
	}

	ctx.JSON(iris.Map{
		"status":"success",
	})

}

func IsEmailSend(ctx iris.Context,auth authbase.AuthAuthorization){
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	email  := params.Str("email","email")
	value  := params.Str("value","value")
	password := params.Str("password","password")
	nickname := params.Str("nickname","nickname")

	v, err := redis.String(cache.Redis.Do(constants.DbNumberModel, "get", value))
	if err == nil && v == email {
		var account db.Account
		account.EmailValidated = true
		account.Email = email
		account.Password = password
		account.Nickname = nickname
		db.Driver.Create(&account)
		   ctx.JSON(iris.Map{
			   "id":    account.Id,
		   })
	} else {
		panic(AccountException.AccountNotFount())
	}
}

func RegisterByPhone(ctx iris.Context,auth authbase.AuthAuthorization){
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	phone := params.Str("phone","phone")
	//todo 电话格式校验
	//todo 电话验证码

	password := params.Str("password","password")
	nickname := params.Str("nickname","nickname")

	var v db.Account
	if err := db.Driver.Where("nickname = ?",nickname).First(&v);err == nil{
		panic("nickname exist")
	}

	if err := db.Driver.Where("phone = ?",phone).First(&v);err == nil{
		panic("phone exist")
	}

	account := db.Account{
		Nickname:       nickname,
		Role:           accountEnums.RoleUser,
		Phone:          phone,
		PhoneValidated: false,
		Password:       password,
	}

	//todo 邮箱格式校验
	if params.Has("email"){
		account.Email = params.Str("email","email")
		account.EmailValidated = true
	}
	db.Driver.Create(&account)

	ctx.JSON(iris.Map{
		"id":account.Id,
	})

}





