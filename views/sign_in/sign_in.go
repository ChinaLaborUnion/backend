package sign_in

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	"grpc-demo/models/db"
	"time"
)

func SignIn(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckLogin()

	date := time.Now().Format("2006-01-02")

	var signIn db.SignIn
	if err := db.Driver.Where("account_id = ? and date = ?",auth.AccountModel().Id,date).First(&signIn).Error;err != nil{
		s := db.SignIn{
			AccountID:  auth.AccountModel().Id,
			Date:      date,
		}
		db.Driver.Create(&s)
	}


	ctx.JSON(iris.Map{
		"status":"success",
	})
}

func SignInList(ctx iris.Context,auth authbase.AuthAuthorization){
	auth.CheckLogin()

	var lists []struct {
		Id int `json:"id"`
		AccountID int `json:"account_id"`
		Date string `json:"date"`
		CreateTime int64 `json:"create_time"`
	}

	var count int

	table := db.Driver.Table("sign_in")

	if !auth.IsAdmin(){
		table = table.Where("account_id = ?",auth.AccountModel().Id)
	}

	if author := ctx.URLParamIntDefault("author_id", 0); author != 0 && auth.IsAdmin() {
		table = table.Where("account_id = ?", author)
	}

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)



	table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Select("id, account_id,date,create_time").Find(&lists)
	ctx.JSON(iris.Map{
		"sign_in": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}
