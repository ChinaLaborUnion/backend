package account

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	accountException "grpc-demo/exceptions/account"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func GetAccountInfo(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckLogin()
	id := auth.AccountModel().Id
	var account db.AccountInfo
	err := db.Driver.Where("id = ?",id).First(&account)
	if err != nil{
		panic(accountException.AccountNotFount())
	}
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
	id := auth.AccountModel().Id
	var account db.AccountInfo
	err := db.Driver.GetOne("account_info",id,account)
	if err != nil{
		panic(accountException.AccountNotFount())
	}

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

func ListAccount(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckLogin()
	var lists []struct{
		ID         int   `json:"id"`
		CreateTime int64 `json:"create_time"`
	}
	var count int

	table := db.Driver.Table("account_info")
	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	//if !auth.IsAdmin() {
	//	table = table.Where("account_id = ?", auth.AccountModel().Id)
	//}
	//
	//if author := ctx.URLParamIntDefault("author_id", 0); author != 0 && auth.IsAdmin() {
	//	table = table.Where("account_id = ?", author)
	//}
	//
	//if status := ctx.URLParamIntDefault("status", 0); status != 0 {
	//	table = table.Where("status = ?", status)
	//}
	//
	//if startTime := ctx.URLParamInt64Default("start_time", 0); startTime != 0 {
	//	endTime := ctx.URLParamInt64Default("end_time", 0)
	//	table = table.Where("create_time between ? and ?", startTime, endTime)
	//}

	table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Select("id, create_time").Find(&lists)

	ctx.JSON(iris.Map{
		"orders": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

func MgetAccounts(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckAdmin()
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	orders := db.Driver.GetMany("account_info",ids,db.AccountInfo{})
	for _,o := range orders{
		func(data *[]interface{}){
			*data = append(*data,paramsUtils.ModelToDict(o,[]string{"Id","Avator","Nickname","Email"}))
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}